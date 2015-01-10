CREATE TABLE users (
   id               SERIAL                   PRIMARY KEY,
   username         TEXT           NOT NULL  UNIQUE,
   email            TEXT           NOT NULL  UNIQUE,
   password         TEXT           NOT NULL,
   join_date        TIMESTAMP      NOT NULL  DEFAULT now(),
   default_location POINT          NOT NULL
);

CREATE TABLE pages (
   id               SERIAL                   PRIMARY KEY,
   title            TEXT           NOT NULL,
   slug             TEXT           NOT NULL,
   category         INT            NOT NULL  REFERENCES categorys(id),
   description      TEXT           NOT NULL,
   user_id          INT            NOT NULL  REFERENCES users(id),
   location         POINT          NOT NULL,
   address          TEXT           NOT NULL,
   date_created     TIMESTAMP      NOT NULL  DEFAULT now()
);

CREATE TABLE categorys (
   id               SERIAL                   PRIMARY KEY,
   name             TEXT           NOT NULL
);

CREATE TABLE posts (
   id               SERIAL                   PRIMARY KEY,
   user_id          INT            NOT NULL  REFERENCES users(id),
   page_id          INT            NOT NULL  REFERENCES pages(id),
   body             TEXT           NOT NULL
);

CREATE TABLE communities (
   id               SERIAL                   PRIMARY KEY,
   name             TEXT           NOT NULL
);

CREATE TABLE community_memberships (
   user_id          INT                      PRIMARY KEY REFERENCES users(id),
   community_id     INT                      PRIMARY KEY REFERENCES communities(id),
);
