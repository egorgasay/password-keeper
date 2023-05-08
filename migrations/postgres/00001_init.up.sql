CREATE TABLE chats (
    chat_id INTEGER PRIMARY KEY,
    chat_lang VARCHAR(20)
);
CREATE TABLE services (
    service TEXT,
    login TEXT,
    password TEXT,
    owner INTEGER REFERENCES chats(chat_id)
);