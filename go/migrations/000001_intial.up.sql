CREATE TABLE "artists" (
    "id" SERIAL,
    "name" TEXT NOT NULL,
    "starred" BOOL NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "artist_unique" UNIQUE ("name")
);

CREATE TABLE "albums" (
    "id" SERIAL,
    "artist_id" INTEGER NOT NULL REFERENCES "artists"("id"),
    "name" TEXT NOT NULL,
    "created" TIMESTAMPTZ NOT NULL,
    "disambiguator" TEXT NOT NULL,
    "starred" BOOL NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "album_unique" UNIQUE ("artist_id", "name", "disambiguator")
);

CREATE TABLE "arts" (
    "id" SERIAL,
    "hash" TEXT NOT NULL,
    "path" TEXT NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "art_hash_unique" UNIQUE ("hash"),
    CONSTRAINT "art_path_unique" UNIQUE ("path")
);

CREATE TABLE "genres" (
    "id" SERIAL,
    "name" TEXT NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "genre_unique" UNIQUE ("name")
);

CREATE TABLE "songs" (
    "id" SERIAL,
    "artist" TEXT NOT NULL,
    "album_id" INTEGER NOT NULL REFERENCES "albums"("id"),
    "path" TEXT NOT NULL,
    "title" TEXT NOT NULL,
    "track" INTEGER NOT NULL,
    "disc" INTEGER NOT NULL,
    "genre_id" INTEGER REFERENCES "genres"("id"),
    "year" INTEGER NOT NULL,
    "art" TEXT NULL REFERENCES "arts"("path"),
    "created" TIMESTAMPTZ NOT NULL,
    "size" INTEGER NOT NULL,
    "bitrate" INTEGER NOT NULL,
    "duration" INTEGER NOT NULL,
    "token" TEXT NOT NULL,
    "provider_id" TEXT NOT NULL,
    "starred" BOOL NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE "playlists" (
    "id" SERIAL,
    "name" TEXT NOT NULL,
    "comment" TEXT NOT NULL,
    "public" BOOL NOT NULL,
    "created" TIMESTAMPTZ NOT NULL,
    "changed" TIMESTAMPTZ NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE "playlist_entries" (
    "id" SERIAL,
    "playlist_id" INTEGER NOT NULL REFERENCES "playlists"("id"),
    "song_id" INTEGER NOT NULL REFERENCES "songs"("id"),
    "index" INTEGER NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "playlist_entry_unique" UNIQUE ("playlist_id", "index")
);