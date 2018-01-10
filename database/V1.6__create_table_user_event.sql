CREATE TABLE "authgo"."user_event" (
    "id" BIGSERIAL NOT NULL,
    "user_id" BIGINT NOT NULL,
    "created_by" BIGINT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "type" BIGINT NOT NULL,
    "description" TEXT NOT NULL,

    PRIMARY KEY ("id"),
    FOREIGN KEY ("user_id") REFERENCES "user" ("id"),
    FOREIGN KEY ("created_by") REFERENCES "user" ("id"),
    FOREIGN KEY ("type") REFERENCES "event_type" ("id")
);
