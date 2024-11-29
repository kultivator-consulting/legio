-- Package Itinerary Notes Table
CREATE TABLE package_itinerary_notes (
                                  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                                  created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                  modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                  package_itinerary_id uuid NOT NULL,
                                  note text NOT NULL,
                                  CONSTRAINT package_itinerary_id_fkey FOREIGN KEY (package_itinerary_id) REFERENCES package_itinerary
                                      ON DELETE CASCADE
);

CREATE INDEX package_itinerary_notes_package_itinerary_id_index ON package_itinerary_notes USING HASH (package_itinerary_id);

