CREATE TABLE IF NOT EXISTS "users_yookassa_settings" (
    "id" SERIAL NOT NULL UNIQUE,
    "user_id" INTEGER NOT NULL,
    "shop_id" VARCHAR NOT NULL,
    "shop_key" VARCHAR NOT NULL,
    PRIMARY KEY("id")
);

ALTER TABLE "users_yookassa_settings"
    ADD FOREIGN KEY ("user_id") REFERENCES "users"("id")
        ON UPDATE NO ACTION ON DELETE CASCADE;