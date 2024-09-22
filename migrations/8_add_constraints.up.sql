ALTER TABLE prices ADD CONSTRAINT unique_event_currency UNIQUE (event_id, currency);
ALTER TABLE users  ADD CONSTRAINT unique_user_email UNIQUE (email);
