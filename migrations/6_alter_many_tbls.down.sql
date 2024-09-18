ALTER TABLE "events"
    DROP COLUMN created_at,
    DROP COLUMN updated_at; 

ALTER TABLE "prices"
    DROP COLUMN created_at,
    DROP COLUMN updated_at;

ALTER TABLE "tickets"
    DROP COLUMN created_at,
    DROP COLUMN updated_at,
    DROP COLUMN first_name,
    DROP COLUMN last_name;

ALTER TABLE "users_yookassa_settings"
    DROP COLUMN created_at,
    DROP COLUMN updated_at;

ALTER TABLE "users"
    ALTER COLUMN created_at DROP DEFAULT,
    ALTER COLUMN updated_at DROP DEFAULT;

DROP INDEX IF EXISTS idx_user_email;