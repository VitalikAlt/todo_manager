CREATE TABLE users (
    id serial UNIQUE NOT NULL,
    token varchar NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

REVOKE ALL ON TABLE users FROM PUBLIC;
ALTER TABLE public.users OWNER TO SESSION_USER;
ALTER TABLE ONLY users ADD CONSTRAINT users_pkey PRIMARY KEY (id);

CREATE INDEX user_id_token_idx ON users (id, token);