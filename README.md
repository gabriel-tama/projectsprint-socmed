# Week 2 - Social Media App

## Documents

Dokumen pendukung untuk project untuk project ini bisa dilihat di link dibawah ini.

- [Notion](https://openidea-projectsprint.notion.site/ProjectSprint-2-Social-Media-App-ff9e8ce81afb4310bdc0e38f8ad8ab67)

## Setup Project

TL;DR:

```bash
cp .env.example .env
make docker-up
make migrate-up
```

LONG VERSION:

### env

Jangan lupa populate file `.env` dengan file `.env.example` yang sudah ada.

```bash
cp .env.example .env
```

### Database

Database bisa di run dengan menggunakan `docker-compose` dengan perintah dibawah ini. Docker ini akan menjalankan database `postgres` dan `adminer` untuk management database. Database akan di mapping di port `5432` dan adminer pada `localhost:8080`.

```bash
make docker-up
```

### Migration

Tools untuk migrasi database, dengan menggunakan [migrate](https://github.com/golang-migrate/migrate/tree/master) dan `makefile`. Untuk migrasi database dengan perintah dibawah ini, nanti file sql akan ada di folder `db/migrations`.

#### create migration

```bash
make migrate-create filename.sql
```

#### run migration

```bash
make migrate-up
```

### Optional

#### Hot Reload

Kita bisa pake `air` untuk hot reload, dengan perintah dibawah ini.

```bash
go install github.com/cosmtrek/air@latest
alias air='$(go env GOPATH)/bin/air'
air
```
