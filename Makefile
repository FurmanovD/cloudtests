SHELL := /bin/bash

export APP_NAME=usersvc
export APP_VERSION=$(shell cat ./VERSION)
export GIT_COMMIT=$$(git log -1 --format="%H")
export BUILD_TIME=$$(date -u "+%F_%T")
export APP_BACKEND_IMAGE=${APP_NAME}-backend:${APP_VERSION}
export APP_BACKEND_IMAGE_LATEST=${APP_NAME}-backend:latest

build-image:
	# an app image will be built one time only
	@echo "check and build if not exist an image for the ${APP_NAME} "
	@docker image inspect ${APP_BACKEND_IMAGE} > /dev/null || docker build \
		--build-arg APP_NAME=${APP_NAME} \
		--build-arg BUILD_TIME=${BUILD_TIME} \
		--build-arg GIT_COMMIT=${GIT_COMMIT} \
		--build-arg GOOS_TYPE=linux \
		--tag ${APP_BACKEND_IMAGE} \
		--tag ${APP_BACKEND_IMAGE_LATEST} \
		--file build/docker/Dockerfile .

rebuild-image:
	# an app image will be built
	@echo "building an app ${APP_NAME} image"
	@docker build \
		--build-arg APP_NAME=${APP_NAME} \
		--build-arg BUILD_TIME=${BUILD_TIME} \
		--build-arg GIT_COMMIT=${GIT_COMMIT} \
		--build-arg GOOS_TYPE=linux \
		--tag ${APP_BACKEND_IMAGE} \
		--tag ${APP_BACKEND_IMAGE_LATEST} \
		--file build/docker/Dockerfile .

build: rebuild-image
	# builds a binary file
	@echo "building ${APP_NAME} (commit:${GIT_COMMIT})"
	@CID=$$(docker create ${APP_BACKEND_IMAGE}) && \
	docker cp $${CID}:/${APP_NAME} ${APP_NAME} && \
	docker rm $${CID}

test_docker: test_deps_stop stop rebuild-image
	docker-compose -f docker-compose.yaml -f docker-compose.tests.yaml up -d tests && docker logs -f cap_tests_1 && docker-compose -f docker-compose.yaml -f docker-compose.tests.yaml down

test:
	@go test ./... -p=1

test_deps_run: run-deps

test_deps_stop:
	docker-compose -f docker-compose.yaml down

test_deps_restart: test_deps_stop test_deps_run

run: rebuild-image
	docker-compose -f docker-compose.yaml up -d

run-deps:
	docker-compose -f docker-compose.yaml up -d redis
	
stop:
	docker-compose -f docker-compose.yaml down

clean:
	rm ${APP_NAME}
