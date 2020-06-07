CREATE TABLE "main"."role" (
  "id"         INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
  "name"       TEXT NOT NULL UNIQUE,
  "created_at" INTEGER NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" INTEGER NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER role_updated_at BEFORE UPDATE ON "main"."role" FOR EACH ROW
BEGIN
  UPDATE "role" SET "updated_at" = CURRENT_TIMESTAMP WHERE "id" = new."id";
END;

--

CREATE TABLE "main"."door" (
  "id"         INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
  "name"       TEXT NOT NULL UNIQUE,
  "created_at" INTEGER NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" INTEGER NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER door_updated_at BEFORE UPDATE ON "main"."door" FOR EACH ROW
BEGIN
  UPDATE "door" SET "updated_at" = CURRENT_TIMESTAMP WHERE "id" = new."id";
END;

--

CREATE TABLE "main"."door_role" (
  "door_id"               INTEGER NOT NULL REFERENCES "door"(id) ON DELETE CASCADE,
  "role_id"               INTEGER NOT NULL REFERENCES "role"(id) ON DELETE CASCADE,
  "daytime_begin_seconds" INTEGER DEFAULT NULL,
  "daytime_end_seconds"   INTEGER DEFAULT NULL,
  "created_at"            INTEGER NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"            INTEGER NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE("door_id", "role_id")
);

-- CREATE TRIGGER door_role_updated_at BEFORE UPDATE ON "main"."door_role" FOR EACH ROW
-- BEGIN
--   UPDATE "door_role" SET "updated_at" = CURRENT_TIMESTAMP WHERE "id" = new."id";
-- END;

--

CREATE TABLE "main"."member" (
  "id"         INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
  "name"       TEXT NOT NULL UNIQUE,
  "created_at" INTEGER NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" INTEGER NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER member_updated_at BEFORE UPDATE ON "main"."member" FOR EACH ROW
BEGIN
  UPDATE "member" SET "updated_at" = CURRENT_TIMESTAMP WHERE "id" = new."id";
END;

--

CREATE TABLE "main"."member_role" (
  "member_id"             INTEGER NOT NULL REFERENCES "member"(id) ON DELETE CASCADE,
  "role_id"               INTEGER NOT NULL REFERENCES "role"(id) ON DELETE CASCADE,
  "expires_at"            INTEGER DEFAULT NULL,
  "created_at"            INTEGER NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"            INTEGER NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE("member_id", "role_id")
);

-- CREATE TRIGGER member_role_updated_at BEFORE UPDATE ON "main"."member_role" FOR EACH ROW
-- BEGIN
--   UPDATE "member_role" SET "updated_at" = CURRENT_TIMESTAMP WHERE "id" = new."id";
-- END;

--

CREATE TABLE "main"."key" (
  "id"         INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
  "member_id"  INTEGER REFERENCES "member"(id) ON DELETE SET NULL,
  "secret"     TEXT NOT NULL UNIQUE,
  "access_key" TEXT NOT NULL UNIQUE,
  "created_at" INTEGER NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" INTEGER NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER key_updated_at BEFORE UPDATE ON "main"."key" FOR EACH ROW
BEGIN
  UPDATE "key" SET "updated_at" = CURRENT_TIMESTAMP WHERE "id" = new."id";
END;
