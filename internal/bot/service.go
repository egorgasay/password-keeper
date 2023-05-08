package bot

// Group of constants for bot messages
const (
	startMessage = `Привет! Я буду хранить твои пароли, чтобы ты не запоминал каждый)
Мои команды:
/set имя_сервиса логин пароль - сохранит пароль для указанного сервиса
/get имя_сервиса - покажет твой пароль для указанного сервиса
/del имена_серверов - удалит пароль для указанного сервиса
`
	setMessage = "Сохранено!"
)

// Group of constants for handling messages from user.
const (
	start = "start"
	set   = "set"
	get   = "get"
	del   = "del"
)
