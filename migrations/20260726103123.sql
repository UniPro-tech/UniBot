-- Modify "rss_settings" table
ALTER TABLE "public"."rss_settings" ADD COLUMN "webhook_url" character varying(255) NOT NULL;
