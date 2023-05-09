package postgres

import (
	"context"
	"github.com/egorgasay/dockerdb/v2"
	"log"
	"os"
	"password-keeper/internal/entity"
	prep "password-keeper/internal/storage/queries"
	"reflect"
	"testing"
)

var st *Postgres

const pathToMigrations = "file://../../../migrations/postgres"

func TestMain(m *testing.M) {
	// Write code here to run before tests
	ctx := context.TODO()
	cfg := dockerdb.CustomDB{
		DB: dockerdb.DB{
			Name:     "postgres_test_keeper",
			User:     "admin",
			Password: "admin",
		},
		Port:   "12545",
		Vendor: dockerdb.Postgres15,
	}

	err := dockerdb.Pull(ctx, dockerdb.Postgres15)
	if err != nil {
		log.Fatal(err)
		return
	}

	vdb, err := dockerdb.New(ctx, cfg)
	if err != nil {
		log.Println(err)
		os.Exit(0)
	}

	st, err = New(vdb.DB, pathToMigrations)
	if err != nil {
		log.Fatal(err)
	}

	err = prep.Prepare(st.DB.DB, "postgres")
	if err != nil {
		log.Fatal(err)
	}

	r := m.Run()

	queries := []string{
		"DROP SCHEMA public CASCADE;",
		"CREATE SCHEMA public;",
		"GRANT ALL ON SCHEMA public TO public;",
		"COMMENT ON SCHEMA public IS 'standard public schema';",
	}

	tx, err := st.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}

	for _, query := range queries {
		_, err := tx.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(r)
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.Exec(
				"INSERT INTO services (service, login, password, owner)  VALUES ($1, $2, $3, $4)",
				tt.args.service, "test", "test", tt.args.chatID,
			)
			if err != nil {
				t.Errorf("can't insert the record: %v", err)
				return
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
					"INSERT INTO services (service, login, password, owner)  VALUES ($1, $2, $3, $4)",
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
					"INSERT INTO chats (chat_id, chat_lang)  VALUES ($1, $2)",
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
