CREATE TABLE IF NOT EXISTS banner(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    content JSON,
    feature_id BIGSERIAL NOT NULL REFERENCES feature(id)
);

CREATE TABLE IF NOT EXISTS tag(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS feature(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS banner_tag(
    banner_id BIGSERIAL NOT NULL REFERENCES banner(id),
    tag_id BIGSERIAL NOT NULL NOT NULL REFERENCES tag(id)
);