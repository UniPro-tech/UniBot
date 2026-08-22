-- Create schema "public"
CREATE SCHEMA IF NOT EXISTS "public";
-- Create enum type "discord_status_type"
CREATE TYPE "public"."discord_status_type" AS ENUM ('online', 'dnd', 'idle', 'invisible', 'offline');
-- Create enum type "discord_activity_type"
CREATE TYPE "public"."discord_activity_type" AS ENUM ('custom', 'game', 'competing', 'listening', 'watching', 'streaming');
-- Create "set_updated_at" function
CREATE FUNCTION "public"."set_updated_at" () RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
-- Create "guilds" table
CREATE TABLE "public"."guilds" (
  "id" bigint NOT NULL,
  "is_admin" boolean NOT NULL DEFAULT false,
  "is_suspended" boolean NOT NULL DEFAULT false,
  "is_tts_enabled" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_guilds_is_suspended" to table: "guilds"
CREATE INDEX "idx_guilds_is_suspended" ON "public"."guilds" ("is_suspended");
-- Set comment to column: "id" on table: "guilds"
COMMENT ON COLUMN "public"."guilds"."id" IS 'Discord Guild ID (Snowflake ID)';
-- Set comment to column: "is_admin" on table: "guilds"
COMMENT ON COLUMN "public"."guilds"."is_admin" IS 'Which is it admin guild';
-- Set comment to column: "is_suspended" on table: "guilds"
COMMENT ON COLUMN "public"."guilds"."is_suspended" IS 'Is this guild suspended';
-- Set comment to column: "is_tts_enabled" on table: "guilds"
COMMENT ON COLUMN "public"."guilds"."is_tts_enabled" IS 'Is available TTS in this guild';
-- Create "channels" table
CREATE TABLE "public"."channels" (
  "id" bigint NOT NULL,
  "guild_id" bigint NOT NULL,
  "is_suspended" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_guild_id" FOREIGN KEY ("guild_id") REFERENCES "public"."guilds" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_channels_is_suspended" to table: "channels"
CREATE INDEX "idx_channels_is_suspended" ON "public"."channels" ("is_suspended");
-- Set comment to column: "id" on table: "channels"
COMMENT ON COLUMN "public"."channels"."id" IS 'Discord Channel ID (Snowflake ID)';
-- Set comment to column: "guild_id" on table: "channels"
COMMENT ON COLUMN "public"."channels"."guild_id" IS 'Discord Guild ID (Snowflake ID)';
-- Set comment to column: "is_suspended" on table: "channels"
COMMENT ON COLUMN "public"."channels"."is_suspended" IS 'Is this channel suspended';
-- Create trigger "trg_channels_updated_at"
CREATE TRIGGER "trg_channels_updated_at" BEFORE UPDATE ON "public"."channels" FOR EACH ROW EXECUTE FUNCTION "public"."set_updated_at"();
-- Create trigger "trg_guilds_updated_at"
CREATE TRIGGER "trg_guilds_updated_at" BEFORE UPDATE ON "public"."guilds" FOR EACH ROW EXECUTE FUNCTION "public"."set_updated_at"();
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" bigint NOT NULL,
  "is_suspended" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_is_suspended" to table: "users"
CREATE INDEX "idx_users_is_suspended" ON "public"."users" ("is_suspended");
-- Set comment to column: "id" on table: "users"
COMMENT ON COLUMN "public"."users"."id" IS 'Discord User ID (Snowflake ID)';
-- Set comment to column: "is_suspended" on table: "users"
COMMENT ON COLUMN "public"."users"."is_suspended" IS 'Is this user suspended';
-- Create "members" table
CREATE TABLE "public"."members" (
  "user_id" bigint NOT NULL,
  "guild_id" bigint NOT NULL,
  "is_suspended" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("guild_id", "user_id"),
  CONSTRAINT "fk_guild_id" FOREIGN KEY ("guild_id") REFERENCES "public"."guilds" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_user_id" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_members_is_suspended" to table: "members"
CREATE INDEX "idx_members_is_suspended" ON "public"."members" ("is_suspended");
-- Set comment to column: "user_id" on table: "members"
COMMENT ON COLUMN "public"."members"."user_id" IS 'Discord User ID (Snowflake ID)';
-- Set comment to column: "guild_id" on table: "members"
COMMENT ON COLUMN "public"."members"."guild_id" IS 'Discord Guild ID (Snowflake ID)';
-- Set comment to column: "is_suspended" on table: "members"
COMMENT ON COLUMN "public"."members"."is_suspended" IS 'Is this user suspended';
-- Create trigger "trg_members_updated_at"
CREATE TRIGGER "trg_members_updated_at" BEFORE UPDATE ON "public"."members" FOR EACH ROW EXECUTE FUNCTION "public"."set_updated_at"();
-- Create "system_preference" table
CREATE TABLE "public"."system_preference" (
  "id" uuid NOT NULL DEFAULT uuidv7(),
  "status_type" "public"."discord_status_type" NOT NULL DEFAULT 'online',
  "activity_type" "public"."discord_activity_type" NULL,
  "activity_summary" character varying(100) NULL,
  PRIMARY KEY ("id")
);
-- Create trigger "trg_users_updated_at"
CREATE TRIGGER "trg_users_updated_at" BEFORE UPDATE ON "public"."users" FOR EACH ROW EXECUTE FUNCTION "public"."set_updated_at"();
-- Create "pin_settings" table
CREATE TABLE "public"."pin_settings" (
  "channel_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "guild_id" bigint NOT NULL,
  "content" text NOT NULL,
  PRIMARY KEY ("channel_id"),
  CONSTRAINT "fk_channel_id" FOREIGN KEY ("channel_id") REFERENCES "public"."channels" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_pin_member_id" FOREIGN KEY ("user_id", "guild_id") REFERENCES "public"."members" ("user_id", "guild_id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create "remind_settings" table
CREATE TABLE "public"."remind_settings" (
  "id" uuid NOT NULL DEFAULT uuidv7(),
  "user_id" bigint NOT NULL,
  "guild_id" bigint NOT NULL,
  "channel_id" bigint NOT NULL,
  "content" text NOT NULL,
  "last_run" timestamptz NOT NULL,
  "next_run" timestamptz NOT NULL,
  "cron_expr" character varying(100) NULL,
  "created_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_remind_channel_id" FOREIGN KEY ("channel_id") REFERENCES "public"."channels" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_remind_settings_member_id" FOREIGN KEY ("user_id", "guild_id") REFERENCES "public"."members" ("user_id", "guild_id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_next_run" to table: "remind_settings"
CREATE INDEX "idx_next_run" ON "public"."remind_settings" ("next_run");
-- Create "rss_settings" table
CREATE TABLE "public"."rss_settings" (
  "id" uuid NOT NULL DEFAULT uuidv7(),
  "user_id" bigint NOT NULL,
  "guild_id" bigint NOT NULL,
  "channel_id" bigint NOT NULL,
  "title" character varying(255) NULL,
  "webhook_url" character varying(255) NOT NULL,
  "url" text NOT NULL,
  "last_run" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_rss_settings_channel_id" FOREIGN KEY ("channel_id") REFERENCES "public"."channels" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_rss_settings_member_id" FOREIGN KEY ("user_id", "guild_id") REFERENCES "public"."members" ("user_id", "guild_id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_last_run" to table: "rss_settings"
CREATE INDEX "idx_last_run" ON "public"."rss_settings" ("last_run");
-- Create "tts_connections" table
CREATE TABLE "public"."tts_connections" (
  "guild_id" bigint NOT NULL,
  "channel_id" bigint NOT NULL,
  PRIMARY KEY ("guild_id"),
  CONSTRAINT "fk_channel_id" FOREIGN KEY ("channel_id") REFERENCES "public"."channels" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_guild_id" FOREIGN KEY ("guild_id") REFERENCES "public"."guilds" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create "tts_dictionary" table
CREATE TABLE "public"."tts_dictionary" (
  "id" uuid NOT NULL DEFAULT uuidv7(),
  "user_id" bigint NOT NULL,
  "guild_id" bigint NOT NULL,
  "yomi" character varying(64) NOT NULL,
  "word" character varying(64) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "uq_word_guild" UNIQUE ("guild_id", "word"),
  CONSTRAINT "fk_tts_dictionary_member_id" FOREIGN KEY ("user_id", "guild_id") REFERENCES "public"."members" ("user_id", "guild_id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_tts_dictionary_word" to table: "tts_dictionary"
CREATE INDEX "idx_tts_dictionary_word" ON "public"."tts_dictionary" ("word");
-- Create "tts_member_preference" table
CREATE TABLE "public"."tts_member_preference" (
  "user_id" bigint NOT NULL,
  "guild_id" bigint NOT NULL,
  "speaker_id" integer NOT NULL DEFAULT 0,
  "speed" integer NOT NULL DEFAULT 100,
  PRIMARY KEY ("user_id", "guild_id"),
  CONSTRAINT "fk_tts_pref_member_id" FOREIGN KEY ("user_id", "guild_id") REFERENCES "public"."members" ("user_id", "guild_id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "chk_tts_user_pref_speed" CHECK ((speed >= 0) AND (speed <= 200))
);
-- Create "tts_user_preference" table
CREATE TABLE "public"."tts_user_preference" (
  "user_id" bigint NOT NULL,
  "speaker_id" integer NOT NULL DEFAULT 0,
  "speed" integer NOT NULL DEFAULT 100,
  PRIMARY KEY ("user_id"),
  CONSTRAINT "fk_user_id" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "chk_tts_user_pref_speed" CHECK ((speed >= 0) AND (speed <= 200))
);
