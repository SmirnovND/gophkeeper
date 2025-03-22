package interfaces

import "github.com/spf13/cobra"

type Command interface {
	Login() *cobra.Command
	UploadCmd() *cobra.Command
}
