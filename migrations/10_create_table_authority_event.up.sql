CREATE TABLE "authgo"."authority_event" (
    "authority_id" BIGINT NOT NULL,
    "event_id" BIGINT NOT NULL,

    PRIMARY KEY ("event_id"),
    FOREIGN KEY ("authority_id") REFERENCES "authgo"."authority" ("id"),
    FOREIGN KEY ("event_id") REFERENCES "authgo"."event" ("id")
);
