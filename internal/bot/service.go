package bot

import tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var (
	startMessageRU = `Привет!👋
Я буду хранить твои пароли, чтобы ты не запоминал каждый 🔐. 

ℹ️ Мои команды:

/set имя_сервиса логин пароль - сохранит пароль для указанного сервиса
/get имя_сервиса - покажет твой пароль для указанного сервиса
/del имена_серверов - удалит пароль для указанного сервиса

Я буду удалять наши сообщения каждые %d секунд, чтобы никто не мог узнать какие пароли ты вводил 🤫.
`
	startMessageEN = `Hi!👋 I'm a password bot 🤖. 
I'm going to save your passwords so that you don't forget them 🔐.

ℹ️ My commands: 
/set service_name login password - saves your password for the specified service
/get service_name - shows your password for the specified service
/del service_names - deletes your password for the specified service

I'll delete my messages every %d seconds, so that nobody can see what you've entered 🤫.`
)

var allMessages = map[string]messages{
	start: {
		Russian: startMessageRU,
		English: startMessageEN,
	},
	set: {
		Russian: setMessageRU,
		English: setMessageEN,
	},
	setErr: {
		Russian: setErrMessageRU,
		English: setErrMessageEN,
	},
	get: {
		Russian: getMessageRU,
		English: getMessageEN,
	},
	getErr: {
		Russian: getErrMessageRU,
		English: getErrMessageEN,
	},
	del: {
		Russian: delMessageRU,
		English: delMessageEN,
	},
	delErr: {
		Russian: delErrMessageRU,
		English: delErrMessageEN,
	},

	wrongInputErr: {
		Russian: wrongInputErrRU,
		English: wrongInputErrEN,
	},
	serviceNotFoundErr: {
		Russian: serviceNotFoundErrRU,
		English: serviceNotFoundErrEN,
	},
}

// Group of constants for bot messages
const (
	setMessageRU    = "Сохранено! ✅"
	setErrMessageRU = "Что-то пошло не так! ⛔️"
	setMessageEN    = "Saved! ✅"
	setErrMessageEN = "Something went wrong! ⛔️"

	delMessageRU    = "Удаленно! 🗑"
	delErrMessageRU = "Ошибка при удалении! ⛔️"
	delMessageEN    = "Deleted! 🗑"
	delErrMessageEN = "Error during deletion! ⛔️"

	getMessageRU    = "🔐 %s\n👤 Логин: %s\n🔑 Пароль: %s\n"
	getErrMessageRU = "Что-то пошло не так! ⚒"
	getMessageEN    = "🔐 %s\n👤 Login: %s\n🔑 Password: %s\n"
	getErrMessageEN = "Something went wrong! ⚒"

	wrongInputErrRU = "Неправильные аргументы для команды ⛔️"
	wrongInputErrEN = "Wrong input for command ⛔️"

	serviceNotFoundErrRU = "Сервис не найден ❌"
	serviceNotFoundErrEN = "Service not found ❌"
)

// Group of constants for handling messages from user.
const (
	start = "start"

	set    = "set"
	setErr = "setErr"

	get    = "get"
	getErr = "getErr"

	del    = "del"
	delErr = "delErr"

	hide = "hide"

	changeLang = "changeLang"
	change     = "change"

	wrongInputErr      = "Wrong input for command"
	serviceNotFoundErr = "Service not found"
)

const (
	hideKeyboard    = "hideKeyboard"
	setLangKeyboard = "setLangKeyboard"
	startKeyboard   = "startKeyboard"
)

// Map of  keyboard buttons.
var allKeyboards = map[string]keyboards{
	hideKeyboard: {
		Russian: tgapi.NewInlineKeyboardMarkup(
			tgapi.NewInlineKeyboardRow(
				tgapi.NewInlineKeyboardButtonData("Спрятать \U0001FAE3", hide),
			),
		),
		English: tgapi.NewInlineKeyboardMarkup(
			tgapi.NewInlineKeyboardRow(
				tgapi.NewInlineKeyboardButtonData("Hide \U0001FAE3", hide),
			),
		),
	},

	setLangKeyboard: {
		Russian: tgapi.NewInlineKeyboardMarkup(
			tgapi.NewInlineKeyboardRow(
				tgapi.NewInlineKeyboardButtonData("Русский 🇷🇺", change+"::ru"),
				tgapi.NewInlineKeyboardButtonData("English 🇺🇸", change+"::en"),
			),
		),
		English: tgapi.NewInlineKeyboardMarkup(
			tgapi.NewInlineKeyboardRow(
				tgapi.NewInlineKeyboardButtonData("Русский 🇷🇺", change+"::ru"),
				tgapi.NewInlineKeyboardButtonData("English 🇺🇸", change+"::en"),
			),
		),
	},

	startKeyboard: {
		Russian: tgapi.NewInlineKeyboardMarkup(
			tgapi.NewInlineKeyboardRow(
				tgapi.NewInlineKeyboardButtonData("Сменить язык 🌍", changeLang),
			),
		),
		English: tgapi.NewInlineKeyboardMarkup(
			tgapi.NewInlineKeyboardRow(
				tgapi.NewInlineKeyboardButtonData("Change language 🌍", changeLang),
			),
		),
	},
}
