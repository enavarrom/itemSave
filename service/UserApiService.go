package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type UserApiResultDTO struct {
	NickName string `json:"nickname"`
}

type UserApiService struct {
	apiUrl            string
	redisCacheService *RedisCacheService
}

func NewUserApiService(redisCacheService *RedisCacheService) *UserApiService {
	return &UserApiService{apiUrl: "https://api.mercadolibre.com/users/%d", redisCacheService: redisCacheService}
}

func (c UserApiService) Get(userId int) (UserApiResultDTO, error) {
	cacheKey := fmt.Sprintf("user:%d", userId)

	jsonData := c.redisCacheService.Get(cacheKey)

	if jsonData != nil {
		var cacheData UserApiResultDTO
		err := json.Unmarshal([]byte(jsonData.(string)), &cacheData)
		if err != nil {
			return UserApiResultDTO{}, err
		}
		return cacheData, nil
	}

	//Construir api url
	resp, err := http.Get(fmt.Sprintf(c.apiUrl, userId))
	if err != nil {
		return UserApiResultDTO{}, nil
	}

	defer resp.Body.Close()
	// Verificar el código de estado de la respuesta
	if resp.StatusCode != http.StatusOK {
		return UserApiResultDTO{}, errors.New(fmt.Sprintf("Error en la solicitud. Código de estado: %d", resp.StatusCode))
	}

	//Decodificando el response
	var response UserApiResultDTO
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return UserApiResultDTO{}, errors.New(fmt.Sprintf("Error al decodificar el cuerpo de la respuesta JSON: %v", err))
	}

	jsonData, err = json.Marshal(response)
	if err != nil {
		log.Println(err)
	}
	c.redisCacheService.Set(cacheKey, jsonData)

	return response, nil

}
