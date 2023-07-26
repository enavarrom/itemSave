package service

import (
	"fmt"
	"itemSave/model"
	"itemSave/repository"
	"log"
	"os"
	"strconv"
	"sync"
)

type ItemService struct {
	*ItemApiService
	*CategoryApiService
	*CurrencyApiService
	*UserApiService
	*repository.MongoRepository
}

type ItemStreamDTO struct {
	ID      string `json:"id"`
	Site    string `json:"site"`
	EventId string
}

func NewItemService() *ItemService {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		fmt.Println("Error Redis Port:", err)
		return nil
	}
	redisCacheService := NewRedisCacheService(redisHost, redisPort)

	return &ItemService{ItemApiService: NewItemApiService(), CategoryApiService: NewCategoryApiService(redisCacheService),
		CurrencyApiService: NewCurrencyApiService(redisCacheService), UserApiService: NewUserApiService(redisCacheService),
		MongoRepository: repository.NewMongoRepository(),
	}
}

func (c *ItemService) ProcessItemStream(itemsStream []ItemStreamDTO) {
	var keys []string
	var items []model.Item

	for _, item := range itemsStream {
		keys = append(keys, item.Site+item.ID)
	}

	itemApiResponses, err := c.ItemApiService.Get(keys)
	if err != nil {
		fmt.Println(err)
	}

	for _, dto := range itemsStream {
		for _, itemApiResponse := range itemApiResponses {
			if dto.Site+dto.ID == itemApiResponse.ID {
				var item model.Item
				item.Id = dto.ID
				item.Site = dto.Site
				if itemApiResponse.ErrorMessage == "" {
					item.Price = itemApiResponse.Price
					item.StartTime = itemApiResponse.StartTime
					c.addOtherInfo(&item, &itemApiResponse)
				} else {
					item.Error = itemApiResponse.ErrorMessage
				}
				items = append(items, item)
			}
		}
	}
	c.save(items)
}

func (c *ItemService) save(items []model.Item) {
	var documents []interface{}
	for _, item := range items {
		documents = append(documents, item)
	}
	c.MongoRepository.InsertMany(documents, "items")
}

func (c *ItemService) getItemById(id string, itemsStream []ItemStreamDTO) *ItemStreamDTO {
	for _, item := range itemsStream {
		if item.Site+item.ID == id {
			return &item
		}
	}
	return nil
}

func (c *ItemService) addOtherInfo(item *model.Item, itemApiResponse *ItemApiResultDTO) {
	// Crear un canal para almacenar los resultados de las consultas en paralelo
	resultsChan := make(chan interface{}, 3)

	// Crear un wait group para esperar a que todas las goroutines terminen
	var wg sync.WaitGroup

	// Consultar category en paralelo
	wg.Add(1)
	go func(categoryId string) {
		defer wg.Done()
		category, err := c.CategoryApiService.Get(categoryId)
		if err != nil {
			log.Println(err)
			resultsChan <- nil
			return
		}
		resultsChan <- category
	}(itemApiResponse.CategoryId)

	// Consultar currency en paralelo
	wg.Add(1)
	go func(currencyId string) {
		defer wg.Done()
		currency, err := c.CurrencyApiService.Get(currencyId)
		if err != nil {
			log.Println(err)
			resultsChan <- nil
			return
		}
		resultsChan <- currency
	}(itemApiResponse.CurrencyId)

	// Consultar user en paralelo
	wg.Add(1)
	go func(userId int) {
		defer wg.Done()
		user, err := c.UserApiService.Get(userId)
		if err != nil {
			log.Println(err)
			resultsChan <- nil
			return
		}
		resultsChan <- user
	}(itemApiResponse.SellerId)

	// Esperar a que todas las goroutines terminen
	wg.Wait()

	// Cerrar el canal después de que todas las goroutines hayan terminado
	close(resultsChan)

	// Recopilar los resultados del canal y setearlo a item
	for result := range resultsChan {
		switch result.(type) {
		case CategoryApiResultDTO:
			category := result.(CategoryApiResultDTO)
			item.Name = category.Name
		case CurrencyApiResultDTO:
			currency := result.(CurrencyApiResultDTO)
			item.Description = currency.Description
		case UserApiResultDTO:
			user := result.(UserApiResultDTO)
			item.Nickname = user.NickName
		default:
			log.Println("Resultado no válido")
		}
	}

}
