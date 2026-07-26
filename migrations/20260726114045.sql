-- Modify "system_preference" table
ALTER TABLE "public"."system_preference" ADD COLUMN "id" uuid NOT NULL DEFAULT uuidv7();
