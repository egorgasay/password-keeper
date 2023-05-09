package usecase

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"password-keeper/internal/entity"
	"password-keeper/internal/storage"
	"strings"
)

type UseCase struct {
	storage *storage.Storage
	cipher  cipher.Block
}

const defaultLanguage = "en"

func New(storage *storage.Storage, key string) (*UseCase, error) {
	cipher, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	return &UseCase{
		storage: storage,
		cipher:  cipher,
	}, nil
}

func (uc *UseCase) Get(chatID int64, service string) (entity.Pair, error) {
	service, err := uc.Hash(service)
	if err != nil {
		log.Println(fmt.Errorf("usecase.Hash: %w", err))
		return entity.Pair{}, fmt.Errorf("usecase.Hash: %w", err)
	}

	pair, err := uc.storage.Get(chatID, service)
	if err != nil {
		log.Println(err)
		return entity.Pair{}, err
	}
	pair.Login = uc.Decrypt(pair.Login)
	pair.Password = uc.Decrypt(pair.Password)
	return pair, nil
}

func (uc *UseCase) Save(chatID int64, service, login, password string) error {
	login, password = uc.Encrypt(login), uc.Encrypt(password)
	service, err := uc.Hash(service)
	if err != nil {
		log.Println(fmt.Errorf("usecase.Hash: %w", err))
		return fmt.Errorf("usecase.Hash: %w", err)
	}

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

func (uc *UseCase) Encrypt(text string) string {
	if text == "" {
		return ""
	} else if len(text) < aes.BlockSize {
		text += strings.Repeat(" ", aes.BlockSize-len(text))
	}

	ciphertext := make([]byte, len(text))
	uc.cipher.Encrypt(ciphertext, []byte(text))
	return string(ciphertext)
}

func (uc *UseCase) Decrypt(text string) string {
	if text == "" {
		return ""
	} else if len(text) < aes.BlockSize {
		return ""
	}

	ciphertext := []byte(text)
	plaintext := make([]byte, len(text))
	uc.cipher.Decrypt(plaintext, ciphertext)
	return strings.TrimSpace(string(plaintext))
}

func (uc *UseCase) Hash(text string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(text))
	if err != nil {
		return "", err
	}

	return string(hash.Sum(nil)), nil
}
