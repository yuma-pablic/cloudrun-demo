-- Add new schema named "public"
CREATE SCHEMA IF NOT EXISTS public;
-- Set comment to schema: "public"
COMMENT ON SCHEMA public IS 'standard public schema';
-- Create "users" table
CREATE TABLE public.users (
    id serial NOT NULL,
    account_name character varying(64) NOT NULL,
    passhash character varying(128) NOT NULL,
    authority boolean NOT NULL DEFAULT false,
    del_flg boolean NOT NULL DEFAULT false,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT users_account_name_key UNIQUE (account_name)
);
