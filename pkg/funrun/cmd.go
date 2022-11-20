package funrun

import (
	"context"
	"os/exec"
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

func (c *Command) createCmd(ctx context.Context) *exec.Cmd {
	// Create the command with the context
	cmd := exec.CommandContext(ctx, c.conf.Cmd, c.conf.Args...)

	// Set the working directory
	cmd.Dir = c.conf.WorkDir
	if cmd.Dir == "" {
		cmd.Dir = "."
	}

	// Add the environment variables
	cmd.Env = append(cmd.Env, c.conf.Envs...)

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
			c.wout.Logf("Starting command...\n")

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
					c.wout.Logf("Restarting command\n")
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
				c.wout.Logf("Restarting command\n")
				continue runloop
			}

			// Otherwise, break out of the loop
			c.wout.Logf("Command finished\n")
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
