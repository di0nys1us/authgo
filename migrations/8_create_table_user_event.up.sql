CREATE TABLE "authgo"."user_event" (
    "user_id" BIGINT NOT NULL,
    "event_id" BIGINT NOT NULL,

    PRIMARY KEY ("event_id"),
    FOREIGN KEY ("user_id") REFERENCES "authgo"."user" ("id"),
    FOREIGN KEY ("event_id") REFERENCES "authgo"."event" ("id")
);
