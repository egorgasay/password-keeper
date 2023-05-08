package usecase

import (
	"password-keeper/internal/entity"
	"password-keeper/internal/storage"
)

type UseCase struct {
	storage *storage.Storage
}

func New(storage *storage.Storage) *UseCase {
	return &UseCase{
		storage: storage,
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
