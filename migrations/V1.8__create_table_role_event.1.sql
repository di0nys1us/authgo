CREATE TABLE "authgo"."role_event" (
    "role_id" BIGINT NOT NULL,
    "event_id" BIGINT NOT NULL,

    PRIMARY KEY ("event_id"),
    FOREIGN KEY ("role_id") REFERENCES "role" ("id"),
    FOREIGN KEY ("event_id") REFERENCES "event" ("id")
);
