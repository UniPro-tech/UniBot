table "members" {
  schema = schema.public
  column "user_id" {
    type    = bigint
    null    = false
    comment = "Discord User ID (Snowflake ID)"
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
    comment = "Is this user suspended"
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
    columns = [column.guild_id, column.user_id]
  }
  foreign_key "fk_guild_id" {
    columns     = [column.guild_id]
    ref_columns = [table.guilds.column.id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
  foreign_key "fk_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
  index "idx_members_is_suspended" {
    columns = [column.is_suspended]
  }
}


trigger "trg_members_updated_at" {
  on = table.members
  before {
    update = true
  }
  for = "ROW"
  execute {
    function = function.set_updated_at
  }
}