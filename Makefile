.PHONY: build release

REL_VERSION=$(shell git rev-parse HEAD)
REL_BUCKET=sem-cli-releases

go.install:
	cd /tmp
	sudo curl -O https://storage.googleapis.com/golang/go1.9.1.linux-amd64.tar.gz
	sudo tar -xvf go1.9.1.linux-amd64.tar.gz
	sudo mv go /usr/local
	export PATH=$PATH:/usr/local/go/bin
	cd -

gsutil.configure:
	gcloud auth activate-service-account deploy-from-semaphore@semaphore2-prod.iam.gserviceaccount.com --key-file ~/gce-creds.json
	gcloud config set project semaphore2-prod
	gcloud container clusters get-credentials prod --zone us-east4

go.get:
	go get -t -d -v ./... && go build -v ./...

go.fmt:
	go fmt ./...

test:
	go test -v ./...

build:
	env GOOS=$(OS) GOARCH=$(ARCH) go build -o sem
	tar -czvf /tmp/sem.tar.gz sem

release:
	$(MAKE) build OS=$(OS) ARCH=$(ARCH)
	gsutil cp /tmp/sem.tar.gz gs://$(REL_BUCKET)/$(REL_VERSION)-$(OS)-$(ARCH).tar.gz
	gsutil acl -R ch -u AllUsers:R gs://$(REL_BUCKET)/$(REL_VERSION)-$(OS)-$(ARCH).tar.gz
	gsutil setmeta -h "Cache-Control:private, max-age=0, no-transform" gs://$(REL_BUCKET)/$(REL_VERSION)-$(OS)-$(ARCH).tar.gz
	echo "https://storage.googleapis.com/$(REL_BUCKET)/$(REL_VERSION)-$(OS)-$(ARCH).tar.gz"

release.all:
	$(MAKE) release OS=linux   ARCH=386
	$(MAKE) release OS=linux   ARCH=amd64
	$(MAKE) release OS=darwin  ARCH=386
	$(MAKE) release OS=darwin  ARCH=amd64
	# $(MAKE) release OS=windows ARCH=386    # mousetrap issues?
	# $(MAKE) release OS=windows ARCH=amd64

release.stable:
	$(MAKE) release.all REL_VERSION=stable

release.edge:
	$(MAKE) release.all REL_VERSION=edge

release.stable.install.script:
	gsutil cp scripts/get gs://$(REL_BUCKET)/get.sh
	gsutil acl -R ch -u AllUsers:R gs://$(REL_BUCKET)/get.sh
	gsutil setmeta -h "Cache-Control:private, max-age=0, no-transform" gs://$(REL_BUCKET)/get.sh
	echo "https://storage.googleapis.com/$(REL_BUCKET)/get.sh"

release.edge.install.script:
	gsutil cp scripts/get-edge gs://$(REL_BUCKET)/get-edge.sh
	gsutil acl -R ch -u AllUsers:R gs://$(REL_BUCKET)/get-edge.sh
	gsutil setmeta -h "Cache-Control:private, max-age=0, no-transform" gs://$(REL_BUCKET)/get-edge.sh
	echo "https://storage.googleapis.com/$(REL_BUCKET)/get-edge.sh"
