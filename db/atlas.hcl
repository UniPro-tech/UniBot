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