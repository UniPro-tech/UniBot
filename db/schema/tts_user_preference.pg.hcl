table "tts_user_preference" {
  schema = schema.public
  column "user_id" {
    type = bigint
    null = false
  }
  column "speaker_id" {
    type    = int
    null    = false
    default = 0
  }
  column "speed" {
    type    = int
    null    = false
    default = 100
  }
  primary_key {
    columns = [column.user_id]
  }
  foreign_key "fk_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
  check "chk_tts_user_pref_speed" {
    expr = "speed >= 0 AND speed <= 200"
  }
}