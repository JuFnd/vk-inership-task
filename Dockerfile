FROM golang:latest

WORKDIR /app

COPY /configs /app/configs

COPY /cmd/authorization/main.go /app/cmd/authorization/main.go
COPY /cmd/films/main.go /app/cmd/films/main.go

RUN apt-get update && apt-get install -y redis-server

RUN apt-get install -y postgresql

RUN service postgresql start && \
    psql -c "CREATE USER boss WITH superuser login password 'boss';" && \
    psql -c "ALTER ROLE boss WITH PASSWORD 'boss';" && \
    createdb -O boss auth_service && \
    createdb -O boss films_service

COPY /database/auth_service_migrations.sql /app/database/auth_service_migrations.sql
COPY /database/films_service_migrations.sql /app/database/films_service_migrations.sql

RUN service postgresql start && \
    psql -U boss -d auth_service -f /app/database/auth_service_migrations.sql && \
    psql -U boss -d films_service -f /app/database/films_service_migrations.sql

# Build the Go binaries
RUN go build -o /app/authorization /app/cmd/authorization/main.go
RUN go build -o /app/films /app/cmd/films/main.go

# Expose the necessary ports
EXPOSE 6379
EXPOSE 5432
EXPOSE 50051
EXPOSE 8080
EXPOSE 8081

# Start the Redis, PostgreSQL, and gRPC services with nohup
CMD service redis-server start && service postgresql start && \
    nohup /app/authorization > /dev/null 2>&1 & && \
    nohup /app/films > /dev/null 2>&1 &
