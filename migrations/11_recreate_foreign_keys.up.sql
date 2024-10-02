ALTER TABLE "prices"
    DROP CONSTRAINT prices_event_id_fkey;
ALTER TABLE "prices"
    ADD FOREIGN KEY("event_id") REFERENCES "events"("id")
        ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE "events"
    DROP CONSTRAINT events_creator_id_fkey;
ALTER TABLE "events"
    ADD FOREIGN KEY("creator_id") REFERENCES "users"("id")
        ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE "tickets"
    DROP CONSTRAINT tickets_user_id_fkey;
ALTER TABLE "tickets"
    ADD FOREIGN KEY("user_id") REFERENCES "users"("id")
        ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE "tickets"
    DROP CONSTRAINT tickets_event_id_fkey;
ALTER TABLE "tickets"
    ADD FOREIGN KEY("event_id") REFERENCES "events"("id")
        ON UPDATE NO ACTION ON DELETE NO ACTION;