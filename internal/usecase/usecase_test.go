package usecase

import (
	"password-keeper/internal/entity"
	"password-keeper/internal/storage"
	"reflect"
	"testing"
)

func newUseCase(t *testing.T) *UseCase {
	s, err := storage.New("test", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	uc, err := New(s, "1234567890123456")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	return uc
}

func TestUseCase_Decrypt(t *testing.T) {
	uc := newUseCase(t)

	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "too small text",
			args: args{
				text: "test1",
			},
			want: "test1",
		},
		{
			name: "ok",
			args: args{
				text: uc.Encrypt("test"),
			},
			want: "test",
		},
		{
			name: "ok #2",
			args: args{
				text: uc.Encrypt("dqwfqwedfefqfqfhkqfjqjfgqwdqwfgqwefhqvdjvqwvf"),
			},
			want: "dqwfqwedfefqfqfhkqfjqjfgqwdqwfgqwefhqvdjvqwvf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := uc.Decrypt(tt.args.text); got != tt.want {
				t.Errorf("Decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_Delete(t *testing.T) {
	uc := newUseCase(t)

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
				chatID:  190,
				service: "test",
			},
		},
		{
			name: "ok #2",
			args: args{
				chatID:  191,
				service: "teqdwqwdqdst",
			},
		},
		{
			name: "not found",
			args: args{
				chatID:  197,
				service: "teqdwqwdqdst",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				err := uc.Save(tt.args.chatID, tt.args.service, "test", "test")
				if err != nil {
					t.Errorf("Save() error = %v", err)
				}
			}
			if err := uc.Delete(tt.args.chatID, tt.args.service); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if _, err := uc.Get(tt.args.chatID, tt.args.service); err == nil {
					t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestUseCase_Encrypt(t *testing.T) {
	uc := newUseCase(t)

	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "small text",
			args: args{
				text: "test1",
			},
		},
		{
			name: "big text",
			args: args{
				text: "dqwfqwedfefqfqfhkqfjqjfgqwdqwfgqwefhqvdjvqwvf,qwld,qlfmkqmqfqjdnqf",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := uc.Encrypt(tt.args.text)

			if got := uc.Decrypt(got); got != tt.args.text {
				t.Errorf("Decrypt() = %v, want %v", got, tt.args.text)
			}
		})
	}
}

func TestUseCase_Get(t *testing.T) {
	uc := newUseCase(t)

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
				chatID:  290,
				service: "test",
			},
			want: entity.Pair{
				Login:    "test login",
				Password: "test password",
			},
		},
		{
			name: "ok #2",
			args: args{
				chatID:  291,
				service: "teqdwqwdqdst",
			},
			want: entity.Pair{
				Login:    "teqdwqwdqdst",
				Password: "XXXXXXXXXXXXXXXXXXXXX",
			},
		},
		{
			name: "not found",
			args: args{
				chatID:  291,
				service: "not found",
			},
			want:    entity.Pair{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				if err := uc.Save(tt.args.chatID, tt.args.service, tt.want.Login, tt.want.Password); err != nil {
					t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			got, err := uc.Get(tt.args.chatID, tt.args.service)
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

func TestUseCase_GetLang(t *testing.T) {
	uc := newUseCase(t)

	type args struct {
		chatID int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{
				chatID: 390,
			},
			want: "ru",
		},
		{
			name: "ok  #2",
			args: args{
				chatID: 391,
			},
			want: "en",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc.SetLang(tt.args.chatID, tt.want)
			if got := uc.GetLang(tt.args.chatID); got != tt.want {
				t.Errorf("GetLang() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_Hash(t *testing.T) {
	uc := newUseCase(t)

	type args struct {
		text string
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
				text: "test",
			},
			want:    "n4bQgYhMfWWaL+qgxVrQFaO/TxsrC4Is0V1sFbDwCgg",
			wantErr: false,
		},
		{
			name: "ok #2",
			args: args{
				text: "fnqjlnfjkqndfjkqnfjqndfqjnj",
			},
			want:    "iSRSDMOs2PAGQ5hS4CBvaU53UrRpfF8TXSOoONJpv0w",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := uc.Hash(tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Hash() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_Save(t *testing.T) {
	uc := newUseCase(t)
	type args struct {
		chatID   int64
		service  string
		login    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				chatID:   490,
				service:  "test",
				login:    "test login",
				password: "XXXXXXXXXXXXX",
			},
		},
		{
			name: "ok #2",
			args: args{
				chatID:   491,
				service:  "test2",
				login:    "test login2",
				password: "XXXXXXXXXXXXX",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := uc.Save(tt.args.chatID, tt.args.service, tt.args.login, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if p, err := uc.Get(tt.args.chatID, tt.args.service); err != nil {
					t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				} else if p.Login != tt.args.login || p.Password != tt.args.password {
					t.Errorf("Get() got = %v, want %v", p, tt.args)
				}
			}
		})
	}
}

func TestUseCase_SetLang(t *testing.T) {
	uc := newUseCase(t)
	type args struct {
		chatID int64
		lang   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ok",
			args: args{
				chatID: 590,
				lang:   "ru",
			},
		},
		{
			name: "ok #2",
			args: args{
				chatID: 590,
				lang:   "en",
			},
		},
		{
			name: "ok #3",
			args: args{
				chatID: 591,
				lang:   "ru",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc.SetLang(tt.args.chatID, tt.args.lang)
			if got := uc.GetLang(tt.args.chatID); got != tt.args.lang {
				t.Errorf("GetLang() = %v, want %v", got, tt.args.lang)
			}
		})
	}
}
