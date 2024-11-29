-- Account Table
CREATE TABLE account (
                         id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                         created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                         modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                         deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
                         is_system boolean NOT NULL DEFAULT FALSE,
                         first_name varchar(300) DEFAULT NULL,
                         last_name varchar(300) DEFAULT NULL,
                         access_Level int DEFAULT NULL,
                         credit_balance double precision NOT NULL DEFAULT '0',
                         username varchar(400) NOT NULL,
                         user_email_address varchar(800) NOT NULL,
                         user_avatar_url varchar(4000),
                         user_zone_info varchar(60),
                         user_locale varchar(20),
                         user_password varchar(200),
                         is_locked boolean NOT NULL DEFAULT FALSE,
                         last_login timestamp with time zone DEFAULT ('-infinity'::timestamp AT TIME ZONE 'UTC'),
                         reset_password_token varchar(200) DEFAULT NULL,
                         reset_password_token_expiry timestamp with time zone DEFAULT NULL,
                         UNIQUE (username)
);

CREATE INDEX account_username_index ON account USING HASH (username);
CREATE INDEX account_user_email_address_index ON account USING HASH (user_email_address);
