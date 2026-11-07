table "guilds" {
  schema = schema.public
  column "id" {
    type    = bigint
    null    = false
    comment = "Discord Guild ID (Snowflake ID)"
  }
  column "is_admin" {
    type    = bool
    null    = false
    default = false
    comment = "Which is it admin guild"
  }
  column "is_suspended" {
    type    = bool
    null    = false
    default = false
    comment = "Is this guild suspended"
  }
  column "is_tts_enabled" {
    type    = bool
    null    = false
    default = false
    comment = "Is available TTS in this guild"
  }
  column "created_at" {
    type    = timestamptz
    null    = false
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    type    = timestamptz
    null    = false
    default = sql("CURRENT_TIMESTAMP")
  }
  column "deleted_at" {
    type    = timestamptz
    null    = true
    default = null
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_guilds_is_suspended" {
    columns = [column.is_suspended]
  }
}

trigger "trg_guilds_updated_at" {
  on = table.guilds
  before {
    update = true
  }
  for = "ROW"
  execute {
    function = function.set_updated_at
  }
}