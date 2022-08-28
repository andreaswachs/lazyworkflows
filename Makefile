.PHONY: build-unit-tests
build-unit-tests:
	docker build -t lazyworkflows-unit-tests:latest -f Dockerfiles/Dockerfile.test .

.PHONY: unit-tests
unit-tests:
	docker run -it -v $(shell pwd):/app lazyworkflows-unit-tests:latest

.PHONY: unit-tests-interactive
unit-tests-continuous:
	docker run -it -v $(shell pwd):/app lazyworkflows-unit-tests:latest bash