package funrun

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// A command's status
type CmdStatus int

const (
	CmdNotStarted CmdStatus = iota // The command hasn't started yet
	CmdRunning                     // The command is running
	CmdDone                        // The command has finished successfully
	CmdFailed                      // The command has finished with an error
	CmdStopped                     // The command has been stopped
)

type Command struct {
	conf   *ProcConf          // The configuration for this command
	cmd    *exec.Cmd          // The command object
	err    error              // The error returned by the command
	status CmdStatus          // The current status of the command
	cancel context.CancelFunc // The cancel function for the command context
	wout   *PrefixWriter
	werr   *PrefixWriter
	sync.RWMutex
}

func NewCommand(conf *ProcConf) *Command {
	return &Command{
		conf: conf,
	}
}

func (c *Command) Name() string {
	return c.conf.Name
}

func (c *Command) SetOutputs(wout, werr *PrefixWriter) {
	c.Lock()
	defer c.Unlock()
	c.wout = wout
	c.werr = werr
}

func (c *Command) makeEnvGetter() func(string) string {
	return func(key string) string {
		// Check for an environment variable set explicitly
		v, ok := c.conf.Envs[key]
		if ok {
			return v
		}

		// Check if the rest of the environment is available
		if c.conf.ClearEnvs {
			return ""
		}

		// Otherwise, return the environment variable
		return os.Getenv(key)
	}
}

func (c *Command) fmtEnvSlice() []string {
	// Create a slice for env vars...
	var envs []string

	// Add the rest of the environment if we're not clearing it...
	if !c.conf.ClearEnvs {
		for _, v := range os.Environ() {
			envs = append(envs, v)
		}
	}

	// Add the explicit env vars...
	for k, v := range c.conf.Envs {
		e := fmt.Sprintf("%s=%s", k, v)
		envs = append(envs, e)
	}

	// Return the slice
	return envs
}

func (c *Command) makeSingleCmd(ctx context.Context) *exec.Cmd {
	envGetter := c.makeEnvGetter()

	// Expand the command text
	cmdTxt := os.Expand(c.conf.Cmd, envGetter)

	// Expand the arguments
	args := make([]string, len(c.conf.Args))
	for i, arg := range c.conf.Args {
		args[i] = os.Expand(arg, envGetter)
	}

	// Create the command...
	return exec.CommandContext(
		ctx,
		cmdTxt,
		args...,
	)
}

func (c *Command) makeMultiCmd(ctx context.Context) *exec.Cmd {
	// Get the env getter...
	envGetter := c.makeEnvGetter()

	// Expand the commands...
	cmds := make([]string, len(c.conf.Cmds))
	for i, cmd := range c.conf.Cmds {
		cmds[i] = os.Expand(cmd, envGetter)
	}

	// Expand the arguments...
	args := make([]string, len(c.conf.Args))
	for i, arg := range c.conf.Args {
		args[i] = os.Expand(arg, envGetter)
	}

	// Join the commands...
	cmdTxt := strings.Join(cmds, "; ")

	// Join together the command and the arguments...
	allArgs := append(
		[]string{"-c", cmdTxt},
		c.conf.Args...,
	)

	// Create the command and return...
	return exec.CommandContext(
		ctx,
		"sh",
		allArgs...,
	)
}

func (c *Command) createCmd(ctx context.Context) *exec.Cmd {
	// Create the command with the context
	var cmd *exec.Cmd
	if c.conf.Cmd == "" {
		cmd = c.makeMultiCmd(ctx)
	} else {
		cmd = c.makeSingleCmd(ctx)
	}

	if c.conf.Cmd == "" {
		cmdTxt := strings.Join(c.conf.Cmds, "; ")
		args := append(
			[]string{"-c", cmdTxt},
			c.conf.Args...,
		)
		cmd = exec.CommandContext(ctx, "sh", args...)
	}

	// Set the working directory
	cmd.Dir = c.conf.WorkDir
	if cmd.Dir == "" {
		cmd.Dir = "."
	}

	// Add the environment variables
	cmd.Env = c.fmtEnvSlice()

	// Set the outputs
	cmd.Stdout = c.wout
	cmd.Stderr = c.werr

	// Return the command
	return cmd
}

func (c *Command) setStatus(status CmdStatus) {
	c.Lock()
	defer c.Unlock()
	c.status = status
}

func (c *Command) Status() CmdStatus {
	c.RLock()
	defer c.RUnlock()
	return c.status
}

func (c *Command) Run(ctx context.Context) error {
	// Start a loop...
runloop:
	for {
		// Create the context for the command
		ctx, cancel := context.WithCancel(ctx)
		c.cancel = cancel

		// Create the command
		c.cmd = c.createCmd(ctx)

		select {
		case <-ctx.Done():
			// The context has been cancelled
			c.setStatus(CmdStopped)
			break runloop

		default:
			c.wout.Logf("Starting...\n")

			// Start the command
			c.setStatus(CmdRunning)
			err := c.cmd.Start()
			if err != nil {
				// Store the error and cmd state
				c.setStatus(CmdFailed)
				c.err = err

				// Log the error
				c.wout.Logf("Error starting command: %s\n", err)

				// Should we restart?
				if c.conf.Restart == RestartOnFail || c.conf.Restart == RestartAlways {
					c.wout.Logf("Restarting...\n")
					continue runloop
				}

				// Otherwise, break out of the loop
				return err
			}

			// Wait for the command to finish
			err = c.cmd.Wait()
			if err != nil {
				c.err = err
				c.setStatus(CmdFailed)
			} else {
				c.setStatus(CmdDone)
				c.err = nil
			}

			// Should we restart?
			if c.conf.Restart == RestartAlways || (c.conf.Restart == RestartOnFail && c.err != nil) {
				c.wout.Logf("Restarting...\n")
				continue runloop
			}

			// Otherwise, break out of the loop
			c.wout.Logf("Finished\n")
			break runloop
		}
	}

	// Return success!
	return nil
}

func (c *Command) Cancel() {
	c.setStatus(CmdStopped)
	if c.cancel != nil {
		c.cancel()
	}
}

func (c *Command) setError(err error) {
	c.Lock()
	defer c.Unlock()
	c.err = err
}

func (c *Command) Error() error {
	c.RLock()
	defer c.RUnlock()
	return c.err
}
