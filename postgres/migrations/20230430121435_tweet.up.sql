CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS tweets(
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v1(),
    body VARCHAR(255)  NOT NULL,
    user_id UUID NOT NULL references users (id) on delete cascade ,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
    );