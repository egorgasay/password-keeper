package sqlite

import (
	"database/sql"
	"log"
	"os"
	"password-keeper/internal/entity"
	"password-keeper/internal/storage/queries"
	"reflect"
	"testing"
)

var st *Sqlite3
var dbName = "test.db"

func TestMain(m *testing.M) {
	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		log.Fatalf("can't opening the db: %v", err)
	}
	defer cleanup(dbName)
	defer db.Close()

	st, err = New(db, "file://..//..//..//migrations/sqlite")
	if err != nil {
		log.Fatalf("can't creating the storage: %v", err)
	}

	err = queries.Prepare(db, "sqlite")
	if err != nil {
		log.Fatalf("error preparing db: %v", err)
	}

	m.Run()

}

func cleanup(filename string) {
	if err := os.Remove(filename); err != nil {
		log.Fatalf("can't remove the db: %v", err)
	}
}

func TestDB_Delete(t *testing.T) {
	type args struct {
		chatID  int64
		service string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				chatID:  1,
				service: "vk",
			},
		},
		{
			name: "ok 2",
			args: args{
				chatID:  2,
				service: "yandex",
			},
		},
		{
			name: "not found",
			args: args{
				chatID:  3,
				service: "test",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				_, err := st.Exec(
					"INSERT INTO services (service, login, password, owner)  VALUES (?, ?, ?, ?)",
					tt.args.service, "test", "test", tt.args.chatID,
				)
				if err != nil {
					t.Errorf("can't insert the record: %v", err)
					return
				}
			}

			if err := st.Delete(tt.args.chatID, tt.args.service); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_Get(t *testing.T) {
	type args struct {
		chatID  int64
		service string
	}
	tests := []struct {
		name    string
		args    args
		want    entity.Pair
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				chatID:  1,
				service: "vk.com",
			},
			want: entity.Pair{
				Login:    "test",
				Password: "test",
			},
		},
		{
			name: "ok 2",
			args: args{
				chatID:  2,
				service: "yandex.ru",
			},
			want: entity.Pair{
				Login:    "test",
				Password: "test",
			},
		},
		{
			name: "not found",
			args: args{
				chatID:  3,
				service: "test.com",
			},
			wantErr: true,
			want:    entity.Pair{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				_, err := st.Exec(
					"INSERT INTO services (service, login, password, owner)  VALUES (?, ?, ?, ?)",
					tt.args.service, "test", "test", tt.args.chatID,
				)
				if err != nil {
					t.Errorf("can't insert the record: %v", err)
					return
				}
			}
			got, err := st.Get(tt.args.chatID, tt.args.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_GetLang(t *testing.T) {
	type args struct {
		chatID int64
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				chatID: 11,
			},
			want: "ru",
		},
		{
			name: "ok 2",
			args: args{
				chatID: 22,
			},
			want: "en",
		},
		{
			name: "not found",
			args: args{
				chatID: 33,
			},
			wantErr: true,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				_, err := st.Exec(
					"INSERT INTO chats (chat_id, chat_lang)  VALUES (?, ?)",
					tt.args.chatID, tt.want)
				if err != nil {
					t.Errorf("can't insert the record: %v", err)
				}
			}
			got, err := st.GetLang(tt.args.chatID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLang() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetLang() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_Save(t *testing.T) {
	type args struct {
		chatID  int64
		service string
		pair    entity.Pair
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				chatID:  111,
				service: "vk.ru.com",
				pair: entity.Pair{
					Login:    "test",
					Password: "XXXX",
				},
			},
		},
		{
			name: "ok 2",
			args: args{
				chatID:  222,
				service: "yan2x.ru",
				pair: entity.Pair{
					Login:    "test",
					Password: "XXXX",
				},
			},
		},
		{
			name: "duplicate",
			args: args{
				chatID:  222,
				service: "yan2x.ru",
				pair: entity.Pair{
					Login:    "test",
					Password: "XXXX",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := st.Save(tt.args.chatID, tt.args.service, tt.args.pair); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				got, err := st.Get(tt.args.chatID, tt.args.service)
				if err != nil {
					t.Errorf("can't get the record: %v", err)
				}
				if !reflect.DeepEqual(got, tt.args.pair) {
					t.Errorf("Save() got = %v, want %v", got, tt.args.pair)
				}
			}
		})
	}
}

func TestDB_SetLang(t *testing.T) {
	type args struct {
		chatID int64
		lang   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				chatID: 111,
				lang:   "ru",
			},
		},
		{
			name: "ok 2",
			args: args{
				chatID: 222,
				lang:   "en",
			},
		},
		{
			name: "duplicate",
			args: args{
				chatID: 222,
				lang:   "en",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := st.SetLang(tt.args.chatID, tt.args.lang); (err != nil) != tt.wantErr {
				t.Errorf("SetLang() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				got, err := st.GetLang(tt.args.chatID)
				if err != nil {
					t.Errorf("can't get the record: %v", err)
				}
				if got != tt.args.lang {
					t.Errorf("Save() got = %v, want %v", got, tt.args.lang)
				}
			}
		})
	}
}
