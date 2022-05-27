.PHONY: build release

REL_VERSION=$(shell git rev-parse HEAD)
REL_BUCKET=sem-cli-releases

install.goreleaser:
	curl -L https://github.com/goreleaser/goreleaser/releases/download/v1.9.1/goreleaser_Linux_x86_64.tar.gz -o /tmp/goreleaser.tar.gz
	tar -xf /tmp/goreleaser.tar.gz -C /tmp
	sudo mv /tmp/goreleaser /usr/bin/goreleaser

go.install:
	cd /tmp
	sudo curl -O https://dl.google.com/go/go1.16.linux-amd64.tar.gz
	sudo tar -xf go1.16.linux-amd64.tar.gz
	sudo mv go /usr/local
	cd -

gsutil.configure:
	gcloud auth activate-service-account deploy-from-semaphore@semaphore2-prod.iam.gserviceaccount.com --key-file ~/gce-creds.json
	gcloud config set project semaphore2-prod

go.get:
	go get

go.fmt:
	go fmt ./...

test:
	go test -v ./...

build:
	env GOOS=$(OS) GOARCH=$(ARCH) go build -ldflags "-s -w -X cmd.VERSION=$(shell git describe --tags --abbrev=0)" -o sem
	tar -czvf /tmp/sem.tar.gz sem

# Automation of CLI tagging.

tag.major:
	git fetch --tags
	latest=$$(git tag | sort --version-sort | tail -n 1); new=$$(echo $$latest | cut -c 2- | awk -F '.' '{ print "v" $$1+1 ".0.0" }');          echo $$new; git tag $$new; git push origin $$new

tag.minor:
	git fetch --tags
	latest=$$(git tag | sort --version-sort | tail -n 1); new=$$(echo $$latest | cut -c 2- | awk -F '.' '{ print "v" $$1 "." $$2 + 1 ".0" }');  echo $$new; git tag $$new; git push origin $$new

tag.patch:
	git fetch --tags
	latest=$$(git tag | sort --version-sort | tail -n 1); new=$$(echo $$latest | cut -c 2- | awk -F '.' '{ print "v" $$1 "." $$2 "." $$3+1 }'); echo $$new; git tag $$new; git push origin $$new


#
# These two scripts update generate a new installation script based on the
# current git tag on Semaphore.

release.stable.install.script:
	sed 's/VERSION_PLACEHOLDER/$(shell git describe --tags --abbrev=0)/' scripts/get.template.sh > scripts/get.sh
	gsutil cp scripts/get.sh gs://$(REL_BUCKET)/get.sh
	gsutil acl -R ch -u AllUsers:R gs://$(REL_BUCKET)/get.sh
	gsutil setmeta -h "Cache-Control:private, max-age=0, no-transform" gs://$(REL_BUCKET)/get.sh

release.edge.install.script:
	sed 's/VERSION_PLACEHOLDER/$(shell git describe --tags --abbrev=0)/' scripts/get.template.sh > scripts/get-edge.sh
	gsutil cp scripts/get-edge.sh gs://$(REL_BUCKET)/get-edge.sh
	gsutil acl -R ch -u AllUsers:R gs://$(REL_BUCKET)/get-edge.sh
	gsutil setmeta -h "Cache-Control:private, max-age=0, no-transform" gs://$(REL_BUCKET)/get-edge.sh
	echo "https://storage.googleapis.com/$(REL_BUCKET)/get-edge.sh"
