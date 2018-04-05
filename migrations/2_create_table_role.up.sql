CREATE TABLE "authgo"."role" (
    "id" UUID NOT NULL,
    "version" BIGINT NOT NULL DEFAULT 0,
    "name" VARCHAR(32) NOT NULL,
    "events" JSONB NOT NULL DEFAULT '[]',

    PRIMARY KEY ("id"),
    UNIQUE ("name")
);
