-- Modify "pin_settings" table
ALTER TABLE "public"."pin_settings" DROP CONSTRAINT "fk_user_id", ADD COLUMN "guild_id" bigint NOT NULL, ADD CONSTRAINT "fk_pin_member_id" FOREIGN KEY ("user_id", "guild_id") REFERENCES "public"."members" ("user_id", "guild_id") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "remind_settings" table
ALTER TABLE "public"."remind_settings" DROP CONSTRAINT "fk_channel_id", DROP CONSTRAINT "fk_user_id", ADD COLUMN "guild_id" bigint NOT NULL, ADD CONSTRAINT "fk_remind_channel_id" FOREIGN KEY ("channel_id") REFERENCES "public"."channels" ("id") ON UPDATE CASCADE ON DELETE CASCADE, ADD CONSTRAINT "fk_remind_settings_member_id" FOREIGN KEY ("user_id", "guild_id") REFERENCES "public"."members" ("user_id", "guild_id") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "rss_settings" table
ALTER TABLE "public"."rss_settings" DROP CONSTRAINT "fk_channel_id", DROP CONSTRAINT "fk_user_id", ADD COLUMN "guild_id" bigint NOT NULL, ADD CONSTRAINT "fk_rss_settings_channel_id" FOREIGN KEY ("channel_id") REFERENCES "public"."channels" ("id") ON UPDATE CASCADE ON DELETE CASCADE, ADD CONSTRAINT "fk_rss_settings_member_id" FOREIGN KEY ("user_id", "guild_id") REFERENCES "public"."members" ("user_id", "guild_id") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "tts_dictionary" table
ALTER TABLE "public"."tts_dictionary" ADD CONSTRAINT "fk_tts_dictionary_member_id" FOREIGN KEY ("user_id", "guild_id") REFERENCES "public"."members" ("user_id", "guild_id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "tts_member_preference" table
ALTER TABLE "public"."tts_member_preference" DROP CONSTRAINT "fk_member_id", ADD CONSTRAINT "fk_tts_pref_member_id" FOREIGN KEY ("user_id", "guild_id") REFERENCES "public"."members" ("user_id", "guild_id") ON UPDATE CASCADE ON DELETE CASCADE;
