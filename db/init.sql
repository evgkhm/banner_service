CREATE TABLE IF NOT EXISTS banner(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    --content JSON,
    title VARCHAR(255),
    text VARCHAR(255),
    url VARCHAR(255),
    feature_id BIGSERIAL NOT NULL REFERENCES feature(id),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS tag(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    --name VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS feature(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    --name VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS banner_tag(
    banner_id BIGSERIAL NOT NULL REFERENCES banner(id),
    tag_id BIGSERIAL NOT NULL NOT NULL REFERENCES tag(id),
    PRIMARY KEY(banner_id, tag_id),
    --created_at TIMESTAMP NOT NULL DEFAULT now(),
    --updated_at TIMESTAMP NOT NULL DEFAULT now()
);