CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    avatar_url TEXT,
    
    -- Уникальные ограничения
    CONSTRAINT unique_username UNIQUE (username),
    CONSTRAINT unique_email UNIQUE (email)
);

CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    user_id bigint NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    title TEXT,
    type VARCHAR(255),
    
    -- Внешний ключ для связи с таблицей users
    CONSTRAINT fk_event_user 
        FOREIGN KEY (user_id) 
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE TABLE friends (
    user_id INT NOT NULL,
    friend_id INT NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'accepted', 'rejected')),
    request_at TIMESTAMP NOT NULL,
    confirmed_at TIMESTAMP,
    
    -- Внешний ключ для связи с таблицей users (пользователь)
    CONSTRAINT fk_friend_user
        FOREIGN KEY (user_id) 
        REFERENCES users(id)
        ON DELETE CASCADE,
    
    -- Внешний ключ для связи с таблицей users (друг)
    CONSTRAINT fk_friend_friend
        FOREIGN KEY (friend_id) 
        REFERENCES users(id)
        ON DELETE CASCADE,
    
    -- Проверяет уникальность записи user_id, friend_id (однонаправленная уникальность)
    CONSTRAINT unique_friendship UNIQUE (user_id, friend_id)
);

CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    owner_id bigint NOT NULL,
    name VARCHAR(255) NOT NULL,
    confirmed_time TIMESTAMP,
    
    -- Внешний ключ для связи с таблицей users
    CONSTRAINT fk_group_owner
        FOREIGN KEY (owner_id)
        REFERENCES users(id)
        ON DELETE RESTRICT
);

CREATE TABLE group_members (
    group_id INT NOT NULL,
    user_id bigint NOT NULL,
    confirmed_time TIMESTAMP,
    
    -- Внешний ключ для связи с таблицей groups
    CONSTRAINT fk_group_member_group
        FOREIGN KEY (group_id) 
        REFERENCES groups(id)
        ON DELETE CASCADE,
    
    -- Внешний ключ для связи с таблицей users
    CONSTRAINT fk_group_member_user
        FOREIGN KEY (user_id) 
        REFERENCES users(id)
        ON DELETE CASCADE
);