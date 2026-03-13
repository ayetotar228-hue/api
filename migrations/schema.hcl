
table "users" {
  schema = schema.public
  column "id" {
    null = false
    type = serial
  }
  column "email" {
    null = false
    type = varchar(255)
  }
  column "created_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_email" {
    unique = true
    columns = [column.email]
  }
}

schema "public" {
}