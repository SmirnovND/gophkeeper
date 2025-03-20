package command

import (
	"fmt"
	"github.com/spf13/cobra"
)

func (c *Command) Login() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Авторизация в сервисе",
		Run: func(cmd *cobra.Command, args []string) {
			var username, password string
			fmt.Print("Введите логин: ")
			fmt.Scanln(&username)

			fmt.Print("Введите пароль: ")
			fmt.Scanln(&password)

			err := c.clientUseCase.Login(username, password)
			if err != nil {
				fmt.Println("Ошибка авторизации:", err)
				return
			}

			fmt.Println("Успешная авторизация!")
		},
	}
}
