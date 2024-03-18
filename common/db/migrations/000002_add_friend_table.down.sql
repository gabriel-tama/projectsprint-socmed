ALTER TABLE "friends" DROP CONSTRAINT "user_id_fk";
ALTER TABLE "friends" DROP CONSTRAINT "friend_id_fk";

DROP TABLE IF EXISTS "friends";

