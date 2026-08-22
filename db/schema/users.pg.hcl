table "users" {
  schema = schema.public
  column "id" {
    type    = bigint
    null    = false
    comment = "Discord User ID (Snowflake ID)"
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
    columns = [column.id]
  }
  index "idx_users_is_suspended" {
    columns = [column.is_suspended]
  }
}

trigger "trg_users_updated_at" {
  on = table.users
  before {
    update = true
  }
  for = "ROW"
  execute {
    function = function.set_updated_at
  }
}