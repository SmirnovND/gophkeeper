package command

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"path/filepath"
)

// filePathCompleter предоставляет автодополнение для путей к файлам
func filePathCompleter(d prompt.Document) []prompt.Suggest {
	path := filepath.Clean(d.Text)
	if path == "." {
		path = ""
	}

	files, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return nil
	}

	var suggestions []prompt.Suggest
	for _, file := range files {
		suggestions = append(suggestions, prompt.Suggest{Text: file})
	}
	return suggestions
}

func (c *Command) UploadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upload",
		Short: "Хранение текстовых/бинарных данных",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Введите путь и имя файла (используйте TAB для автодополнения):")
			filePath := prompt.Input("> ", filePathCompleter)

			fmt.Println("Введите уникальное название (label) сохраняемого объекта:")
			label := prompt.Input("> ", func(d prompt.Document) []prompt.Suggest { return nil })

			resp, err := c.clientUseCase.Upload(filePath, label)
			if err != nil {
				fmt.Println("Ошибка:", err)
				return
			}

			fmt.Println("Файл успешно загружен:", resp)
		},
	}
}
