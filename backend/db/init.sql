CREATE TYPE user_role AS ENUM ('guest', 'user', 'admin');

CREATE TYPE user_tier AS ENUM ('guest', 'free', 'premium');

CREATE TABLE "users" (
    "id" bigserial,
    "guid" uuid NOT NULL DEFAULT gen_random_uuid(),
    "email" text NOT NULL,
    "password" text NOT NULL,
    "role" user_role NOT NULL DEFAULT 'guest',
    "tier" user_tier NOT NULL DEFAULT 'guest',
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    PRIMARY KEY ("id")
)
