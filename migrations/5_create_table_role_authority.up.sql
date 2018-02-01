CREATE TABLE "authgo"."role_authority" (
    "role_id" BIGINT NOT NULL,
    "authority_id" BIGINT NOT NULL,

    PRIMARY KEY ("role_id", "authority_id"),
    FOREIGN KEY ("role_id") REFERENCES "role" ("id"),
    FOREIGN KEY ("authority_id") REFERENCES "authority" ("id")
);
