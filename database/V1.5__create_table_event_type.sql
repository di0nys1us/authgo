CREATE TABLE "authgo"."event_type" (
    "id" BIGSERIAL NOT NULL,
    "name" VARCHAR(32) NOT NULL,

    PRIMARY KEY ("id"),
    UNIQUE ("name")
);
