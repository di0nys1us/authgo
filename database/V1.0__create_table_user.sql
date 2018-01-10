CREATE TABLE "authgo"."user" (
    "id" BIGSERIAL NOT NULL,
    "version" BIGINT NOT NULL DEFAULT 0,
    "first_name" VARCHAR(255) NOT NULL,
    "last_name" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "password" TEXT NOT NULL,
    "enabled" BOOLEAN NOT NULL DEFAULT TRUE,

    PRIMARY KEY ("id"),
    UNIQUE ("email")
);
