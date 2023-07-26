package main

import (
	"itemSave/consumer"
	"itemSave/service"
)

func main() {
	itemService := service.NewItemService()
	consumer.NewRedisStreamConsumer("localhost", 6378,
		"saveItems", "ItemApp", itemService).Consume()
}
