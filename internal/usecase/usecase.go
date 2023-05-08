package usecase

import (
	"log"
	"password-keeper/internal/entity"
	"password-keeper/internal/storage"
	"sync"
)

type UseCase struct {
	storage     *storage.Storage
	key         string
	langStorage *sync.Map
}

func New(storage *storage.Storage, key string) *UseCase {
	return &UseCase{
		storage:     storage,
		key:         key,
		langStorage: &sync.Map{},
	}
}

func (uc *UseCase) Get(chatID int64, service string) (entity.Pair, error) {
	return uc.storage.Get(chatID, service)
}

func (uc *UseCase) Save(chatID int64, service, login, password string) error {
	return uc.storage.Save(chatID, service, entity.Pair{Login: login, Password: password})
}

func (uc *UseCase) Delete(chatID int64, service string) error {
	return uc.storage.Delete(chatID, service)
}

func (uc *UseCase) GetLang(chatID int64) string {
	lang, _ := uc.langStorage.LoadOrStore(chatID, "en")

	l, ok := lang.(string)
	if !ok {
		log.Println("lang is not string")
		return "en"
	}

	return l
}

func (uc *UseCase) SetLang(chatID int64, lang string) {
	uc.langStorage.Store(chatID, lang)
}
