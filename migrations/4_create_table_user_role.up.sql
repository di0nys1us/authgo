CREATE TABLE "authgo"."user_role" (
    "user_id" BIGINT NOT NULL,
    "role_id" BIGINT NOT NULL,

    PRIMARY KEY ("user_id", "role_id"),
    FOREIGN KEY ("user_id") REFERENCES "authgo"."user" ("id"),
    FOREIGN KEY ("role_id") REFERENCES "authgo"."role" ("id")
);