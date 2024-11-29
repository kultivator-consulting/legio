CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP ROLE IF EXISTS cortexdbuser;

CREATE USER cortexdbuser WITH PASSWORD 'TkSvWM44uUBLW9';

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO cortexdbuser;

GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO cortexdbuser;

ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO cortexdbuser;

ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON SEQUENCES TO cortexdbuser;
