package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {

	var rootCmd = &cobra.Command{
		Use:   "home",
		Short: "home is a tool to control a custom home server with multiple smart home integrations",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	var serverCmd = &cobra.Command{
		Use:   "serve",
		Short: "Run the home server",
		Run: func(cmd *cobra.Command, args []string) {
			server(8600)
		},
	}

	var connectCmd = &cobra.Command{
		Use:   "connect",
		Short: "Establish a connection with the server as a client",
		Run: func(cmd *cobra.Command, args []string) {
			connect()
		},
	}

	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(connectCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
		server(8600)
	}
}
