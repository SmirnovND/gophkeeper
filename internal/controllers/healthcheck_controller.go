package controllers

import (
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"net/http"
)

type HealthcheckController struct {
	DB interfaces.DB
}

func NewHealthcheckController(DB interfaces.DB) *HealthcheckController {
	return &HealthcheckController{
		DB: DB,
	}
}

func (hc *HealthcheckController) HandlePing(w http.ResponseWriter, r *http.Request) {
	err := hc.DB.Ping()
	if err != nil {
		http.Error(w, "Failed to connect DB", http.StatusInternalServerError)
		return
	} else {
		w.Write([]byte("pong"))
		w.WriteHeader(http.StatusOK)
	}
}
