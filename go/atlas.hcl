env "local" {
  src = "file://schema.sql"

  dev = "postgres://admin:admin@localhost:5432/app?sslmode=disable"

  migration {
    dir    = "file://migrations"
    format = "atlas"
  }

  url = "postgres://admin:admin@localhost:5432/app?sslmode=disable"

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}