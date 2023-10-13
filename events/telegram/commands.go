package telegram

import (
	"errors"
	errl "example/hello/lib"
	"example/hello/storage"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("Получил новую команду '%s' от '%s'", text, username)

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sandRndPage(chatID, username)
	case HelpCmd:
		return p.sandHelp(chatID)
	case StartCmd:
		return p.sandStart(chatID)
	default:
		// return p.tg.SendMessage(chatID, msgUnknownCommand)
		return nil
	}
}

func (p *Processor) savePage(chatID int, pageURL, username string) (err error) {
	defer func() { err = errl.WrapIfErr("Не удалось сохранить страницу", err) }()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	IsExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}

	if IsExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaveCommand); err != nil {
		return err
	}

	return nil

}

func (p *Processor) sandRndPage(chatID int, username string) (err error) {
	defer func() {
		err = errl.WrapIfErr("Не удалось выдать случайную страницу", err)
	}()

	page, err := p.storage.PickRandom(username)

	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavePage)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func (p *Processor) sandHelp(chatID int) (err error) {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sandStart(chatID int) (err error) {
	return p.tg.SendMessage(chatID, msgStart)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
