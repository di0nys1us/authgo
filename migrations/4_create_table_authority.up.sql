CREATE TABLE "authgo"."authority" (
    "id" UUID NOT NULL DEFAULT uuid_generate_v1mc(),
    "version" BIGINT NOT NULL DEFAULT 0,
    "name" VARCHAR(32) NOT NULL,
    "events" JSONB NOT NULL DEFAULT '[]',

    PRIMARY KEY ("id"),
    UNIQUE ("name")
);
