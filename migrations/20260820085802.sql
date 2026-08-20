-- Modify "tts_member_preference" table
ALTER TABLE "public"."tts_member_preference" ALTER COLUMN "speaker_id" DROP NOT NULL, ALTER COLUMN "speaker_id" DROP DEFAULT, ALTER COLUMN "speed" DROP NOT NULL, ALTER COLUMN "speed" DROP DEFAULT;
