package telegram

import "example/hello/clients/telegram"

type Processor struct {
	tg     *telegram.Client
	offset int
}
