package command

import "github.com/SmirnovND/gophkeeper/internal/interfaces"

type Command struct {
	clientUseCase interfaces.ClientUseCase
}

func NewCommand(
	ClientUseCase interfaces.ClientUseCase,
) interfaces.Command {
	return &Command{
		clientUseCase: ClientUseCase,
	}
}
