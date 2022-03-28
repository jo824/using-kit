GO_SERVER_PORT=8008:8008
APP_NAME=kit-server



dkbuild:
	docker build --no-cache -t kit-server  -f using-kit/Dockerfile .

dkrun:
	docker run -d  --name $(APP_NAME) -p $(GO_SERVER_PORT) kit-server

dkstop:  ## Stop and remove a running container
	docker stop $(APP_NAME)

dkrm: dkstop ## Stop and remove a running container
	docker rm $(APP_NAME)


