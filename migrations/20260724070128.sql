-- Create "tts_dictionary" table
CREATE TABLE "public"."tts_dictionary" (
  "user_id" bigint NOT NULL,
  "guild_id" bigint NOT NULL,
  "yomi" character varying(64) NOT NULL,
  "word" character varying(64) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("guild_id", "word")
);
-- Create index "idx_tts_dictionary_word" to table: "tts_dictionary"
CREATE INDEX "idx_tts_dictionary_word" ON "public"."tts_dictionary" ("word");
