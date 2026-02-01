# Copyright 2020 Changkun Ou. All rights reserved.
# Use of this source code is governed by a GPL-3.0
# license that can be found in the LICENSE file.
all: build
build:
	CGO_ENABLED=0 GOOS=linux go build
	docker buildx build -t static:latest --load .
up:
	docker compose up -d
down:
	docker compose down