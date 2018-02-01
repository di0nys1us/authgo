CREATE TABLE "authgo"."authority_event" (
    "authority_id" BIGINT NOT NULL,
    "event_id" BIGINT NOT NULL,

    PRIMARY KEY ("event_id"),
    FOREIGN KEY ("authority_id") REFERENCES "authority" ("id"),
    FOREIGN KEY ("event_id") REFERENCES "event" ("id")
);
