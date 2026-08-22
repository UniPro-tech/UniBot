variable "pg_dsn" {
    type    = string
    default = getenv("PG_DSN")
}

env "local" {
  url = getenv("PG_DSN")
  dev = "docker://postgres/18/dev"
  schema {
    src = "file://schema"
  }
  migration {
    dir =  "file://migrations"
  }
}

env "prod" {
    url = var.pg_dsn

    migration {
        dir     = "file://migrations"
    }
}