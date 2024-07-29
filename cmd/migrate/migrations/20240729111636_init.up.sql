CREATE TABLE IF NOT EXISTS users (
    `id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(30) NOT NULL,
    `password` VARCHAR(30) NOT NULL,
    `createdAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updatedAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (id),
    UNIQUE KEY (name)
);

INSERT INTO users (id, name, password)
VALUES 
('123e4567-e89b-12d3-a456-426614174000', 'anon', '1234'),
('123e4567-e89b-12d3-a456-426614174001', 'nona', '1234');

CREATE TABLE IF NOT EXISTS servers (
    `id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(30) NOT NULL,
    `createdAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (id),
    UNIQUE KEY (name)
);

CREATE TABLE IF NOT EXISTS server_roles (
    `id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(30) NOT NULL,

    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS server_chat_categories (
    `id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(30) NOT NULL,
    `server_id` VARCHAR(36),

    PRIMARY KEY (id),
    FOREIGN KEY (server_id) REFERENCES servers(id)
);

CREATE TABLE IF NOT EXISTS chats (
    `id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(30),
    `server_id` VARCHAR(36),
    `server_role_id` VARCHAR(36),
    `server_chat_category_id` VARCHAR(36),

    PRIMARY KEY (id),
    FOREIGN KEY (server_id) REFERENCES servers(id),
    FOREIGN KEY (server_role_id) REFERENCES server_roles(id),
    FOREIGN KEY (server_chat_category_id) REFERENCES server_chat_categories(id)
);

CREATE TABLE IF NOT EXISTS chat_or_server_users (
    `user_id` VARCHAR(36) NOT NULL,
    `chat_id` VARCHAR(36),
    `server_id` VARCHAR(36),
    `server_role_id` VARCHAR(36),

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (chat_id) REFERENCES chats(id),
    FOREIGN KEY (server_id) REFERENCES servers(id),
    FOREIGN KEY (server_role_id) REFERENCES server_roles(id)
);

CREATE TABLE IF NOT EXISTS chat_messages (
    `user_id` VARCHAR(36) NOT NULL,
    `chat_id` VARCHAR(36) NOT NULL,
    `message`  TEXT NOT NULL,
    `createdAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (chat_id) REFERENCES chats(id)
);

