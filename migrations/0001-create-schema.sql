-- Normally I'd create a database over a schema, but schemas offer DROP SCHEMA IF EXISTS
-- sql functionality which makes getting things set up for a demo easier.
-- In other words, I want to make it as painless as possible for folks to get this app up and running.
DROP SCHEMA IF EXISTS palindrome;
CREATE SCHEMA palindrome;

-- Create a db trigger to set updated at
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
