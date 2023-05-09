CREATE TABLE chats (
    chat_id INTEGER PRIMARY KEY,
    chat_lang VARCHAR(20)
);
CREATE TABLE services (
     owner INTEGER PRIMARY KEY,
     service TEXT,
     login TEXT,
     password TEXT
);