package command

import (
	"fmt"
	"github.com/spf13/cobra"
)

// Переменные, которые будут устанавливаться при сборке
var (
	version   string
	buildDate string
)

// SetVersionInfo устанавливает информацию о версии и дате сборки
func SetVersionInfo(v, date string) {
	version = v
	buildDate = date
}

// VersionCmd возвращает команду для получения информации о версии
func (c *Command) VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Показать информацию о версии и дате сборки",
		Long:  "Показать информацию о версии и дате сборки клиента",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Версия: %s\n", version)
			fmt.Printf("Дата сборки: %s\n", buildDate)
		},
	}
}