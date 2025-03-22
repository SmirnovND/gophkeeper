package command

import (
	"fmt"
	"github.com/spf13/cobra"
)

func (c *Command) UploadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upload",
		Short: "Хранение текстовых/бинарных данных",
		Run: func(cmd *cobra.Command, args []string) {
			var filePath string
			fmt.Print("Введите путь и имя файла: ")
			fmt.Scanln(&filePath)

			var label string
			fmt.Print("Введите уникальное название(label) сохраняемого объекта: ")
			fmt.Scanln(&label)

			resp, err := c.clientUseCase.Upload(filePath, label)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(resp)
		},
	}
}
