# postgres

The PostgreSQL instance shared by the services on this host.

It was originally a service inside the urlstat compose project, which tied
every other consumer to that project's network and lifecycle. It lives here
so consumers are equal clients of shared infrastructure.

## Consumers

| Service | Database | Role |
| --- | --- | --- |
| [urlstat](https://github.com/changkun/urlstat) | `urlstat` | `urlstat` |

A consumer joins the `postgres_internal` network and connects to host
`postgres` on port 5432. The instance publishes no host port and is not on
`traefik_proxy`, so nothing reaches it from the ingress path.

## Running

```sh
$ cp .env.template .env   # then fill in POSTGRES_PASSWORD
$ docker compose up -d
```

## Schema changes

The image only runs `/docker-entrypoint-initdb.d` when the data directory
is empty, which it no longer is. Apply DDL directly instead:

```sh
$ docker exec -i postgres psql -U urlstat -d urlstat < change.sql
```

## Adding a consumer

Create a database and role for it rather than sharing the `urlstat` one:

```sh
$ docker exec -i postgres psql -U urlstat -d postgres <<'SQL'
CREATE ROLE myservice LOGIN PASSWORD '...';
CREATE DATABASE myservice OWNER myservice;
SQL
```

Then declare `postgres_internal` as an external network in that service's
compose file and attach the service to it.
