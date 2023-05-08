package usecase

import (
	"database/sql"
	"errors"
	"log"
	"password-keeper/internal/entity"
	"password-keeper/internal/storage"
)

type UseCase struct {
	storage *storage.Storage
	key     string
}

const defaultLanguage = "en"

func New(storage *storage.Storage, key string) *UseCase {

	return &UseCase{
		storage: storage,
		key:     key,
	}
}

func (uc *UseCase) Get(chatID int64, service string) (entity.Pair, error) {
	pair, err := uc.storage.Get(chatID, service)
	if err != nil {
		log.Println(err)
		return entity.Pair{}, err
	}
	return pair, nil
}

func (uc *UseCase) Save(chatID int64, service, login, password string) error {
	if err := uc.storage.Save(chatID, service, entity.Pair{Login: login, Password: password}); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (uc *UseCase) Delete(chatID int64, service string) error {
	if err := uc.storage.Delete(chatID, service); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (uc *UseCase) GetLang(chatID int64) string {
	l, err := uc.storage.GetLang(chatID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			uc.SetLang(chatID, defaultLanguage)
			return defaultLanguage
		}
		log.Println(err)
		return defaultLanguage
	}
	return l
}

func (uc *UseCase) SetLang(chatID int64, lang string) {
	err := uc.storage.SetLang(chatID, lang)
	if err != nil {
		log.Println(err)
	}
}
