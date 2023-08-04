.PHONY: build release

REL_VERSION=$(shell git rev-parse HEAD)
REL_BUCKET=sem-cli-releases

install.goreleaser:
	curl -L https://github.com/goreleaser/goreleaser/releases/download/v1.14.1/goreleaser_Linux_x86_64.tar.gz -o /tmp/goreleaser.tar.gz
	tar -xf /tmp/goreleaser.tar.gz -C /tmp
	sudo mv /tmp/goreleaser /usr/bin/goreleaser

gsutil.configure:
	gcloud auth activate-service-account $(GCP_REGISTRY_WRITER_EMAIL) --key-file ~/gce-registry-writer-key.json
	gcloud --quiet auth configure-docker
	gcloud --quiet config set project semaphore2-prod

go.get:
	docker-compose run --rm cli go get

go.fmt:
	docker-compose run --rm cli go fmt ./...

test:
	docker-compose run --rm cli gotestsum --format short-verbose --junitfile results.xml --packages="./..." -- -p 1

build:
	docker-compose run --rm cli env GOOS=$(OS) GOARCH=$(ARCH) go build -ldflags "-s -w -X cmd.VERSION=$(shell git describe --tags --abbrev=0)" -o sem
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
