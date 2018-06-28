.PHONY: build release

REL_VERSION=$(shell git rev-parse HEAD)
REL_BUCKET=sem-cli-releases

OS=linux
ARCH=amd64

go.get:
	go get -t -d -v ./... && go build -v ./...

go.fmt:
	go fmt ./...

test:
	go test -v ./...

build:
	go build
	tar -czvf /tmp/sem.tar.gz sem

gsutil.configure:
	./scripts/install-gcloud
	gcloud auth activate-service-account deploy-from-semaphore@semaphore2-prod.iam.gserviceaccount.com --key-file ~/semaphore2-prod-2fd29ae99af8.json
	gcloud config set project semaphore2-prod
	gcloud container clusters get-credentials prod --zone us-east4

release: build
	gsutil cp /tmp/sem.tar.gz gs://$(REL_BUCKET)/$(REL_VERSION)-$(OS)-$(ARCH).tar.gz
	gsutil acl -R ch -u AllUsers:R gs://$(REL_BUCKET)/$(REL_VERSION)-$(OS)-$(ARCH).tar.gz
	gsutil setmeta -h "Cache-Control:private, max-age=0, no-transform" gs://$(REL_BUCKET)/$(REL_VERSION)-$(OS)-$(ARCH).tar.gz
	echo "https://storage.googleapis.com/$(REL_BUCKET)/$(REL_VERSION)-$(OS)-$(ARCH).tar.gz"

release.stable:
	$(MAKE) release REL_VERSION=stable

release.install.script:
	gsutil cp scripts/get gs://$(REL_BUCKET)/get.sh
	gsutil acl -R ch -u AllUsers:R gs://$(REL_BUCKET)/get.sh
	gsutil setmeta -h "Cache-Control:private, max-age=0, no-transform" gs://$(REL_BUCKET)/get.sh
	echo "https://storage.googleapis.com/$(REL_BUCKET)/get.sh"
