package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"itemSave/service"
	"log"
	"time"
)

type RedisItemStreamConsumer struct {
	client       *redis.Client
	context      context.Context
	streamName   string
	streamGroup  string
	consumerName string
	itemService  *service.ItemService
}

func NewRedisStreamConsumer(redisHost string, redisPort int, streamName string, streamGroup string, itemService *service.ItemService) *RedisItemStreamConsumer {
	return &RedisItemStreamConsumer{
		itemService: itemService,
		client: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", redisHost, redisPort),
		}),
		context:      context.Background(),
		streamName:   streamName,
		streamGroup:  streamGroup,
		consumerName: uuid.NewString(),
	}
}

func (c *RedisItemStreamConsumer) Consume() {
	c.createGroupConsumer()
	c.read()
}

func (c *RedisItemStreamConsumer) processItems(itemsStream []service.ItemStreamDTO) {
	c.itemService.ProcessItemStream(itemsStream)
	//Logica para identificar y actualizar los eventos de los items procesados correctamente
	for _, item := range itemsStream {
		// Después de procesar el mensaje, confirma que se ha procesado
		c.client.XAck(c.context, c.streamName, c.streamGroup, item.EventId)
	}
}

func (c *RedisItemStreamConsumer) createGroupConsumer() {
	// Crear el grupo de consumidores y un consumidor
	_, err := c.client.XGroupCreateMkStream(c.context, c.streamName, c.streamGroup, "$").Result()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		log.Fatalf("Error al unirse o crear el grupo de consumidores: %v", err)
	}
}

func (c *RedisItemStreamConsumer) read() {
	// Leer los mensajes desde el grupo
	for {
		messages, err := c.client.XReadGroup(c.context, &redis.XReadGroupArgs{
			Group:    c.streamGroup,
			Consumer: c.consumerName,
			Streams:  []string{c.streamName, ">"},
			Block:    time.Second, // Bloquear hasta que haya mensajes disponibles
			Count:    20,          // Procesar solo un mensaje a la vez
		}).Result()

		if err != nil {
			if err != redis.Nil {
				fmt.Printf("Error al leer mensajes del grupo: %v\n", err)
			}
		} else {
			var items []service.ItemStreamDTO
			// Procesar los mensajes
			for _, message := range messages {
				for _, event := range message.Messages {
					fmt.Printf("Mensaje: ID=%s Datos=%v\n", event.ID, event.Values)
					// Aquí procesas el mensaje y haces lo que necesites con él

					var item service.ItemStreamDTO
					err := json.Unmarshal([]byte(event.Values["message"].(string)), &item)
					if err != nil {
						fmt.Println("Error al deserializar el JSON:", err)
					}
					item.EventId = event.ID
					items = append(items, item)
				}
			}
			c.processItems(items)
		}
	}
}
