
build-dev:
	goreleaser build --single-target --snapshot --rm-dist

build:
	goreleaser build --rm-dist

release-and-publish:
	goreleaser release --rm-dist

release-no-publish:
	goreleaser release --skip-publish --rm-dist
