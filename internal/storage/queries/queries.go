package queries

import (
	"database/sql"
	"errors"
	"fmt"
)

// Query text of query.
type Query string

// Name number of query.
type Name int

// Query names.
const (
	AddService = iota
	AddOrUpdateChatLang
	GetService
	GetLang
	DeleteService
)

var queriesSqlite = map[Name]Query{
	AddService:          "INSERT INTO services (service, login, password, owner) VALUES (?, ?, ?, ?)",
	AddOrUpdateChatLang: "INSERT INTO chats (chat_id, chat_lang) VALUES (?, ?) ON CONFLICT DO UPDATE SET chat_lang = ? WHERE chat_id = ?",
	GetService:          "SELECT login, password FROM services WHERE service = ? and owner = ?",
	GetLang:             "SELECT chat_lang FROM chats WHERE chat_id = ?",
	DeleteService:       "DELETE FROM services WHERE service = ? and owner = ?",
}

var queriesPostgres = map[Name]Query{
	AddService:          "INSERT INTO services (service, login, password, owner) VALUES ($1, $2, $3, $4)",
	AddOrUpdateChatLang: "INSERT INTO chats (chat_id, chat_lang) VALUES ($1, $2) ON CONFLICT DO UPDATE SET chat_lang = $3 WHERE chat_id = $4",
	GetService:          "SELECT login, password FROM services WHERE service = $1 and owner = $2",
	GetLang:             "SELECT chat_lang FROM chats WHERE chat_id = $1",
	DeleteService:       "DELETE FROM services WHERE service = $1 and owner = $2",
}

// ErrNotFound occurs when query was not found.
var ErrNotFound = errors.New("the query was not found")

// ErrNilStatement occurs query statement is nil.
var ErrNilStatement = errors.New("query statement is nil")

var statements = make(map[Name]*sql.Stmt, 10)

// Prepare prepares all queries for db instance.
func Prepare(DB *sql.DB, vendor string) error {
	var queries map[Name]Query
	switch vendor {
	case "sqlite":
		queries = queriesSqlite
	case "postgres":
		queries = queriesPostgres
	}

	for n, q := range queries {
		prep, err := DB.Prepare(string(q))
		if err != nil {
			return err
		}
		statements[n] = prep
	}
	return nil
}

// GetPreparedStatement returns *sql.Stmt by name of query.
func GetPreparedStatement(name int) (*sql.Stmt, error) {
	stmt, ok := statements[Name(name)]
	if !ok {
		return nil, ErrNotFound
	}

	if stmt == nil {
		return nil, ErrNilStatement
	}

	return stmt, nil
}

// Close closes all prepared statements.
func Close() error {
	for _, stmt := range statements {
		err := stmt.Close()
		if err != nil {
			return fmt.Errorf("error closing statement: %w", err)
		}
	}

	return nil
}
