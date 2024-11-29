-- Operator Table
CREATE TABLE operator (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
    operator_name varchar(300) NOT NULL,
    location_id uuid DEFAULT NULL,
    operator_code varchar(200) DEFAULT NULL,
    instruction_cust text DEFAULT NULL,
    email_address varchar(200) DEFAULT NULL,
    operator_description text,
    operator_bio text,
    video_url varchar(4000) DEFAULT NULL,
    star_rating int NOT NULL,
    operator_image varchar(4000) DEFAULT NULL,
    website_url varchar(4000) DEFAULT NULL,
    zone_info varchar(200) NOT NULL,
    locale varchar(100) NOT NULL,
    fare_url varchar(4000) NOT NULL,
    phone varchar(60) NOT NULL,
    is_pass_code boolean NOT NULL DEFAULT FALSE,
    is_active boolean NOT NULL DEFAULT TRUE,
    CONSTRAINT location_id_fkey FOREIGN KEY (location_id) REFERENCES location
        ON DELETE NO ACTION
);

CREATE INDEX operator_name_index ON operator USING HASH (operator_name);
CREATE INDEX operator_location_id_index ON operator USING HASH (location_id);
CREATE INDEX operator_email_address_index ON operator USING HASH (email_address);
CREATE INDEX operator_operator_code_index ON operator USING HASH (operator_code);
CREATE INDEX operator_is_active_index ON operator USING HASH (is_active);
