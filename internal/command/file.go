package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func (c *Command) UploadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upload",
		Short: "Хранение текстовых/бинарных данных",
		Run: func(cmd *cobra.Command, args []string) {
			var filePath, label string

			fmt.Println("Введите путь и имя файла:")
			fmt.Print("> ")
			fmt.Fscanln(os.Stdin, &filePath)

			fmt.Println("Введите уникальное название (label) сохраняемого объекта:")
			fmt.Print("> ")
			fmt.Fscanln(os.Stdin, &label)

			resp, err := c.clientUseCase.Upload(filePath, label)
			if err != nil {
				fmt.Println("Ошибка:", err)
				return
			}

			fmt.Println("Файл успешно загружен:", resp)
		},
	}
}

func (c *Command) DownloadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "download",
		Short: "Скачивание файла с сервера",
		Run: func(cmd *cobra.Command, args []string) {
			var label string

			fmt.Println("Введите уникальное название (label) файла для скачивания:")
			fmt.Print("> ")
			fmt.Fscanln(os.Stdin, &label)

			// Передаем пустую строку в качестве пути, чтобы использовать директорию загрузок по умолчанию
			err := c.clientUseCase.Download(label)
			if err != nil {
				fmt.Println("Ошибка при скачивании файла:", err)
				return
			}
		},
	}
}
