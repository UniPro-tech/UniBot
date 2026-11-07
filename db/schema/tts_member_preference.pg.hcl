table "tts_member_preference" {
  schema = schema.public
  column "user_id" {
    type = bigint
    null = false
  }
  column "guild_id" {
    type = bigint
    null = false
  }
  column "speaker_id" {
    type    = int
    null    = true
  }
  column "speed" {
    type    = int
    null    = true
  }
  primary_key {
    columns = [column.user_id, column.guild_id]
  }
  foreign_key "fk_tts_pref_member_id" {
    columns     = [column.user_id, column.guild_id]
    ref_columns = [table.members.column.user_id, table.members.column.guild_id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
  check "chk_tts_user_pref_speed" {
    expr = "speed >= 0 AND speed <= 200"
  }
}