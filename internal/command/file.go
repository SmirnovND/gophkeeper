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
