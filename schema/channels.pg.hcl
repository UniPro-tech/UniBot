table "channels" {
  schema = schema.public
  column "id" {
    type    = bigint
    null    = false
    comment = "Discord Channel ID (Snowflake ID)"
  }
  column "guild_id" {
    type    = bigint
    null    = false
    comment = "Discord Guild ID (Snowflake ID)"
  }
  column "is_suspended" {
    type    = bool
    null    = false
    default = false
    comment = "Is this channel suspended"
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
  foreign_key "fk_guild_id" {
    columns     = [column.guild_id]
    ref_columns = [table.guilds.column.id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
  index "idx_channels_is_suspended" {
    columns = [column.is_suspended]
  }
}

trigger "trg_channels_updated_at" {
  on = table.channels
  before {
    update = true
  }
  for = "ROW"
  execute {
    function = function.set_updated_at
  }
}