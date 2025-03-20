package main

import (
	"fmt"
	"github.com/SmirnovND/gophkeeper/internal/container/client"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "passcli",
	Short: "CLI клиент для управления паролями",
}

func main() {
	diContainer := client.NewContainer()
	var Command interfaces.Command
	diContainer.Invoke(func(cmd interfaces.Command) {
		Command = cmd
	})
	rootCmd.AddCommand(Command.Login())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
