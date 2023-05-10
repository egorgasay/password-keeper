package usecase

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"password-keeper/internal/entity"
	"password-keeper/internal/storage"
	"strings"
)

// UseCase is the main struct for the application logic.
type UseCase struct {
	storage *storage.Storage
	cipher  cipher.Block
}

const defaultLanguage = "en"

// New creates a new UseCase.
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

// Get returns the pair from the storage.
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

// Save saves the pair to the storage.
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

// Delete deletes the pair from the storage.
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

// GetLang returns the language of the user.
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

// SetLang sets the language of the user.
func (uc *UseCase) SetLang(chatID int64, lang string) {
	err := uc.storage.SetLang(chatID, lang)
	if err != nil {
		log.Println(err) // TODO: USE LOGGER
	}
}

// Encrypt encrypts the text.
func (uc *UseCase) Encrypt(text string) string {
	if text == "" {
		return ""
	} else if len(text) < aes.BlockSize {
		text += strings.Repeat(" ", aes.BlockSize-len(text))
	}

	plainText := []byte(text)
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Println(err)
		return ""
	}

	stream := cipher.NewCFBEncrypter(uc.cipher, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return base64.RawStdEncoding.EncodeToString(cipherText)
}

// Decrypt decrypts the text.
func (uc *UseCase) Decrypt(text string) string {
	if text == "" || len(text) < aes.BlockSize {
		return text
	}

	ciphertext, err := base64.RawStdEncoding.DecodeString(text)
	if err != nil {
		log.Println(err)
		return ""
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(uc.cipher, iv)
	cfb.XORKeyStream(ciphertext, ciphertext)

	return strings.Trim(string(ciphertext), " ")
}

// Hash hashes the text.
func (uc *UseCase) Hash(text string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(text))
	if err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(hash.Sum(nil)), nil
}
