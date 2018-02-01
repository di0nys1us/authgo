CREATE TABLE "authgo"."event" (
    "id" BIGSERIAL NOT NULL,
    "created_by" BIGINT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "type" VARCHAR(32) NOT NULL,
    "description" TEXT NOT NULL,

    PRIMARY KEY ("id"),
    FOREIGN KEY ("created_by") REFERENCES "authgo"."user" ("id"),
    FOREIGN KEY ("type") REFERENCES "authgo"."event_type" ("name")
);
