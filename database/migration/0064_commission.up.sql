-- Commission Table
CREATE TABLE commission (
                          id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                          created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                          modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                          deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
                          code varchar(30) NOT NULL,
                          purpose varchar(100) NOT NULL, -- API, OPERATOR
                          operator_id uuid DEFAULT NULL,
                          type varchar(30) NOT NULL, -- PERCENT, FIXED
                          value decimal NOT NULL,
                          description text,
                          is_active boolean NOT NULL DEFAULT true,
                          CONSTRAINT operator_id_fkey FOREIGN KEY (operator_id) REFERENCES operator
                              ON DELETE CASCADE
);

CREATE INDEX commission_code_index ON commission USING HASH (code);

