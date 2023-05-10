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
	"go.uber.org/zap"
	"io"
	"password-keeper/internal/entity"
	"password-keeper/internal/storage"
	"strings"
)

// UseCase is the main struct for the application logic.
type UseCase struct {
	storage *storage.Storage
	cipher  cipher.Block
	logger  *zap.Logger
}

const defaultLanguage = "en"

// New creates a new UseCase.
func New(storage *storage.Storage, key string, logger *zap.Logger) (*UseCase, error) {
	cipher, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	return &UseCase{
		storage: storage,
		cipher:  cipher,
		logger:  logger,
	}, nil
}

// Get returns the pair from the storage.
func (uc *UseCase) Get(chatID int64, service string) (entity.Pair, error) {
	service, err := uc.Hash(service)
	if err != nil {
		err = fmt.Errorf("usecase.Hash: %w", err)
		uc.logger.Warn(err.Error())
		return entity.Pair{}, err
	}

	pair, err := uc.storage.Get(chatID, service)
	if err != nil {
		err = fmt.Errorf("usecase.Get: %w", err)
		uc.logger.Warn(err.Error())
		return entity.Pair{}, err
	}
	pair.Login, err = uc.Decrypt(pair.Login)
	if err != nil {
		err = fmt.Errorf("usecase.Decrypt: %w", err)
		uc.logger.Warn(err.Error())
		return entity.Pair{}, err
	}

	pair.Password, err = uc.Decrypt(pair.Password)
	if err != nil {
		err = fmt.Errorf("usecase.Decrypt: %w", err)
		uc.logger.Warn(err.Error())
		return entity.Pair{}, err
	}

	return pair, nil
}

// Save saves the pair to the storage.
func (uc *UseCase) Save(chatID int64, service, login, password string) (err error) {
	login, err = uc.Encrypt(login)
	if err != nil {
		err = fmt.Errorf("usecase.Encrypt: %w", err)
		uc.logger.Warn(err.Error())
		return err
	}

	password, err = uc.Encrypt(password)
	if err != nil {
		err = fmt.Errorf("usecase.Encrypt: %w", err)
		uc.logger.Warn(err.Error())
		return err
	}

	service, err = uc.Hash(service)
	if err != nil {
		err = fmt.Errorf("usecase.Hash: %w", err)
		uc.logger.Warn(err.Error())
		return err
	}

	if err := uc.storage.Save(chatID, service, entity.Pair{Login: login, Password: password}); err != nil {
		err = fmt.Errorf("usecase.Save: %w", err)
		uc.logger.Warn(err.Error())
		return err
	}

	return nil
}

// Delete deletes the pair from the storage.
func (uc *UseCase) Delete(chatID int64, service string) (err error) {
	service, err = uc.Hash(service)
	if err != nil {
		err = fmt.Errorf("usecase.Hash: %w", err)
		uc.logger.Warn(err.Error())
		return err
	}
	if err := uc.storage.Delete(chatID, service); err != nil {
		err = fmt.Errorf("usecase.Delete: %w", err)
		uc.logger.Warn(err.Error())
		return err
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
		err = fmt.Errorf("usecase.GetLang: %w", err)
		uc.logger.Warn(err.Error())
		return defaultLanguage
	}
	return l
}

// SetLang sets the language of the user.
func (uc *UseCase) SetLang(chatID int64, lang string) {
	err := uc.storage.SetLang(chatID, lang)
	if err != nil {
		err = fmt.Errorf("usecase.SetLang: %w", err)
		uc.logger.Warn(err.Error())
	}
}

// Encrypt encrypts the text.
func (uc *UseCase) Encrypt(text string) (string, error) {
	if text == "" {
		return "", nil
	} else if len(text) < aes.BlockSize {
		text += strings.Repeat(" ", aes.BlockSize-len(text))
	}

	plainText := []byte(text)
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		err = fmt.Errorf("io.ReadFull: %w", err)
		uc.logger.Warn(err.Error())
		return "", err
	}

	stream := cipher.NewCFBEncrypter(uc.cipher, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return base64.RawStdEncoding.EncodeToString(cipherText), nil
}

// Decrypt decrypts the text.
func (uc *UseCase) Decrypt(text string) (string, error) {
	if text == "" || len(text) < aes.BlockSize {
		return text, nil
	}

	ciphertext, err := base64.RawStdEncoding.DecodeString(text)
	if err != nil {
		err = fmt.Errorf("base64.RawStdEncoding.DecodeString: %w", err)
		uc.logger.Warn(err.Error())
		return "", err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(uc.cipher, iv)
	cfb.XORKeyStream(ciphertext, ciphertext)

	return strings.Trim(string(ciphertext), " "), nil
}

// Hash hashes the text.
func (uc *UseCase) Hash(text string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(text))
	if err != nil {
		err = fmt.Errorf("hash.Write: %w", err)
		uc.logger.Warn(err.Error())
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(hash.Sum(nil)), nil
}
