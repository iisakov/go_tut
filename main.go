package main

import (
	"flag"
	"log"

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
	token := flag.String("tg-bot-token", "", "Токен для использования telegram бота")

	flag.Parse()

	if *token == "" {
		log.Fatal("Нет токена")
	}
	return *token
}

// func mustHost() string {
// 	host := flag.String("host-bot-host", "", "Хост для использования telegram бота")

// 	flag.Parse()

// 	if *host == "" {
// 		log.Fatal("Нет хоста")
// 	}
// 	return *host
// }
