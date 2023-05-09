package usecase

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
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
		return fmt.Errorf("usecase.Save: %w", err)
	}
	return nil
}

func (uc *UseCase) Delete(chatID int64, service string) (err error) {
	service, err = uc.Hash(service)
	if err != nil {
		log.Println(fmt.Errorf("usecase.Hash: %w", err))
		return fmt.Errorf("usecase.Hash: %w", err)
	}
	if err := uc.storage.Delete(chatID, service); err != nil {
		log.Println(err)
		return fmt.Errorf("usecase.Delete: %w", err)
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
		log.Println(err) // TODO: USE LOGGER
	}
}

func (uc *UseCase) Encrypt(text string) string {
	if text == "" {
		return ""
	} else if len(text) < aes.BlockSize {
		text += strings.Repeat(" ", aes.BlockSize-len(text))
	}

	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(uc.cipher, bytes)
	ciphertext := make([]byte, len(plainText))
	cfb.XORKeyStream(ciphertext, plainText)

	return base64.RawStdEncoding.EncodeToString(ciphertext)
}

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func (uc *UseCase) Decrypt(text string) string {
	if text == "" || len(text) < aes.BlockSize {
		return text
	}

	ciphertext, err := base64.RawStdEncoding.DecodeString(text)
	if err != nil {
		log.Println(err)
		return ""
	}

	cfb := cipher.NewCFBDecrypter(uc.cipher, bytes)
	plainText := make([]byte, len(ciphertext))
	cfb.XORKeyStream(plainText, ciphertext)

	return strings.Trim(string(plainText), " ")
}

func (uc *UseCase) Hash(text string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(text))
	if err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(hash.Sum(nil)), nil
}
