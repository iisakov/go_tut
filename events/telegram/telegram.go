package telegram

import (
	"errors"
	"example/hello/clients/telegram"
	"example/hello/events"
	errl "example/hello/lib"
	"example/hello/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknowEventsType = errors.New("неизвестный тип события")
	ErrUnknowMetaType   = errors.New("неизвестный тип Меты")
)

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		offset:  0,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, errl.Wrap("Не удалось получить события", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	result := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		result = append(result, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return result, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.EventsType {
	case events.Message:
		return p.processMessage(event)
	default:
		return errl.Wrap("Неизвестный тип события", ErrUnknowEventsType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return errl.Wrap("Ошибка в получении Меты", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return errl.Wrap("Не удалось исполнить команду", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	result, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, errl.Wrap("Мета пустая.", ErrUnknowMetaType)
	}

	return result, nil
}

func event(upd telegram.Updates) events.Event {
	updType := fetchType(upd)
	result := events.Event{
		EventsType: updType,
		Text:       fetchText(upd),
	}

	if updType == events.Message {
		result.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return result
}

func fetchType(upd telegram.Updates) events.EventsType {
	if upd.Message != nil {
		return events.Message
	}

	return events.Unknown
}

func fetchText(upd telegram.Updates) string {
	if upd.Message != nil {
		return upd.Message.Text
	}

	return ""
}
