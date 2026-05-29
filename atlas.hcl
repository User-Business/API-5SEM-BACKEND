env "local" {
  src = "file://schema.sql"

  dev = "postgres://denarius:denarius@localhost:5433/denarius?sslmode=disable&search_path=public"

  url = "postgres://denarius:denarius@localhost:5433/denarius?sslmode=disable&search_path=public"
}