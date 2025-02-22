FROM migrate/migrate

COPY ./ /migrations

CMD ["-path", "/migrations", "-database", "postgres://postgres:password@pg-postgresql:5432/postgres?sslmode=disable", "up"]
