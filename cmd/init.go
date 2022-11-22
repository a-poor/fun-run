/*
Copyright © 2022 Austin Poor <code@austinpoor.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/a-poor/fun-run/pkg/funrun"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [{- | CONFIG_FILE_PATH}]",
	Short: "Initialize a new fun-run config file.",
	Long: `Initialize a new fun-run config file. If no config
file path is provided, the default config file path 
("fun-run.yaml") will be used.

To print to stdout instead of a file, use "-" as the 
config file path.`,
	Args:    cobra.MaximumNArgs(1),
	Aliases: []string{"i"},
	Run: func(cmd *cobra.Command, args []string) {
		// Create the sample config file...
		cfg := funrun.Conf{
			Procs: []*funrun.ProcConf{
				{
					Name: "say-hello",
					Cmds: []string{
						"echo hello...",
						"sleep 1",
						"echo ...world",
						"sleep 1",
					},
					Restart: funrun.RestartAlways,
				},
				{
					Name:    "print-the-date",
					Cmd:     "date",
					Restart: funrun.RestartNever,
				},
				{
					Name: "greet-fun-run",
					Cmd:  "echo",
					Args: []string{
						"Hello, ${NAME}!",
					},
					Envs: map[string]string{
						"NAME": "fun-run",
					},
					Restart: funrun.RestartOnFail,
				},
			},
		}

		// Get the filepath to write to...
		p := "fun-run.yaml"
		if len(args) > 0 {
			p = args[0]
		}

		// Write the config file...
		b, err := yaml.Marshal(&cfg)
		if err != nil {
			fmt.Printf("Error marshaling config file as yaml: %s\n", err)
			os.Exit(1)
		}

		// Should we write to stdout? (Instead of a file)
		if p == "-" {
			cmd.Println(string(b))
			return
		}

		// Create the file (or overwrite it if it exists)...
		f, err := os.Create(p)
		if err != nil {
			cmd.Printf("Error creating config file: %s\n", err)
			os.Exit(1)
		}
		defer f.Close()

		// Write the config file...
		_, err = f.Write(b)
		if err != nil {
			fmt.Printf("Error writing config file: %s\n", err)
			f.Close()
			os.Exit(1)
		}

		// Done!
		cmd.Printf("New fun-run config file written to: %s\n", p)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
