package main

import (
	"fmt"
	"github.com/SmirnovND/gophkeeper/internal/container/client"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"os"

	"github.com/spf13/cobra"
)

var (
	version       string = "dev"
	buildDate     string = "unknown"
	serverAddress string = "127.0.0.1:8080"
)

var rootCmd = &cobra.Command{
	Use:   "passcli",
	Short: "CLI клиент для управления паролями",
}

func main() {
	fmt.Println(serverAddress)
	diContainer := client.NewContainer(serverAddress)
	var Command interfaces.Command
	diContainer.Invoke(func(cmd interfaces.Command) {
		Command = cmd
	})
	rootCmd.AddCommand(Command.Login())
	rootCmd.AddCommand(Command.UploadCmd())
	rootCmd.AddCommand(Command.DownloadCmd())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
