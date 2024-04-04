IMG := "nikogura/configsvc"
TAG := "0.1.0"

.PHONY: build
build:
	docker build -t $(IMG):$(TAG) .

.PHONY: push
push:
	docker push $(IMG):$(TAG)
