package main

import (
	"fmt"
	"github.com/alekseinovikov/gonner/configs"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
	"runtime"
	"sync"
)

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(upCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of the process manager",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Gonner Process Manager " + configs.Version)
	},
}

var rootCmd = &cobra.Command{
	Use:   "gonner",
	Short: "A lightweight process manager",
	Long:  `Gonner Process Manager helps you manage and monitor processes on your server.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Gonner Process Manager! Use --help to see available commands.")
	},
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start all processes defined in the config",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := configs.LoadConfig("configs/config.yml")
		if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}

		var wg sync.WaitGroup

		for _, service := range config.Services {
			wg.Add(1)
			go func(service configs.Service) {
				defer wg.Done()
				startService(service)
			}(service)
		}

		// Ждем завершения всех горутин
		wg.Wait()
		fmt.Println("All services have been started.")
	},
}

func startService(service configs.Service) {
	fmt.Printf("Starting service: %s\n", service.Name)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// Используем cmd.exe для Windows
		cmd = exec.Command("cmd", "/C", service.Command)
	} else {
		// Используем sh для Unix-подобных систем
		cmd = exec.Command("sh", "-c", service.Command)
	}

	cmd.Env = append(cmd.Env, service.Env...)

	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start service %s: %v", service.Name, err)
		return
	}

	log.Printf("Service %s started with PID %d", service.Name, cmd.Process.Pid)

	// Ждем завершения процесса
	if err := cmd.Wait(); err != nil {
		log.Printf("Service %s exited with error: %v", service.Name, err)
	} else {
		log.Printf("Service %s completed successfully", service.Name)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
