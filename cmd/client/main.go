package main

import (
	"fmt"
	"github.com/SmirnovND/gophkeeper/internal/command"
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
	// Устанавливаем информацию о версии и дате сборки
	command.SetVersionInfo(version, buildDate)
	
	diContainer := client.NewContainer(serverAddress)
	var Command interfaces.Command
	diContainer.Invoke(func(cmd interfaces.Command) {
		Command = cmd
	})
	rootCmd.AddCommand(Command.Login())
	rootCmd.AddCommand(Command.RegisterCmd())
	rootCmd.AddCommand(Command.UploadCmd())
	rootCmd.AddCommand(Command.DownloadCmd())
	
	// Добавляем команды для работы с текстовыми данными
	rootCmd.AddCommand(Command.SaveTextCmd())
	rootCmd.AddCommand(Command.GetTextCmd())
	rootCmd.AddCommand(Command.DeleteTextCmd())
	
	// Добавляем команды для работы с данными кредитных карт
	rootCmd.AddCommand(Command.SaveCardCmd())
	rootCmd.AddCommand(Command.GetCardCmd())
	rootCmd.AddCommand(Command.DeleteCardCmd())
	
	// Добавляем команды для работы с учетными данными
	rootCmd.AddCommand(Command.SaveCredentialCmd())
	rootCmd.AddCommand(Command.GetCredentialCmd())
	rootCmd.AddCommand(Command.DeleteCredentialCmd())
	
	// Добавляем команду для получения информации о версии
	rootCmd.AddCommand(Command.VersionCmd())
	
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
