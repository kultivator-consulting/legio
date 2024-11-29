-- Booking Ref Table
CREATE TABLE booking_ref (
                          id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                          created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                          modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                          deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
                          prefix varchar(10) NOT NULL,
                          initial bigint NOT NULL DEFAULT 0,
                          length integer NOT NULL DEFAULT 2,
                          sequence bigint NOT NULL DEFAULT 0
--                           agent_id uuid DEFAULT NULL
);

CREATE INDEX booking_ref_prefix_index ON booking_ref USING HASH (prefix);

