CREATE TABLE "authgo"."user" (
    "id" UUID NOT NULL DEFAULT uuid_generate_v1mc(),
    "version" BIGINT NOT NULL DEFAULT 0,
    "events" JSONB NOT NULL DEFAULT '[]',
    "deleted" BOOLEAN NOT NULL DEFAULT FALSE,
    "first_name" VARCHAR(255) NOT NULL,
    "last_name" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "password" TEXT NOT NULL,
    "enabled" BOOLEAN NOT NULL DEFAULT TRUE,

    PRIMARY KEY ("id"),
    UNIQUE ("email")
);
