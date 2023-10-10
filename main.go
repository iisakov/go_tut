package main

import (
	"example/hello/clients/telegram"
	"flag"
	"log"
)

func main() {
	tgClient = telegram.New(mustHost(), mustToken())
	// fetcher = fetcher.New(tgClient)
	// processor = processor.New(tgClient)
	// consumer.Start(fetcher, processor)
}

func mustToken() string {
	token := flag.String("token-bot-token", "", "Токен для использования telegram бота")

	flag.Parse()

	if *token == "" {
		log.Fatal("Нет токена")
	}
	return *token
}

func mustHost() string {
	host := flag.String("host-bot-host", "", "Хост для использования telegram бота")

	flag.Parse()

	if *host == "" {
		log.Fatal("Нет хоста")
	}
	return *host
}
