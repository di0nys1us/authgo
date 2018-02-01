CREATE TABLE "authgo"."role_authority" (
    "role_id" BIGINT NOT NULL,
    "authority_id" BIGINT NOT NULL,

    PRIMARY KEY ("role_id", "authority_id"),
    FOREIGN KEY ("role_id") REFERENCES "authgo"."role" ("id"),
    FOREIGN KEY ("authority_id") REFERENCES "authgo"."authority" ("id")
);
