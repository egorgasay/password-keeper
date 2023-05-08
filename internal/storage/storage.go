package storage

import (
	"errors"
	"log"
	"password-keeper/internal/entity"
	"sync"
)

type RealStorage interface {
	Save(chatID int64, service string, pair entity.Pair) error
	Get(chatID int64, service string) (entity.Pair, error)
	Delete(chatID int64, service string) error
	GetLang(chatID int64) (string, error)
}

type Storage struct {
	ramStorage  *sync.Map
	realStorage RealStorage
}

var ErrNotFound = errors.New("not found")

func New() (*Storage, error) {
	return &Storage{
		ramStorage: &sync.Map{},
	}, nil
}

func (s *Storage) Save(chatID int64, service string, pair entity.Pair) error {
	us, err := s.getUserStorage(chatID)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}

	us.Store(service, pair)
	return nil
}

func (s *Storage) getUserStorage(chatID int64) (*sync.Map, error) {
	us, _ := s.ramStorage.LoadOrStore(chatID, &sync.Map{})

	storage, ok := us.(*sync.Map)
	if !ok {
		log.Println("storage is not *sync.Map")
		return nil, ErrNotFound
	}

	return storage, nil
}

func (s *Storage) Get(chatID int64, service string) (entity.Pair, error) {
	us, err := s.getUserStorage(chatID)
	if err != nil {
		return entity.Pair{}, err
	}

	value, ok := us.Load(service)
	if !ok {
		return entity.Pair{}, ErrNotFound
	}

	pair, ok := value.(entity.Pair)
	if !ok {
		return entity.Pair{}, ErrNotFound
	}

	return pair, nil
}

func (s *Storage) Delete(chatID int64, service string) error {
	us, err := s.getUserStorage(chatID)
	if err != nil {
		return err
	}

	us.Delete(service)
	return nil
}
