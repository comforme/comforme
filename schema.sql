CREATE TABLE users (
   id               SERIAL                   PRIMARY KEY,
   username         TEXT           NOT NULL  UNIQUE,
   email            TEXT           NOT NULL  UNIQUE,
   password         TEXT           NOT NULL,
   join_date        TIMESTAMP      NOT NULL  DEFAULT now(),
   default_location POINT,
   reset_required   BOOLEAN        NOT NULL  DEFAULT true
);

CREATE TABLE categories (
   id               SERIAL                   PRIMARY KEY,
   name             TEXT           NOT NULL
);

CREATE TABLE pages (
   id               SERIAL                   PRIMARY KEY,
   title            TEXT           NOT NULL,
   slug             TEXT           NOT NULL,
   category         INT            NOT NULL  REFERENCES categories(id),
   UNIQUE (slug, category),
   description      TEXT           NOT NULL,
   user_id          INT            NOT NULL  REFERENCES users(id),
   location         POINT          NOT NULL,
   address          TEXT           NOT NULL,
   date_created     TIMESTAMP      NOT NULL  DEFAULT now()
);
CREATE INDEX pages_title_tsvector_idx ON pages (to_tsvector('english', title));

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
   user_id          INT                      REFERENCES users(id),
   community_id     INT                      REFERENCES communities(id),
   PRIMARY KEY (user_id, community_id)
);

CREATE TABLE sessions (
   id               TEXT                     PRIMARY KEY,
   user_id           INT            NOT NULL  REFERENCES users(id),
   create_date      TIMESTAMP      NOT NULL  DEFAULT now()
);

INSERT INTO public.communities (id, name) VALUES (1, 'Lazy');
INSERT INTO public.communities (id, name) VALUES (2, 'Baboon');
INSERT INTO public.communities (id, name) VALUES (4, 'OCD');
INSERT INTO public.communities (id, name) VALUES (5, 'Stoner');
INSERT INTO public.communities (id, name) VALUES (6, 'Trans');
INSERT INTO public.communities (id, name) VALUES (7, 'BBW');
INSERT INTO public.communities (id, name) VALUES (8, 'Single');
INSERT INTO public.communities (id, name) VALUES (10, 'Boring');
INSERT INTO public.communities (id, name) VALUES (11, 'Business Owner');
INSERT INTO public.communities (id, name) VALUES (12, 'Michelle Obama');
INSERT INTO public.communities (id, name) VALUES (13, 'Genderqueer');
INSERT INTO public.communities (id, name) VALUES (14, 'Heterosexual');

INSERT INTO public.categories (id, name) VALUES (1, 'Medical');
INSERT INTO public.categories (id, name) VALUES (2, 'Food');
INSERT INTO public.categories (id, name) VALUES (3, 'Entertainment');
INSERT INTO public.categories (id, name) VALUES (4, 'Travel Services');
INSERT INTO public.categories (id, name) VALUES (5, 'Home/Garden Renovation');
