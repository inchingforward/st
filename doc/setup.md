# Postgres setup

create role storytellers createdb login;
create database storytellers owner storytellers;

# Migrations

Migrations are done using [migrate](https://github.com/golang-migrate/migrate).

## Installation

See the [docs](https://github.com/golang-migrate/migrate/tree/master/cli) for the migrate cli tool.

## Creating a new migration

From the project's root directory:

    migrate create -dir ./migrations -ext sql initial_tables

## Migrating

From the project's root directory:

    migrate -database postgresql://storytellers@localhost:5432/storytellers?sslmode=disable -path ./migrations up

# General

## Running

- Use the -debug flag when running in a non-production environment.  This will disable template reloading, analytics, etc.
