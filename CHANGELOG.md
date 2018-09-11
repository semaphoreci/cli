# Changelog

### v0.5.0-preview2

- Fix invalid .semaphore/semaphore.yml creation logic introduced in preview1
- Secret name presence is validates in create and update
- Both `sem create secret <name>` `sem create secrets <name>` is now available
  for consistency with other actions
- Prevent raising errors if `sem create invalidname` is executed. Now a list of
  subcommands is displayed.

### v0.5.0-preview1

- Ability to pass `--project-name` and `--repo-url` to `sem init`
- Secrets are using the `v1beta` API
- Files in Secrets
- Edit a secret with `sem edit secret <name>`
- Update a secret with `sem apply -f <file-path>`
- Secrets YAML validation in the CLI
- Display age of secrets when listing with `sem get secrets`
- ApiVersion and Kind always present in Yaml format of Secrets
- Get one resource with `sem get project <name>`, `sem get secret <name>`
- Create an empty secret with `sem create secret <name>`

### v0.4.2

- Verbose sem init git repository erorrs
- Fixed wrong .semaphore/semaphore.yml file permission. Now it is 0644.
- Type fix in sem config
