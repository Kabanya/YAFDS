-- clear_db.sql: truncates all public tables and resets sequences

-- Truncate all tables in public schema
DO $$
DECLARE
    r RECORD;
BEGIN
    FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
        EXECUTE format('TRUNCATE TABLE public.%I CASCADE;', r.tablename);
    END LOOP;
END
$$;

-- Reset all sequences in public schema to start from 1
DO $$
DECLARE
    s RECORD;
BEGIN
    FOR s IN (SELECT sequence_schema, sequence_name FROM information_schema.sequences WHERE sequence_schema = 'public') LOOP
        EXECUTE format('ALTER SEQUENCE %I.%I RESTART WITH 1;', s.sequence_schema, s.sequence_name);
    END LOOP;
END
$$;
