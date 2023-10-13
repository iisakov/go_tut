package files

import (
	"encoding/gob"
	"errors"
	errl "example/hello/lib"
	"example/hello/storage"
	"fmt"
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

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	file.Close()

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = errl.WrapIfErr("Не удалось выдать страницу.", err) }()

	path := filepath.Join(s.basePath, userName)

	//TODO Проверять, есть ли папка в ФЛ, Если нет - выдавать сообщение, у вас ещё ничего не сохранено

	if isDirExists, _ := isDirExists(path); !isDirExists {
		return nil, storage.ErrNoSavedPages
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n].Name()
	return s.decodePage(filepath.Join(path, file))
}

func isDirExists(dirPath string) (bool, error) {
	_, err := os.Stat(dirPath)
	if err != nil && os.IsNotExist(err) {
		return false, err
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errl.Wrap(fmt.Sprintf("Не удалось декодировать файл: %s", filePath), err)
	}
	defer f.Close()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, errl.Wrap(fmt.Sprintf("Не удалось декодировать файл: %s", filePath), err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, errl.Wrap("can't check if page exists", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", path)

		return false, errl.Wrap(msg, err)
	}

	return true, nil
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return errl.Wrap("can't remove page", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove page %s", path)

		return errl.Wrap(msg, err)
	}

	return nil
}
