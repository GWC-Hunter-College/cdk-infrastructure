# Edited from https://docs.aws.amazon.com/lambda/latest/dg/go-image.html, 
# 'Using an AWS OS-only base image -> creating an image from the provided.al2023 base image' section #4.
FROM golang:1.24.4-alpine3.21 AS build
WORKDIR /init

COPY go.mod go.sum ./

COPY main.go .
RUN go build -o main main.go

# Copy artifacts to a clean image
FROM alpine:3.21
COPY --from=build /init/main /main

# Copy migration files
COPY /migrations/07_11_2025_create_databases_up.sql /migrations/07_11_2025_create_databases_up.sql
COPY /migrations/07_11_2025_create_core_tables_up.sql /migrations/07_11_2025_create_core_tables_up.sql
COPY /migrations/07_11_2025_create_core_tables_down.sql /migrations/07_11_2025_create_core_tables_down.sql
COPY /migrations/07_11_2025_create_member_form_migration_table_up.sql /migrations/07_11_2025_create_member_form_migration_table_up.sql
COPY /migrations/07_11_2025_create_member_form_migration_table_down.sql /migrations/07_11_2025_create_member_form_migration_table_down.sql

ENTRYPOINT [ "/main" ]