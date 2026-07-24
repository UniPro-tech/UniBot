table "tts_dictionary" {
  schema = schema.public
  column "user_id" {
    type = bigint
    null = false
  }
  column "guild_id" {
    type = bigint
    null = false
  }
  column "yomi" {
    type = varchar(64)
    null = false
  }
  column "word" {
    type = varchar(64)
    null = false
  }
  column "created_at" {
    type    = timestamptz
    null    = false
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.guild_id, column.word]
  }
  index "idx_tts_dictionary_word" {
    columns = [column.word]
  }
  foreign_key "fk_tts_dictionary_member_id" {
    columns     = [column.user_id, column.guild_id]
    ref_columns = [table.members.column.user_id, table.members.column.guild_id]
  }
}