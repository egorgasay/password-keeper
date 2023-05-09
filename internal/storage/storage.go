package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"password-keeper/internal/entity"
	"password-keeper/internal/storage/postgres"
	"password-keeper/internal/storage/queries"
	"password-keeper/internal/storage/service"
	"password-keeper/internal/storage/sqlite"
	"sync"
)

type RealStorage interface {
	Save(chatID int64, service string, pair entity.Pair) error
	Get(chatID int64, service string) (entity.Pair, error)
	Delete(chatID int64, service string) error
	GetLang(chatID int64) (string, error)
	SetLang(chatID int64, lang string) error
}

type Storage struct {
	ramStorage  *sync.Map
	realStorage RealStorage
	langStorage *sync.Map
}

var ErrNotFound = errors.New("not found")

func New(storageType, dsn string) (*Storage, error) {
	var rs RealStorage

	switch storageType {
	case "postgres":
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			return nil, fmt.Errorf("open db: %w", err)
		}

		rs, err = postgres.New(db, "file://migrations/postgres")
		if err != nil {
			return nil, fmt.Errorf("new postgres: %w", err)
		}

		err = queries.Prepare(db, "postgres")
		if err != nil {
			return nil, fmt.Errorf("prepare db: %w", err)
		}
	case "sqlite", "test":
		db, err := sql.Open("sqlite", dsn)
		if err != nil {
			return nil, fmt.Errorf("open db: %w", err)
		}

		if storageType == "test" {
			rs, err = sqlite.New(db, "file://../../migrations/sqlite")
			if err != nil {
				return nil, fmt.Errorf("new sqlite: %w", err)
			}
		} else {
			rs, err = sqlite.New(db, "file://migrations/sqlite")
			if err != nil {
				return nil, fmt.Errorf("new sqlite: %w", err)
			}
		}

		err = queries.Prepare(db, "sqlite")
		if err != nil {
			return nil, fmt.Errorf("prepare db: %w", err)
		}
	default:
		return nil, fmt.Errorf("unknown storage type: %s", storageType)
	}
	return &Storage{
		ramStorage:  &sync.Map{},
		langStorage: &sync.Map{},
		realStorage: rs,
	}, nil
}

func (s *Storage) Save(chatID int64, service string, pair entity.Pair) error {
	us, err := s.getUserStorage(chatID)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}

	us.Store(service, pair)
	return s.realStorage.Save(chatID, service, pair)
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
		if !errors.Is(err, ErrNotFound) {
			return entity.Pair{}, err
		}

		p, err := s.realStorage.Get(chatID, service)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return entity.Pair{}, ErrNotFound
			}
			return entity.Pair{}, err
		}
		return p, nil
	}

	value, ok := us.Load(service)
	if !ok {
		p, err := s.realStorage.Get(chatID, service)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return entity.Pair{}, ErrNotFound
			}
			return entity.Pair{}, err
		}
		return p, nil
	}

	pair, ok := value.(entity.Pair)
	if !ok {
		return entity.Pair{}, ErrNotFound
	}

	return pair, nil
}

func (s *Storage) Delete(chatID int64, serviceName string) error {
	us, err := s.getUserStorage(chatID)
	if err != nil {
		return err
	}

	us.Delete(serviceName)
	err = s.realStorage.Delete(chatID, serviceName)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("realStorage delete: %w", err)
	}
	return nil
}

func (s *Storage) GetLang(chatID int64) (string, error) {
	lang, loaded := s.langStorage.LoadOrStore(chatID, "en")
	if !loaded {
		lang, err := s.realStorage.GetLang(chatID)
		if err != nil {
			return "", fmt.Errorf("get lang: %w", err)
		}
		s.langStorage.Store(chatID, lang)
		return lang, nil
	}

	l, ok := lang.(string)
	if !ok {
		return "", ErrNotFound
	}

	return l, nil
}

func (s *Storage) SetLang(chatID int64, lang string) error {
	s.langStorage.Store(chatID, lang)
	err := s.realStorage.SetLang(chatID, lang)
	if err != nil {
		return fmt.Errorf("set lang: %w", err)
	}
	return nil
}
