-- +migrate Up
CREATE TYPE user_source AS ENUM (
    'wechat_ios',
    'wechat_android',
    'ios',
    'android',
    'web'
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    phone_number VARCHAR(255) UNIQUE,
    avatar_url VARCHAR(255),
    source user_source,
    onboarded_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_active_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);
