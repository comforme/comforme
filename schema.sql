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
   name             TEXT           NOT NULL,
   slug             TEXT           NOT NULL
);

CREATE TABLE pages (
   id               SERIAL                   PRIMARY KEY,
   title            TEXT           NOT NULL,
   slug             TEXT           NOT NULL,
   category         INT            NOT NULL  REFERENCES categories(id) ON DELETE CASCADE,
   UNIQUE (slug, category),
   description      TEXT           NOT NULL,
   user_id          INT            NOT NULL  REFERENCES users(id) ON DELETE CASCADE,
   location         POINT,
   address          TEXT           NOT NULL,
   website          TEXT           NOT NULL,
   date_created     TIMESTAMP      NOT NULL  DEFAULT now()
);
CREATE INDEX pages_title_tsvector_idx ON pages (to_tsvector('english', title));

CREATE TABLE posts (
   id               SERIAL                   PRIMARY KEY,
   user_id          INT            NOT NULL  REFERENCES users(id) ON DELETE CASCADE,
   page_id          INT            NOT NULL  REFERENCES pages(id) ON DELETE CASCADE,
   body             TEXT           NOT NULL,
   date_created     TIMESTAMP      NOT NULL  DEFAULT now()
);

CREATE TABLE communities (
   id               SERIAL                   PRIMARY KEY,
   name             TEXT           NOT NULL
);

CREATE TABLE community_memberships (
   user_id          INT                      REFERENCES users(id) ON DELETE CASCADE,
   community_id     INT                      REFERENCES communities(id) ON DELETE CASCADE,
   PRIMARY KEY (user_id, community_id)
);

CREATE TABLE sessions (
   id               TEXT                     PRIMARY KEY,
   user_id          INT            NOT NULL  REFERENCES users(id) ON DELETE CASCADE,
   create_date      TIMESTAMP      NOT NULL  DEFAULT now()
);

-- Users are automatically enrolled in community #1 at signup, but are given the option to opt-out.
INSERT INTO public.communities (name, id) VALUES ('Lazy', 1);

INSERT INTO public.communities (name) VALUES ('Baboon');
INSERT INTO public.communities (name) VALUES ('Detail Oriented');
INSERT INTO public.communities (name) VALUES ('OCD');
INSERT INTO public.communities (name) VALUES ('Crazy Cat Person');
INSERT INTO public.communities (name) VALUES ('Python 2 Holdout');
INSERT INTO public.communities (name) VALUES ('Stoner');
INSERT INTO public.communities (name) VALUES ('BBW');
INSERT INTO public.communities (name) VALUES ('Boring');
INSERT INTO public.communities (name) VALUES ('Single');
INSERT INTO public.communities (name) VALUES ('Business Owner');
INSERT INTO public.communities (name) VALUES ('Michelle Obama');
INSERT INTO public.communities (name) VALUES ('Chia Pet Enthusiast');
INSERT INTO public.communities (name) VALUES ('Gay');
INSERT INTO public.communities (name) VALUES ('Lesbian');
INSERT INTO public.communities (name) VALUES ('Transexual/Transgender');
INSERT INTO public.communities (name) VALUES ('Genderqueer');
INSERT INTO public.communities (name) VALUES ('Heterosexual');
INSERT INTO public.communities (name) VALUES ('Homosexual');
INSERT INTO public.communities (name) VALUES ('Asexual');
INSERT INTO public.communities (name) VALUES ('Bisexual');
INSERT INTO public.communities (name) VALUES ('Transexual/Transgender');

INSERT INTO public.categories (name, slug) VALUES ('Medical', 'medical');
INSERT INTO public.categories (name, slug) VALUES ('Food', 'food');
INSERT INTO public.categories (name, slug) VALUES ('Entertainment', 'entertainment');
INSERT INTO public.categories (name, slug) VALUES ('Travel Services', 'travel');
INSERT INTO public.categories (name, slug) VALUES ('Home/Garden', 'home-garden');
