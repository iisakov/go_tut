package files

import (
	"encoding/gob"
	"errors"
	errl "example/hello/lib"
	"example/hello/storage"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

const defaultPerm = 0774

var ErrNoSavedPage = errors.New("нет сохранённых страниц")

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = errl.WrapIfErr("Не удалось сохранить.", err) }()

	fPath := filepath.Join(s.basePath, page.UserName)
	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}

	file.Close()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = errl.WrapIfErr("Не удалось выдать странеицу.", err) }()

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, ErrNoSavedPage
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n].Name()

	return s.decodePage(file)
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errl.Wrap("Не удалось декодировать файл", err)
	}
	defer f.Close()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, errl.Wrap("Не удалось декодировать файл", err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
