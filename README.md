# Sem

Semaphore 2.0 command line interface.

## Install

Edge (latest build on master branch):

``` bash
curl https://storage.googleapis.com/sem-cli-releases/get-edge.sh | bash
```

Stable (latest stable version, manually released):

``` bash
curl https://storage.googleapis.com/sem-cli-releases/get.sh | bash
```

Homebrew (latest stable version)

```bash
brew install semaphoreci/tap/sem
```

## Development

### Releases

We build a new release for every tag in this repository and upload it to Github.

Apart from this, we have two installation scripts:
 - `get.sh` - gets the latest stable version of the CLI
 - `get-edge.sh` - gets the latest edge version of the CLI

The `edge` script is updated every time we release a tag. The `stable` `get.sh`
script needs to be manually approved and released from Semaphore. Follow the
releasing new versions procedure to update and release.

### Releasing new versions

1. Prepare the changes as a PR -> Get a green build -> Merge to master

2. Checkout latest master code and run: `make tag.patch`, `make tag.minor` or
  `make tag.major` depending on the type of the change. We use semantic
   versioning for the CLI.

3. Semaphore will build the release, upload it to Github and our brew taps, and
   update the `get-edge` installation script.

4. The `stable` installation script needs to be updated by manually promoting
   the release on Semaphore. Find the workflow for the tag you want to promote
   to stable, and click on the "Stable" promotion. This will update the `get.sh`
   script.

5. Update the CHANGELOG.md file
