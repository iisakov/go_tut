package storage

import (
	"crypto/sha1"
	"errors"
	errl "example/hello/lib"
	"fmt"
	"io"
)

var ErrNoSavedPages = errors.New("нет сохранённых страниц")

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", errl.Wrap("не вышло посчитать hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", errl.Wrap("не вышло посчитать hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
