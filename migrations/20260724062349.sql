-- Modify "remind_settings" table
ALTER TABLE "public"."remind_settings" ALTER COLUMN "cron_expr" DROP NOT NULL, ALTER COLUMN "created_at" DROP NOT NULL;
