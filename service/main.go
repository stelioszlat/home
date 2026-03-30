package main

import (
	"fmt"
	"os"

	"C"

	"github.com/spf13/cobra"
)
import (
	"context"
	"service/dashboard"
	"service/internal/chat"
	"service/internal/server"
	"service/notifications"
	"time"
)

type HomeService struct {
	ctx    *context.Context
	server *server.Server
}

func main() {
	var port int
	var deviceName string
	var deviceId string
	var allDevices bool
	var message string
	chatInstance := chat.New()

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
			serverInstance := server.New()
			serverInstance.Serve(port)
		},
	}

	var dashboardCmd = &cobra.Command{
		Use:   "dashboard",
		Short: "Run the terminal dashboard",
		Run: func(cmd *cobra.Command, args []string) {
			dashboard.Dashboard()
		},
	}

	var notificationCmd = &cobra.Command{
		Use:   "notify",
		Short: "Send notification to a device",
		Run: func(cmd *cobra.Command, args []string) {
			notifications.Notify(notifications.NotificationArgs{All: allDevices, DeviceId: deviceId, DeviceName: deviceName}, notifications.Notification{Message: message, App: "homecli", CreatedAt: time.Now()})
		},
	}

	var chatCmd = &cobra.Command{
		Use:   "chat",
		Short: "Chat with the home assistant",
		Run: func(cmd *cobra.Command, args []string) {
			chatInstance.Chat(args[0], "deepseek-r1:1.5b")
		},
	}

	var abacusCmd = &cobra.Command{
		Use:   "abacus",
		Short: "Connect to Abacus AI chat server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Use one of the subcommands to interact with Abacus AI")
		},
	}

	var modelsCmd = &cobra.Command{
		Use:   "models",
		Short: "List available models",
		Run: func(cmd *cobra.Command, args []string) {
			chatInstance.Models()
		},
	}

	var pipelinesCmd = &cobra.Command{
		Use:   "pipelines",
		Short: "List available pipelines from Abacus AI",
		Run: func(cmd *cobra.Command, args []string) {
			pipelines := server.GetAbacusPipelines()
			fmt.Println("Available pipelines from Abacus AI:")
			for _, pipeline := range pipelines {
				fmt.Println("-", pipeline)
			}
		},
	}

	serverCmd.PersistentFlags().IntVarP(&port, "port", "p", 8060, "Port to run the server on (shorthand), default is 8600")
	notificationCmd.PersistentFlags().StringVarP(&deviceName, "device", "d", "", "Device to send an instant notification")
	notificationCmd.PersistentFlags().StringVar(&deviceId, "id", "", "Device ID to send an isntant notification")
	notificationCmd.PersistentFlags().BoolVarP(&allDevices, "all", "a", false, "Send an instant notifications to all devices")
	notificationCmd.PersistentFlags().StringVarP(&message, "message", "m", "", "Message to send through the notification")
	abacusCmd.AddCommand(modelsCmd)
	abacusCmd.AddCommand(pipelinesCmd)
	chatCmd.AddCommand(modelsCmd)
	rootCmd.AddCommand(chatCmd)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(dashboardCmd)
	rootCmd.AddCommand(notificationCmd)
	rootCmd.AddCommand(abacusCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
