variable "db_url" {
  type    = string
  default = "postgres://admin:admin@localhost:5432/study_db?sslmode=disable"
}

env "local" {
  src = "file://migrations/schema.hcl"
  url = var.db_url
  dev = "docker://postgres/15/dev"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}