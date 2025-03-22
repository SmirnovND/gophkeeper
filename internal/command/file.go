package command

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

// filePathCompleter предоставляет автодополнение для путей к файлам
func filePathCompleter(d prompt.Document) []prompt.Suggest {
	path := d.Text
	dir := "."
	base := ""

	if path != "" {
		dir = filepath.Dir(path)
		base = filepath.Base(path)

		// Если путь заканчивается на /, то мы ищем в этой директории
		if strings.HasSuffix(path, "/") || strings.HasSuffix(path, "\\") {
			dir = path
			base = ""
		}
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return []prompt.Suggest{}
	}

	var suggestions []prompt.Suggest
	for _, file := range files {
		name := file.Name()
		if strings.HasPrefix(name, base) {
			var text string
			if file.IsDir() {
				text = filepath.Join(dir, name) + "/"
			} else {
				text = filepath.Join(dir, name)
			}

			// Если путь начинается с ./, удаляем его для лучшего отображения
			if strings.HasPrefix(text, "./") {
				text = text[2:]
			}

			suggestions = append(suggestions, prompt.Suggest{
				Text:        text,
				Description: getFileDescription(file),
			})
		}
	}
	return suggestions
}

// getFileDescription возвращает описание файла (директория или размер файла)
func getFileDescription(file os.DirEntry) string {
	if file.IsDir() {
		return "Директория"
	}

	info, err := file.Info()
	if err != nil {
		return "Файл"
	}

	size := info.Size()
	if size < 1024 {
		return fmt.Sprintf("Файл (%d байт)", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("Файл (%.2f КБ)", float64(size)/1024)
	} else {
		return fmt.Sprintf("Файл (%.2f МБ)", float64(size)/(1024*1024))
	}
}

// emptyCompleter - пустой комплитер для полей без автодополнения
func emptyCompleter(d prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{}
}

func (c *Command) UploadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upload",
		Short: "Хранение текстовых/бинарных данных",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Введите путь и имя файла (используйте TAB для автодополнения):")
			filePath := prompt.Input("> ", filePathCompleter)

			fmt.Println("Введите уникальное название(label) сохраняемого объекта:")
			label := prompt.Input("> ", emptyCompleter)

			resp, err := c.clientUseCase.Upload(filePath, label)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(resp)
		},
	}
}
