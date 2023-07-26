package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ItemApiResultDTO struct {
	ID           string    `json:"id"`
	StartTime    time.Time `json:"start_time"`
	Price        float64   `json:"price"`
	CategoryId   string    `json:"category_id"`
	CurrencyId   string    `json:"currency_id"`
	SellerId     int       `json:"seller_id"`
	ErrorMessage string    `json:"message"`
}

type ItemResponse struct {
	Code             int `json:"code"`
	ItemApiResultDTO `json:"body"`
}

type ItemApiService struct {
	apiURL string
}

func NewItemApiService() *ItemApiService {
	return &ItemApiService{
		apiURL: "https://api.mercadolibre.com/items?ids=%s&attributes=%s",
	}
}

func (c ItemApiService) Get(ids []string) ([]ItemApiResultDTO, error) {
	//Atributos Requeridos
	attributes := []string{"id", "price", "start_time", "category_id", "currency_id", "seller_id"}
	//Construir api url
	resp, err := http.Get(fmt.Sprintf(c.apiURL, strings.Join(ids, ","), strings.Join(attributes, ",")))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// Verificar el código de estado de la respuesta
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Error en la solicitud. Código de estado: %d", resp.StatusCode))
	}

	//Decodificando el response
	var responses []ItemResponse
	err = json.NewDecoder(resp.Body).Decode(&responses)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error al decodificar el cuerpo de la respuesta JSON: %v", err))
	}

	return c.getResult(responses), err
}

func (ItemApiService) getResult(responses []ItemResponse) []ItemApiResultDTO {
	var result []ItemApiResultDTO
	for _, response := range responses {
		result = append(result, response.ItemApiResultDTO)
	}
	return result
}
