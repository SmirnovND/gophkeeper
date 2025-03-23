package interfaces

import "github.com/spf13/cobra"

type Command interface {
	Login() *cobra.Command
	RegisterCmd() *cobra.Command
	UploadCmd() *cobra.Command
	DownloadCmd() *cobra.Command
}
