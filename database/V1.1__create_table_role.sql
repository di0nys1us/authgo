CREATE TABLE "authgo"."role" (
    "id" BIGSERIAL NOT NULL,
    "name" VARCHAR(32) NOT NULL,

    PRIMARY KEY ("id"),
    UNIQUE ("name")
);
