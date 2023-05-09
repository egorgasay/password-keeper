package sqllike

import (
	"database/sql"
	"errors"
	"password-keeper/internal/entity"
	"password-keeper/internal/storage/queries"
	"password-keeper/internal/storage/service"
)

type DB struct {
	*sql.DB
}

func (db DB) Save(chatID int64, service string, pair entity.Pair) error {
	prep, err := queries.GetPreparedStatement(queries.AddService)
	if err != nil {
		return err
	}
	_, err = prep.Exec(service, pair.Login, pair.Password, chatID, pair.Login, pair.Password, chatID)
	return err
}

func (db DB) Get(chatID int64, service string) (entity.Pair, error) {
	prep, err := queries.GetPreparedStatement(queries.GetService)
	if err != nil {
		return entity.Pair{}, err
	}

	var pair entity.Pair
	err = prep.QueryRow(service, chatID).Scan(&pair.Login, &pair.Password)
	return pair, err
}

var ErrNotFound = errors.New("not found")

func (db DB) Delete(chatID int64, serviceName string) error {
	prep, err := queries.GetPreparedStatement(queries.DeleteService)
	if err != nil {
		return err
	}

	r, err := prep.Exec(serviceName, chatID)
	if err != nil {
		return err
	}
	a, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if a == 0 {
		return service.ErrNotFound
	}
	return nil
}

func (db DB) GetLang(chatID int64) (string, error) {
	prep, err := queries.GetPreparedStatement(queries.GetLang)
	if err != nil {
		return "", err
	}

	var lang string
	err = prep.QueryRow(chatID).Scan(&lang)
	return lang, err
}

func (db DB) SetLang(chatID int64, lang string) error {
	prep, err := queries.GetPreparedStatement(queries.AddOrUpdateChatLang)
	if err != nil {
		return err
	}
	_, err = prep.Exec(chatID, lang, lang)
	return err
}
