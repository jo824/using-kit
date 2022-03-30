GO_SERVER_PORT=8008:8008
APP_NAME=kit-server

## this may not be how you want to install your own deps
## and could use modules instead .. see
## https://go.dev/blog/using-go-modules
## https://github.com/golang/go/issues/38812
deps:
	export GO111MODULE=off
	go get ./...
	go install ./...


dbuild: gobc
	docker build --no-cache -t $(APP_NAME) -f using-kit/Dockerfile .

drun:
	docker run -d  --name $(APP_NAME) -p $(GO_SERVER_PORT) kit-server

dstop:  ## Stop and remove a running container
	docker stop $(APP_NAME)

drm: dstop ## Stop and remove a running container
	docker rm $(APP_NAME)

dclean: drm
	docker image rm $(APP_NAME)

gob:  ## build binary for local
	go build -o ./using-kit/$(APP_NAME) ./using-kit/main.go

gobc: ## build bin for container
	GOOS=linux GOARCH=amd64 go build -o ./using-kit/$(APP_NAME)-d ./using-kit/main.go

