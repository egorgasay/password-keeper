# goKeeper - telegram bot for storing passwords [![CI](https://github.com/egorgasay/password-keeper/actions/workflows/go.yml/badge.svg)](https://github.com/egorgasay/password-keeper/actions/workflows/go.yml)

### 🔍️ Purpose
The bot is designed so that you don't have to remember every password.

### ✨ Features
- 🗑 Deleting all messages in the interval specified in the config,
- 🤐 Hiding messages from the chat by clicking on the interactive button,
- 🌎 Each user has the opportunity to choose a language to communicate with the bot (Russian or English),
- ℹ️ The ability to choose between two databases: Postgresql and Sqlite,
- 👤 Each user has their own space, so one user will not be able to access the passwords of another.
- ⚡️ All passwords are stored in RAM for the fastest response and in a database to ensure durability.

### ⚙️ Configuration
```python
./keeper

-token=YOUR_BOT_KEY 
example: -token=fQDwfQDqwDQfqqDqdfXDQDqdq3q13e1h1

-key=YOUR_KEY_FOR_ENCRYPTION 
example: -key=e2qed678901qwd56 (16, 24 or 32 symbols)

-storage=sqlite_or_postgres 
example: -storage=postgres

-dsn=CONNECTION_STRING 
example: -dsn=my.db

-interval=DELETION_INTERVAL 
example: -interval=1s
```

### ⏬ Installation

```bash
git clone https://github.com/egorgasay/password-keeper
cd password-keeper
export TELEGRAM_BOT_KEY=YOUR_BOT_KEY
export ENCRYPTION_KEY=YOUR_KEY_FOR_ENCRIPTION
make run
```

### 🐋 Docker
```bash
docker-compose up
```

### ✅ Run tests

```bash
make test
```
