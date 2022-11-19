package funrun

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

// RestartPolicy represents the restart behavior of a process.
type RestartPolicy string

const (
	RestartNever  RestartPolicy = "never"   // Never restart the process
	RestartOnFail RestartPolicy = "on-fail" // Restart the process if it exits with a non-zero exit code
	RestartAlways RestartPolicy = "always"  // Always restart the process
)

type ProcConf struct {
	Name    string        `yaml:"name"`    // Name of the process
	Cmd     string        `yaml:"cmd"`     // Command to run (required)
	Args    []string      `yaml:"args"`    // Arguments to pass to the command
	Envs    []string      `yaml:"envs"`    // Environment variables to set
	Restart RestartPolicy `yaml:"restart"` // Restart policy for the command
	WorkDir string        `yaml:"workdir"` // Working directory for the command
}

type Conf struct {
	Procs []ProcConf `yaml:"procs"`
}

func ReadConf(path string) (*Conf, error) {
	// Open the file
	f := os.Stdin

	// If a path is given, open the file, otherwise use stdin
	if path != "" && path != "-" {
		var err error
		f, err = os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open config file: %w", err)
		}
		defer f.Close()
	}

	// Read the file
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse the file
	var conf Conf
	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate and set defaults...
	for i, p := range conf.Procs {
		// Check that a command is set...
		if p.Cmd == "" {
			return nil, fmt.Errorf("missing command for process %d", i)
		}

		// Set the default restart policy
		if p.Restart == "" {
			p.Restart = RestartNever
		}

		// Set a name if not set...
		if p.Name == "" {
			p.Name = fmt.Sprintf("proc-%d", i)
		}

		// Store the updated process
		conf.Procs[i] = p
	}

	// Return the config successfully!
	return &conf, nil
}

func (c *Conf) maxNameLength() int {
	var max int
	for _, p := range c.Procs {
		if len(p.Name) > max {
			max = len(p.Name)
		}
	}
	return max
}
