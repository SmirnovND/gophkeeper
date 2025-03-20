package command

import "github.com/SmirnovND/gophkeeper/internal/interfaces"

type Command struct {
	clientUseCase interfaces.ClientUseCase
}

func NewCommand(
	ClientUseCase interfaces.ClientUseCase,
) *Command {
	return &Command{
		clientUseCase: ClientUseCase,
	}
}
