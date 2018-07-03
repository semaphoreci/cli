# Sem

## Install

Edge (latest build on master branch):

``` bash
curl https://storage.googleapis.com/sem-cli-releases/get-edge.sh | bash
```

Stable (latest stable version, manually released):

``` bash
curl https://storage.googleapis.com/sem-cli-releases/get.sh | bash
```

## Configure

``` bash
sem config set auth.token <token>
sem config set host <org>.semaphoreci.com
```

## Initialize a project

In the root of the repository:

``` bash
sem init
```

## Low level interface

## Projects

### Create project

``` yaml
# project.yaml

apiVersion: v1alpha
kind: Project
metadata:
  name: test
spec:
  repository:
    url: "git@github.com:<owner>/<name>.git"
```

``` bash
sem create -f project.yml
```

### List projects

``` bash
sem get projects
```

### Describe projects

``` bash
sem describe projects <name>
```

## Secrets (not yet supported)

### Create secret

``` yaml
# aws-secrets.yaml

metadata:
  name: "aws-secrets"
data:
  env_vars:
    - name: "aws-id"
      value: "123"
	  - name: "aws-secret"
		  value: "$ekret"
```

``` bash
sem create -f aws-secrets.yml
```

### List secrets

``` bash
sem get secrets
```

### Describe secrets

``` bash
sem describe secrets <name>
```

### Delete secrets

``` bash
sem delete secrets <name>
```

## TODOs

- [ ] Automate ~/.sem.yaml creation
- [ ] Sem config set validation
- [ ] Releaes v1.0
- [ ] Expose env vars in secrets
- [ ] Expose config files in secrets
- [ ] Update resources via CLI
