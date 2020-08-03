CREATE TABLE IF NOT EXISTS  users
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    first_name character varying(30) NOT NULL,
    last_name character varying(30) NOT NULL,
    email character varying(50) NOT NULL,
    password character varying(150) NOT NULL,
    role integer NOT NULL,
    status integer NOT NULL,
    created_at timestamp with time zone NOT NULL,
    created_by text NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    updated_by text NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX ON users (email);