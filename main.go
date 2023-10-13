package main

import (
	"log"
	"os"

	tgClient "example/hello/clients/telegram"
	event_consumer "example/hello/events/consumer/event-consumer"
	"example/hello/events/telegram"
	"example/hello/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "file_storage"
	batchSize   = 100
)

func main() {
	tgClient := tgClient.New(tgBotHost, mustToken())
	eventProcessor := telegram.New(tgClient, files.New(storagePath))

	log.Print("Сервис запущен")

	consumer := event_consumer.New(eventProcessor, eventProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("Сервис остановился", err)
	}
}

func mustToken() string {
	token, ok := os.LookupEnv("TOKEN")
	if !ok {
		log.Fatal("Нет токена")
	}

	return token
}
