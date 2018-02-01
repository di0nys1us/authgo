CREATE TABLE "authgo"."authority" (
    "id" BIGSERIAL NOT NULL,
    "version" BIGINT NOT NULL DEFAULT 0,
    "name" VARCHAR(32) NOT NULL,

    PRIMARY KEY ("id"),
    UNIQUE ("name")
);
