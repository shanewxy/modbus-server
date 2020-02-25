package main

import (
	rtu "modbus-demo/rtu/cmd"
	tcp "modbus-demo/tcp/cmd"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	name        = "modbus"
	description = ``
)

var allCommands = []*cobra.Command{
	tcp.NewCommand(),
	rtu.NewCommand(),
}

func main() {
	var c = &cobra.Command{
		Use:  name,
		Long: description,
		RunE: func(cmd *cobra.Command, args []string) error {

			var (
				basename  = filepath.Base(os.Args[0])
				targetCmd *cobra.Command
			)
			for _, cmd := range allCommands {
				if cmd.Name() == basename {
					targetCmd = cmd
					break
				}
				for _, alias := range cmd.Aliases {
					if alias == basename {
						targetCmd = cmd
						break
					}
				}
			}
			if targetCmd != nil {
				return targetCmd.Execute()
			}
			return cmd.Help()
		},
	}
	c.AddCommand(allCommands...)

	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
