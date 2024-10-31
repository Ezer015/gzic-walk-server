CREATE TABLE images(
    image_id serial PRIMARY KEY,
    image_path text,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sights(
    sight_id serial PRIMARY KEY,
    sight_name text NOT NULL,
    sight_description text NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE records(
    record_id serial PRIMARY KEY,
    image_id integer REFERENCES images(image_id) NOT NULL,
    sight_id integer REFERENCES sights(sight_id),
    sight_name text NOT NULL,
    copywriting text NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

