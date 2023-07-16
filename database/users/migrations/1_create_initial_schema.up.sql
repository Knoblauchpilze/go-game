
-- Common properties of the DB
SET client_encoding = 'UTF8';

SET search_path = public;
SET default_tablespace = '';

CREATE OR REPLACE FUNCTION update_created_at() RETURNS TRIGGER AS $$
  BEGIN
    NEW.created_at = now();
    RETURN NEW;
  END;
$$ language 'plpgsql';
