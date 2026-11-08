# infra

The infrastructure of changkun.de, golang.design, etc.: the traefik proxy and
the shared services behind it, each a Docker Compose stack in its own folder.

```bash
traefik/    # the proxy: routes each address to its service, holds the certificates
postgres/   # the PostgreSQL instance the services share
dns/        # the DNS records on Cloudflare, applied with DNSControl
dev/        # traefik for a local run on localhost
```

## Services

changkun.de and golang.design are served by
[changkun/website](https://github.com/changkun/website) from the folders
`/www/changkun.de` and `/www/golang.design`, into which most of the
repositories below publish. A folder is live at its path once it is there.

Some of the addresses, and where they come from:

```bash
changkun.de                    # github.com/changkun/website
changkun.de/{s,x}              # github.com/changkun/redir
changkun.de/blog               # github.com/changkun/blog
changkun.de/research           # github.com/changkun/research
changkun.de/modern-cpp         # github.com/changkun/modern-cpp-tutorial
changkun.de/midgard            # github.com/changkun/midgard
changkun.de/urlstat            # github.com/changkun/urlstat
golang.design                  # github.com/golang-design/website
golang.design/{s,x}            # github.com/changkun/redir
golang.design/gossa            # github.com/golang-design/ssaplayground
golang.design/research         # github.com/golang-design/research
golang.design/under-the-hood   # github.com/golang-design/under-the-hood
golang.design/go-questions     # github.com/golang-design/go-questions
poly.red                       # github.com/changkun/redir (the import path of changkun/polyred)
poly.red/study                 # github.com/changkun/rankui-v2
```

## DNS

The records of the zones on Cloudflare are in
[`dns/dnsconfig.js`](dns/dnsconfig.js). Edit the file, then apply it with
[DNSControl](https://dnscontrol.org) from `dns/`:

```bash
export CLOUDFLARE_ACCOUNT_ID=...
export CLOUDFLARE_API_TOKEN=...  # All Domains: Zone Read; DNS, Dynamic URL Redirects: Edit
dnscontrol preview               # what a push would change
dnscontrol push
```

A push deletes records that are not in the file, so read the preview
first. To add a zone, write its current records out with
`dnscontrol get-zones --format=js cloudflare <zone>` and paste them in.

## Rebuild

The code of every service is in git. What lives only on the server, and in
DigitalOcean's weekly backup of it:

```bash
/root/changkun.de/*/.env                         # each service's secrets, as listed in its .env.template
/root/changkun.de/infra/postgres/.env
/root/changkun.de/rankui-v2/study/app/.env
/root/changkun.de/infra/postgres/data/           # the database of redir and urlstat
/root/changkun.de/infra/traefik/conf/acme*.json  # certificates; traefik requests new ones if missing
/root/changkun.de/redir/data/redirconf.yml
/root/changkun.de/midgard/config.yml
/root/changkun.de/midgard/data/
/root/changkun.de/website/data/
/root/changkun.de/rankui-study-data/
/root/golang.design/ssaplayground/data/
/www/                                            # published into by the repositories' workflows
/root/.ssh/authorized_keys                       # includes the keys those workflows deploy with
```

From the backup, which is quicker: create a droplet from the latest backup
in DigitalOcean. Every container has `restart: always` and comes back on
its own.

From scratch:

1. Create an Ubuntu droplet, install Docker with the compose plugin and Go
   (the website's build runs its tests), and add a 2 GB swap file.
2. Restore `/root/.ssh/authorized_keys`. With new keys instead, update the
   `SERVER_KEY` secret of each repository that publishes into `/www`.
3. Clone infra, website, redir, urlstat, midgard and rankui-v2 into
   `/root/changkun.de`, and golang-design/ssaplayground into
   `/root/golang.design`.
4. Restore the files above.
5. Start traefik and postgres first, since they create the `traefik_proxy`
   and `postgres_internal` networks the other services join:

   ```bash
   cd /root/changkun.de/infra/traefik && docker compose up -d
   cd /root/changkun.de/infra/postgres && docker compose up -d
   cd /root/changkun.de/website && make build up
   cd /root/changkun.de/redir && make build up
   cd /root/changkun.de/urlstat && make build up
   cd /root/changkun.de/midgard && make build && docker compose up -d
   cd /root/changkun.de/rankui-v2/study/app && make build up
   cd /root/golang.design/ssaplayground && make build up
   ```

6. Without a backup of `/www`, rerun the publishing workflows.

Either way, if the address changed, set `SERVER` in `dns/dnsconfig.js` and
push it, change the records of changkun.de, and update the `SERVER_DOMAIN`
secrets that hold an address rather than a name. The
[status page](https://status.changkun.de/) shows when everything is
back.

Start a compose project only from its folder here. A container keeps the
path it was created with, so after a folder is moved or renamed, run
`docker compose up -d` in its new place, or the next restart mounts an
empty folder at the old path. After a restore too, this prints any mount
whose folder is gone:

```bash
for s in $(docker inspect -f '{{range .Mounts}}{{.Source}} {{end}}' $(docker ps -q)); do [ -e "$s" ] || echo "missing: $s"; done
```

## License

GPLv3 | Copyright &copy; 2021 [Changkun Ou](https://changkun.de)
