
-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE "public"."users_sex_type_enum" AS ENUM('Unknown', 'Male', 'Female', 'Other');

CREATE TABLE IF NOT EXISTS "users" (
    "id" bigserial NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT now(),
    "updated_at" timestamp NOT NULL DEFAULT now(),
    "deleted_at" timestamp,
    "email" character varying NOT NULL,
    "username" character varying NOT NULL,
    "password" character varying NOT NULL,
    "name" character varying NOT NULL,
    "sex_type" "public"."users_sex_type_enum" DEFAULT 'Unknown',
    "birthdate" date,

    PRIMARY KEY ("id"),
    UNIQUE ("email"),
    UNIQUE ("username")
);

-- +migrate Down
DROP TABLE IF EXISTS "users";

DROP TYPE "public"."users_sex_type_enum";
