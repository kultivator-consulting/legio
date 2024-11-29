-- Pass Itinerary Table
CREATE TABLE pass_itinerary (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    pass_id uuid NOT NULL,
    title varchar(120) NOT NULL,
    description text NOT NULL,
    terms text NOT NULL,
    CONSTRAINT pass_id_fkey FOREIGN KEY (pass_id) REFERENCES pass
        ON DELETE CASCADE
);

CREATE INDEX pass_itinerary_pass_id_index ON pass_itinerary USING HASH (pass_id);
