package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func (c *Command) Login() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Авторизация в сервисе",
		Run: func(cmd *cobra.Command, args []string) {
			var username, password string
			
			fmt.Print("Введите логин: ")
			fmt.Fscanln(os.Stdin, &username)
			
			fmt.Print("Введите пароль: ")
			fmt.Fscanln(os.Stdin, &password)
			
			err := c.clientUseCase.Login(username, password)
			if err != nil {
				fmt.Println("Ошибка авторизации:", err)
				return
			}
			
			fmt.Println("Успешная авторизация!")
		},
	}
}

func (c *Command) RegisterCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "register",
		Short: "Регистрация в сервисе",
		Run: func(cmd *cobra.Command, args []string) {
			var username, password, passwordCheck string
			
			fmt.Print("Введите логин: ")
			fmt.Fscanln(os.Stdin, &username)
			
			fmt.Print("Введите пароль: ")
			fmt.Fscanln(os.Stdin, &password)
			
			fmt.Print("Введите пароль еще раз: ")
			fmt.Fscanln(os.Stdin, &passwordCheck)
			
			err := c.clientUseCase.Register(username, password, passwordCheck)
			if err != nil {
				fmt.Println("Ошибка регистрации:", err)
				return
			}
			
			fmt.Println("Успешная регистрация!")
		},
	}
}
