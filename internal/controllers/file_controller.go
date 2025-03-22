package controllers

import (
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"github.com/SmirnovND/toolbox/pkg/paramsparser"
	"net/http"
)

type FileController struct {
	FileUseCase interfaces.CloudUseCase
}

func NewFileController(FileUseCase interfaces.CloudUseCase) *FileController {
	return &FileController{
		FileUseCase: FileUseCase,
	}
}

func (f *FileController) HandleUploadFile(w http.ResponseWriter, r *http.Request) {
	fileData, err := paramsparser.JSONParse[domain.FileData](w, r)
	if err != nil {
		return
	}
	f.FileUseCase.GenerateUploadLink(w, fileData)
}
