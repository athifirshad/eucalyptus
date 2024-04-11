-- Drop tables in reverse order of creation
DROP TABLE IF EXISTS medication CASCADE;
DROP TABLE IF EXISTS prescription CASCADE;
DROP TABLE IF EXISTS appointment CASCADE;
DROP TABLE IF EXISTS doctor CASCADE;
DROP TABLE IF EXISTS hospital CASCADE;
DROP TABLE IF EXISTS health_record CASCADE;
DROP TABLE IF EXISTS patient CASCADE;
DROP TABLE IF EXISTS profile CASCADE;
DROP TABLE IF EXISTS public.tokens CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Drop types
DROP TYPE IF EXISTS appointment_status;
DROP TYPE IF EXISTS user_type_enum;