-- -------------------------------------------------------------
-- TablePlus 3.3.0(300)
--
-- https://tableplus.com/
--
-- Database: sound_test
-- Generation Time: 2020-04-16 12:30:17.3960 am
-- -------------------------------------------------------------


-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS albums_id_seq;

-- Table Definition
CREATE TABLE "public"."albums" (
    "id" int4 NOT NULL DEFAULT nextval('albums_id_seq'::regclass),
    "artist_id" int4,
    "name" text,
    "created" timestamptz,
    "art" text,
    "genre_id" int4,
    "year" int4,
    "duration" int4,
    "disambiguator" text,
    "starred" bool,
    "song_count" int4,
    "artist_name" text,
    "genre_name" text,
    PRIMARY KEY ("id")
);

-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS artists_id_seq;

-- Table Definition
CREATE TABLE "public"."artists" (
    "id" int4 NOT NULL DEFAULT nextval('artists_id_seq'::regclass),
    "name" text,
    "art" text,
    "starred" bool,
    "duration" int4,
    "album_count" int4,
    PRIMARY KEY ("id")
);

-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS arts_id_seq;

-- Table Definition
CREATE TABLE "public"."arts" (
    "id" int4 NOT NULL DEFAULT nextval('arts_id_seq'::regclass),
    "hash" text,
    "path" text,
    PRIMARY KEY ("id")
);

-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS genres_id_seq;

-- Table Definition
CREATE TABLE "public"."genres" (
    "id" int4 NOT NULL DEFAULT nextval('genres_id_seq'::regclass),
    "name" text,
    PRIMARY KEY ("id")
);

-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS playlist_entries_id_seq;

-- Table Definition
CREATE TABLE "public"."playlist_entries" (
    "id" int4 NOT NULL DEFAULT nextval('playlist_entries_id_seq'::regclass),
    "playlist_id" int4,
    "song_id" int4,
    "index" int4,
    PRIMARY KEY ("id")
);

-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS playlists_id_seq;

-- Table Definition
CREATE TABLE "public"."playlists" (
    "id" int4 NOT NULL DEFAULT nextval('playlists_id_seq'::regclass),
    "name" text,
    "comment" text,
    "public" bool,
    "created" timestamptz,
    "changed" timestamptz,
    "entry_count" int4,
    PRIMARY KEY ("id")
);

-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS songs_id_seq;

-- Table Definition
CREATE TABLE "public"."songs" (
    "id" int4 NOT NULL DEFAULT nextval('songs_id_seq'::regclass),
    "artist" text,
    "album_id" int4,
    "path" text,
    "title" text,
    "track" int4,
    "disc" int4,
    "genre_id" int4,
    "year" int4,
    "art" text,
    "created" timestamptz,
    "size" int8,
    "bitrate" int4,
    "duration" int4,
    "token" text,
    "provider_id" text,
    "starred" bool,
    "album_name" text,
    "album_artist_id" int4,
    "genre_name" text,
    PRIMARY KEY ("id")
);

