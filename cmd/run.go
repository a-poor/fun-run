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
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run {- | CONFIG_FILE}",
	Short: "Start running your commands.",
	Long: `Start running your commands based on the configuration file.

To read from stdin, use '-' as the CONFIG_FILE argument.`,
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"r"},
	Run: func(cmd *cobra.Command, args []string) {
		// Get the config file (or read from stdin)...
		var p string
		if len(args) > 0 {
			p = args[0]
		}

		// Load the config...
		conf, err := funrun.ReadConf(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
			os.Exit(1)
		}

		// Create the process manager...
		man := funrun.NewManager(conf)

		// Get the context...
		ctx := cmd.Context()

		// Start the processes...
		if err := man.Run(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Error running processes: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
