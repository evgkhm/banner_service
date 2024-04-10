CREATE TABLE IF NOT EXISTS feature(
  id BIGSERIAL NOT NULL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS tag(
  id BIGSERIAL NOT NULL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS banner(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    title VARCHAR(255),
    text VARCHAR(255),
    url VARCHAR(255),
    feature_id BIGSERIAL NOT NULL REFERENCES feature(id) ON DELETE CASCADE ON UPDATE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS banner_tag(
    banner_id BIGSERIAL NOT NULL REFERENCES banner(id) ON DELETE CASCADE ON UPDATE CASCADE,
    tag_id BIGSERIAL NOT NULL NOT NULL REFERENCES tag(id) ON DELETE CASCADE ON UPDATE CASCADE,
    PRIMARY KEY(banner_id, tag_id)
);

--insert into tables values
INSERT INTO feature (created_at, updated_at) VALUES (now(), now());
INSERT INTO tag (created_at, updated_at) VALUES (now(), now());
INSERT INTO tag (created_at, updated_at) VALUES (now(), now());
INSERT INTO banner (title, text, url, feature_id, is_active, created_at, updated_at) VALUES ('some_title', 'some_text', 'some_url', 1, true, now(), now());
INSERT INTO banner_tag (banner_id, tag_id) VALUES (1, 1);
INSERT INTO banner_tag (banner_id, tag_id) VALUES (1, 2);

INSERT INTO feature (created_at, updated_at) VALUES (now(), now());
INSERT INTO tag (created_at, updated_at) VALUES (now(), now());
INSERT INTO tag (created_at, updated_at) VALUES (now(), now());
INSERT INTO banner (title, text, url, feature_id, is_active, created_at, updated_at) VALUES ('some_title', 'some_text', 'some_url', 2, false, now(), now());
INSERT INTO banner_tag (banner_id, tag_id) VALUES (2, 3);
INSERT INTO banner_tag (banner_id, tag_id) VALUES (2, 4);