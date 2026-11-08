# infra

The infrastructure of changkun.de, golang.design, etc.: the traefik proxy and
the shared services behind it, each a Docker Compose stack in its own folder.

```bash
traefik/    # the proxy: routes each address to its service, holds the certificates
postgres/   # the PostgreSQL instance the services share
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
changkun.de/upbot              # github.com/changkun/upbot
golang.design                  # github.com/golang-design/website
golang.design/{s,x}            # github.com/changkun/redir
golang.design/gossa            # github.com/golang-design/ssaplayground
golang.design/research         # github.com/golang-design/research
golang.design/under-the-hood   # github.com/golang-design/under-the-hood
golang.design/go-questions     # github.com/golang-design/go-questions
poly.red                       # github.com/changkun/redir (the import path of changkun/polyred)
poly.red/study                 # github.com/changkun/rankui-v2
```

## License

GPLv3 | Copyright &copy; 2021 [Changkun Ou](https://changkun.de)
