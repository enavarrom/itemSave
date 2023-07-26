package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type CategoryApiResultDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CategoryApiService struct {
	apiURL            string
	redisCacheService *RedisCacheService
}

func NewCategoryApiService(redisCacheService *RedisCacheService) *CategoryApiService {
	return &CategoryApiService{apiURL: "https://api.mercadolibre.com/categories/%s", redisCacheService: redisCacheService}
}

func (c CategoryApiService) Get(categoryId string) (CategoryApiResultDTO, error) {

	cacheKey := fmt.Sprintf("category:%s", categoryId)
	jsonData := c.redisCacheService.Get(cacheKey)

	if jsonData != nil {
		var cacheData CategoryApiResultDTO
		err := json.Unmarshal([]byte(jsonData.(string)), &cacheData)
		if err != nil {
			return CategoryApiResultDTO{}, err
		}
		return cacheData, nil
	}

	if categoryId == "" {
		return CategoryApiResultDTO{}, errors.New("categoría consultada es vacío")
	}

	//Construir api url
	resp, err := http.Get(fmt.Sprintf(c.apiURL, categoryId))
	if err != nil {
		return CategoryApiResultDTO{}, err
	}

	defer resp.Body.Close()
	// Verificar el código de estado de la respuesta
	if resp.StatusCode != http.StatusOK {
		return CategoryApiResultDTO{}, errors.New(fmt.Sprintf("Error en la solicitud. Código de estado: %d", resp.StatusCode))
	}

	//Decodificando el response
	var response CategoryApiResultDTO
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return CategoryApiResultDTO{}, errors.New(fmt.Sprintf("Error al decodificar el cuerpo de la respuesta JSON: %v", err))
	}

	jsonData, err = json.Marshal(response)
	if err != nil {
		log.Println(err)
	}
	c.redisCacheService.Set(cacheKey, jsonData)

	return response, nil

}
