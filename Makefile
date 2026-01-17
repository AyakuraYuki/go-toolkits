.PHONY: build-synology-webhook-proxy

build-synology-webhook-proxy:
	@docker buildx build --platform=linux/amd64 -t synology-webhook-proxy:1.0.0 -f docker/synology-webhook-proxy.dockerfile .
