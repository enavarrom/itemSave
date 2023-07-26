package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type CurrencyApiResultDTO struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type CurrencyApiService struct {
	apiUrl            string
	redisCacheService *RedisCacheService
}

func NewCurrencyApiService(redisCacheService *RedisCacheService) *CurrencyApiService {
	return &CurrencyApiService{apiUrl: "https://api.mercadolibre.com/currencies/%s", redisCacheService: redisCacheService}
}

func (c CurrencyApiService) Get(currencyId string) (CurrencyApiResultDTO, error) {
	//Construir api url

	cacheKey := fmt.Sprintf("currency:%s", currencyId)
	jsonData := c.redisCacheService.Get(cacheKey)

	if jsonData != nil {
		var cacheData CurrencyApiResultDTO
		err := json.Unmarshal([]byte(jsonData.(string)), &cacheData)
		if err != nil {
			return CurrencyApiResultDTO{}, err
		}
		return cacheData, nil
	}

	if currencyId == "" {
		return CurrencyApiResultDTO{}, errors.New("currency consultada es vacío")
	}

	resp, err := http.Get(fmt.Sprintf(c.apiUrl, currencyId))
	if err != nil {
		return CurrencyApiResultDTO{}, err
	}

	defer resp.Body.Close()
	// Verificar el código de estado de la respuesta
	if resp.StatusCode != http.StatusOK {
		return CurrencyApiResultDTO{}, errors.New(fmt.Sprintf("Error en la solicitud. Código de estado: %d", resp.StatusCode))
	}

	//Decodificando el response
	var response CurrencyApiResultDTO
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return CurrencyApiResultDTO{}, errors.New(fmt.Sprintf("Error al decodificar el cuerpo de la respuesta JSON: %v", err))
	}

	jsonData, err = json.Marshal(response)
	if err != nil {
		log.Println(err)
	}
	c.redisCacheService.Set(cacheKey, jsonData)

	return response, nil
}
