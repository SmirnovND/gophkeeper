package interfaces

import "github.com/spf13/cobra"

type Command interface {
	Login() *cobra.Command
	RegisterCmd() *cobra.Command
	UploadCmd() *cobra.Command
	DownloadCmd() *cobra.Command
	
	// Команды для работы с текстовыми данными
	SaveTextCmd() *cobra.Command
	GetTextCmd() *cobra.Command
	DeleteTextCmd() *cobra.Command
	
	// Команды для работы с данными кредитных карт
	SaveCardCmd() *cobra.Command
	GetCardCmd() *cobra.Command
	DeleteCardCmd() *cobra.Command
	
	// Команды для работы с учетными данными
	SaveCredentialCmd() *cobra.Command
	GetCredentialCmd() *cobra.Command
	DeleteCredentialCmd() *cobra.Command
	
	// Команда для получения информации о версии
	VersionCmd() *cobra.Command
}
