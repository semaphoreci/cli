.PHONY: build release

BRANCH=$(shell echo "$(BRANCH_NAME)" | sed 's/[^a-z]//g')
REL_BUCKET="gs://sem-cli-releases"
REL_VERSION="$(BRANCH)-$(SEMAPHORE_BUILD_NUMBER)-sha-$(REVISION)"

build:
	go build

release:
	gsutil cp sem $(REL_BUCKET)/$(REL_VERSION)
	gsutil acl -R ch -u AllUsers:R $(REL_BUCKET)/$(REL_VERSION)
