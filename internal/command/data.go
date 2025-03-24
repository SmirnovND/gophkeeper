package command

import (
	"fmt"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/spf13/cobra"
	"os"
)

// SaveTextCmd создает команду для сохранения текстовых данных
func (c *Command) SaveTextCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "save-text",
		Short: "Сохранение текстовых данных",
		Run: func(cmd *cobra.Command, args []string) {
			var label, content string

			fmt.Println("Введите уникальное название (label) для текстовых данных:")
			fmt.Print("> ")
			fmt.Fscanln(os.Stdin, &label)

			fmt.Println("Введите текст для сохранения:")
			fmt.Print("> ")
			// Используем bufio.Scanner для чтения многострочного текста
			var buffer string
			fmt.Fscanln(os.Stdin, &buffer)
			content = buffer

			// Создаем структуру TextData
			textData := &domain.TextData{
				Content: content,
			}

			// Вызываем метод сохранения текста
			err := c.clientUseCase.SaveText(label, textData)
			if err != nil {
				fmt.Println("Ошибка при сохранении текста:", err)
				return
			}

			fmt.Println("Текст успешно сохранен!")
		},
	}
}

// GetTextCmd создает команду для получения текстовых данных
func (c *Command) GetTextCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get-text",
		Short: "Получение текстовых данных",
		Run: func(cmd *cobra.Command, args []string) {
			var label string

			fmt.Println("Введите уникальное название (label) текстовых данных для получения:")
			fmt.Print("> ")
			fmt.Fscanln(os.Stdin, &label)

			// Вызываем метод получения текста
			textData, err := c.clientUseCase.GetText(label)
			if err != nil {
				fmt.Println("Ошибка при получении текста:", err)
				return
			}

			fmt.Println("\nПолученный текст:")
			fmt.Println("------------------")
			fmt.Println(textData.Content)
			fmt.Println("------------------")
		},
	}
}

// DeleteTextCmd создает команду для удаления текстовых данных
func (c *Command) DeleteTextCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete-text",
		Short: "Удаление текстовых данных",
		Run: func(cmd *cobra.Command, args []string) {
			var label string

			fmt.Println("Введите уникальное название (label) текстовых данных для удаления:")
			fmt.Print("> ")
			fmt.Fscanln(os.Stdin, &label)

			// Вызываем метод удаления текста
			err := c.clientUseCase.DeleteText(label)
			if err != nil {
				fmt.Println("Ошибка при удалении текста:", err)
				return
			}

			fmt.Println("Текст успешно удален!")
		},
	}
}

// SaveCardCmd создает команду для сохранения данных кредитной карты
func (c *Command) SaveCardCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "save-card",
		Short: "Сохранение данных кредитной карты",
		Run: func(cmd *cobra.Command, args []string) {
			var label, number, holder, expiryDate, cvv string

			fmt.Println("Введите уникальное название (label) для данных карты:")
			fmt.Print("> ")
			fmt.Fscanln(os.Stdin, &label)

			fmt.Println("Введите номер карты:")
			fmt.Print("> ")
			fmt.Fscanln(os.Stdin, &number)

			fmt.Println("Введите имя держателя карты:")
			fmt.Print("> ")
			fmt.Fscanln(os.Stdin, &holder)

			fmt.Println("Введите срок действия карты (MM/YY):")
			fmt.Print("> ")
			fmt.Fscanln(os.Stdin, &expiryDate)

			fmt.Println("Введите CVV код:")
			fmt.Print("> ")
			fmt.Fscanln(os.Stdin, &cvv)

			// Создаем структуру CardData
			cardData := &domain.CardData{
				Number:     number,
				Holder:     holder,
				ExpiryDate: expiryDate,
				CVV:        cvv,
			}

			// Вызываем метод сохранения данных карты
			err := c.clientUseCase.SaveCard(label, cardData)
			if err != nil {
				fmt.Println("Ошибка при сохранении данных карты:", err)
				return
			}

			fmt.Println("Данные карты успешно сохранены!")
		},
	}
}

// GetCardCmd создает команду для получения данных кредитной карты
func (c *Command) GetCardCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get-card",
		Short: "Получение данных кредитной карты",
		Run: func(cmd *cobra.Command, args []string) {
			var label string

			fmt.Println("Введите уникальное название (label) данных карты для получения:")
			fmt.Print("> ")
			fmt.Fscanln(os.Stdin, &label)

			// Вызываем метод получения данных карты
			cardData, err := c.clientUseCase.GetCard(label)
			if err != nil {
				fmt.Println("Ошибка при получении данных карты:", err)
				return
			}

			fmt.Println("\nДанные кредитной карты:")
			fmt.Println("------------------------")
			fmt.Println("Номер карты:", cardData.Number)
			fmt.Println("Держатель карты:", cardData.Holder)
			fmt.Println("Срок действия:", cardData.ExpiryDate)
			fmt.Println("CVV код:", cardData.CVV)
			fmt.Println("------------------------")
		},
	}
}

// DeleteCardCmd создает команду для удаления данных кредитной карты
func (c *Command) DeleteCardCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete-card",
		Short: "Удаление данных кредитной карты",
		Run: func(cmd *cobra.Command, args []string) {
			var label string

			fmt.Println("Введите уникальное название (label) данных карты для удаления:")
			fmt.Print("> ")
			fmt.Fscanln(os.Stdin, &label)

			// Вызываем метод удаления данных карты
			err := c.clientUseCase.DeleteCard(label)
			if err != nil {
				fmt.Println("Ошибка при удалении данных карты:", err)
				return
			}

			fmt.Println("Данные карты успешно удалены!")
		},
	}
}

// SaveCredentialCmd создает команду для сохранения учетных данных
func (c *Command) SaveCredentialCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "save-credential",
		Short: "Сохранение учетных данных (логин/пароль)",
		Run: func(cmd *cobra.Command, args []string) {
			var label, login, password string

			fmt.Println("Введите уникальное название (label) для учетных данных:")
			fmt.Print("> ")
			fmt.Fscanln(os.Stdin, &label)

			fmt.Println("Введите логин:")
			fmt.Print("> ")
			fmt.Fscanln(os.Stdin, &login)

			fmt.Println("Введите пароль:")
			fmt.Print("> ")
			fmt.Fscanln(os.Stdin, &password)

			// Создаем структуру CredentialData
			credentialData := &domain.CredentialData{
				Login:    login,
				Password: password,
			}

			// Вызываем метод сохранения учетных данных
			err := c.clientUseCase.SaveCredential(label, credentialData)
			if err != nil {
				fmt.Println("Ошибка при сохранении учетных данных:", err)
				return
			}

			fmt.Println("Учетные данные успешно сохранены!")
		},
	}
}

// GetCredentialCmd создает команду для получения учетных данных
func (c *Command) GetCredentialCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get-credential",
		Short: "Получение учетных данных (логин/пароль)",
		Run: func(cmd *cobra.Command, args []string) {
			var label string

			fmt.Println("Введите уникальное название (label) учетных данных для получения:")
			fmt.Print("> ")
			fmt.Fscanln(os.Stdin, &label)

			// Вызываем метод получения учетных данных
			credentialData, err := c.clientUseCase.GetCredential(label)
			if err != nil {
				fmt.Println("Ошибка при получении учетных данных:", err)
				return
			}

			fmt.Println("\nУчетные данные:")
			fmt.Println("---------------")
			fmt.Println("Логин:", credentialData.Login)
			fmt.Println("Пароль:", credentialData.Password)
			fmt.Println("---------------")
		},
	}
}

// DeleteCredentialCmd создает команду для удаления учетных данных
func (c *Command) DeleteCredentialCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete-credential",
		Short: "Удаление учетных данных (логин/пароль)",
		Run: func(cmd *cobra.Command, args []string) {
			var label string

			fmt.Println("Введите уникальное название (label) учетных данных для удаления:")
			fmt.Print("> ")
			fmt.Fscanln(os.Stdin, &label)

			// Вызываем метод удаления учетных данных
			err := c.clientUseCase.DeleteCredential(label)
			if err != nil {
				fmt.Println("Ошибка при удалении учетных данных:", err)
				return
			}

			fmt.Println("Учетные данные успешно удалены!")
		},
	}
}
