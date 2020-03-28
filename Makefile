DOCKER_IMAGE=nonlogical/smt.logger

.PHONY: install
install:
	go install ./...

.PHONY: docker.build
docker.build:
	docker build -t $(DOCKER_IMAGE) .

.PHONY: docker.run
docker.run:
	docker run -it --rm --name smt.logger $(DOCKER_IMAGE)

.PHONY: docker.upload
docker.upload: docker.build
	docker push $(DOCKER_IMAGE):latest

docker.mpush:
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t $(DOCKER_IMAGE):latest --push .

