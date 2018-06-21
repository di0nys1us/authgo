CREATE TABLE "authgo"."role" (
    "id" UUID NOT NULL DEFAULT uuid_generate_v1mc(),
    "version" BIGINT NOT NULL DEFAULT 0,
    "events" JSONB NOT NULL DEFAULT '[]',
    "name" VARCHAR(32) NOT NULL,

    PRIMARY KEY ("id"),
    UNIQUE ("name")
);
