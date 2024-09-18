CREATE TABLE IF NOT EXISTS "users" (
    "id" INTEGER NOT NULL UNIQUE,
    "email" VARCHAR NOT NULL,
    "password" VARCHAR NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL,
    "tg_username" VARCHAR,
    PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "events" (
    "id" INTEGER NOT NULL UNIQUE,
    "title" VARCHAR NOT NULL,
    "capacity" INTEGER NOT NULL,
    "description" TEXT,
    "beginning_time" TIMESTAMP NOT NULL,
    "end_time" TIMESTAMP NOT NULL,
    "creator_id" INTEGER NOT NULL,
    "is_public" BOOLEAN NOT NULL DEFAULT true,
    "location" VARCHAR NOT NULL,
    "is_free" BOOLEAN NOT NULL DEFAULT true,
    "preview_image" VARCHAR,
    "utc_offset" SMALLINT,
    PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "prices" (
    "id" INTEGER NOT NULL UNIQUE,
    "event_id" INTEGER NOT NULL,
    "price" INTEGER NOT NULL,
    "currency" VARCHAR NOT NULL,
    PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "tickets" (
    "id" VARCHAR NOT NULL UNIQUE,
    "user_id" INTEGER NOT NULL,
    "event_id" INTEGER NOT NULL,
    "is_used" BOOLEAN NOT NULL DEFAULT false,
    PRIMARY KEY("id")
);

ALTER TABLE "prices"
    ADD FOREIGN KEY("event_id") REFERENCES "events"("id")
        ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE "events"
    ADD FOREIGN KEY("creator_id") REFERENCES "users"("id")
        ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE "tickets"
    ADD FOREIGN KEY("user_id") REFERENCES "users"("id")
        ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE "tickets"
    ADD FOREIGN KEY("event_id") REFERENCES "events"("id")
        ON UPDATE NO ACTION ON DELETE CASCADE;