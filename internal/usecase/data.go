package usecase

import (
	"encoding/json"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"github.com/SmirnovND/gophkeeper/pkg"
	"net/http"
)

type DataUseCase struct {
	dataService interfaces.DataService
}

func NewDataUseCase(
	dataService interfaces.DataService,
) interfaces.DataUseCase {
	return &DataUseCase{
		dataService: dataService,
	}
}
func (c *DataUseCase) GetCredential(w http.ResponseWriter, r *http.Request, label string) {
	login, err := pkg.ExtractLoginFromToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Ошибка получения логина: "+err.Error(), http.StatusInternalServerError)
		return
	}

	credentialData, err := c.dataService.GetCredential(login, label)
	if err != nil {
		if err.Error() == "учетные данные не найдены" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем данные в ответе
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(credentialData)
}

func (c *DataUseCase) SaveCredential(w http.ResponseWriter, r *http.Request, label string, credentialData *domain.CredentialData) {
	login, err := pkg.ExtractLoginFromToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Ошибка получения логина: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Сохраняем данные
	err = c.dataService.SaveCredential(login, label, credentialData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "учетные данные успешно сохранены"})
}

func (c *DataUseCase) SaveCard(w http.ResponseWriter, r *http.Request, label string, cardData *domain.CardData) {
	login, err := pkg.ExtractLoginFromToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Ошибка получения логина: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Сохраняем данные
	err = c.dataService.SaveCard(login, label, cardData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "данные карты успешно сохранены"})
}

func (c *DataUseCase) GetCard(w http.ResponseWriter, r *http.Request, label string) {
	login, err := pkg.ExtractLoginFromToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Ошибка получения логина: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем данные
	cardData, err := c.dataService.GetCard(login, label)
	if err != nil {
		if err.Error() == "данные карты не найдены" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем данные в ответе
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cardData)
}

func (c *DataUseCase) SaveText(w http.ResponseWriter, r *http.Request, label string, textData *domain.TextData) {
	login, err := pkg.ExtractLoginFromToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Ошибка получения логина: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Сохраняем данные
	err = c.dataService.SaveText(login, label, textData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "текстовые данные успешно сохранены"})
}

func (c *DataUseCase) GetText(w http.ResponseWriter, r *http.Request, label string) {
	login, err := pkg.ExtractLoginFromToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Ошибка получения логина: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем данные
	textData, err := c.dataService.GetText(login, label)
	if err != nil {
		if err.Error() == "текстовые данные не найдены" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем данные в ответе
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(textData)
}
