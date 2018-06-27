.PHONY: build release

REL_VERSION=$(shell git rev-parse HEAD)
REL_BUCKET=sem-cli-releases

go.get:
	go get -t -d -v ./... && go build -v ./...

test:
	go test -v ./...

build:
	go build
	tar -czvf /tmp/sem.tar.gz sem

build-all: go.get
	GOOS=darwin GOARCH=arm64 go build

gsutil.configure:
	./scripts/install-gcloud
	gcloud auth activate-service-account deploy-from-semaphore@semaphore2-prod.iam.gserviceaccount.com --key-file ~/semaphore2-prod-2fd29ae99af8.json
	gcloud config set project semaphore2-prod
	gcloud container clusters get-credentials prod --zone us-east4

release: build
	gsutil cp /tmp/sem.tar.gz gs://$(REL_BUCKET)/$(REL_VERSION)
	gsutil acl -R ch -u AllUsers:R gs://$(REL_BUCKET)/$(REL_VERSION)
	echo "https://storage.googleapis.com/$(REL_BUCKET)/$(REL_VERSION)"
