prepare:
	go mod download

run:
	go build -o bin/main cmd/chat/main.go
	./bin/main

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o bin/main cmd/chat/main.go
	chmod +x bin/main

dkb:
	docker build -t superhero-chat .

dkr:
	docker run -p "8120:8120" superhero-chat

launch: dkb dkr

chat-log:
	docker logs superhero-chat -f

db-log:
	docker logs db -f

rabbitmq-log:
	docker logs rabbitmq -f

rmc:
	docker rm -f $$(docker ps -a -q)

rmi:
	docker rmi -f $$(docker images -a -q)

clear: rmc rmi

chat-ssh:
	docker exec -it superhero-chat /bin/bash

db-ssh:
	docker exec -it db /bin/bash

rabbitmq-ssh:
	docker exec -it rabbitmq /bin/bash

PHONY: prepare build dkb dkr launch chat-log db-log rabbitmq-log chat-ssh db-ssh rabbitmq-ssh rmc rmi clear