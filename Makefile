IMAGE = h-platform-automation-cc-server
VERSION = 0.1

.PHONY: _
_: build publish

.PHONY: build
build:
	docker build -t $(IMAGE):$(VERSION) .

.PHONY: publish
publish:
	docker push $(IMAGE):$(VERSION)
