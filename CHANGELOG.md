# Changelog

### v0.8.14

- Issues with release machinery

### v0.8.13

- List pipelines
- Workflow stop
- Adjust message when editor fails to edit
- archive switch in snapshot cretate (can send arbitrary tgz file as snapshot archive)
- Workflow rebuild

### v0.8.12

- Describe workflow and extended describe for workflow
- Stop running job with `sem stop job <id>`

### v0.8.11

- Default duration of Debug session is 60minutes and can be extended on demand
  with `sem debug project <name> --duration 30m`.
- Require at least one argument in create commands

### v0.8.10

- Fix bug with missing verification of the SSH key while starting debugging
- Create workflow with all files in current directory
- Create default label for snapshot
- Display creation time in get workflow
- Resolve project name based on git URL

### v0.8.9

- Fix bug when sem CLI panics in case of erroneous YAML file

### v0.8.8

- Fix osX binary release name from 'cli' -> 'sem'

### v0.8.7

- Goreleaser and Homebrew

### v0.8.6

- Apply changes to a notification with sem apply -f [notification-file]

### v0.8.5

- Create workflow from label
- Create notification from YAML file
- Edit notification
- Fix notification creation with no branch rule

### v0.8.4

- Fix creation of Slack notifications from `sem create notification`

### v0.8.3

- Get, List, Create and Delete notifications
- Extend pipeline describe with block/job info
- Put archive in snapshot create request

### v0.8.2

- Define machine-type when starting a debug session for projects.
  Example: `sem debug project cli --machine-type e1-standard-4`.

### v0.8.1

- Fix debugging projects with sem CLI tool

### v0.8.0

- Describe a pipeline
- Stop a pipeline
- List workflows
- Rebuild a pipeline
- Snapshot based workflow
- HTTP(s) support in sem init

### v0.7.4

- Create a job from a YAML file.

### v0.7.3

- Debug a job. Example: `sem debug job <existing-job-id>`

### v0.7.2

- Secrets creation from files. Example: `sem create secret aws-secrets -f /home/vagrant/.aws/credentials:/home/semaphore/.aws/credentials`
- On demand job creation. Example: `sem create job hello-world --project cli --command 'echo "Hello World"'`

### v0.7.1

- SSH command error message is displayed if `sem attach` fails. Example: `permission deinied`.
- Debugging a project with `sem debug project <project-name>`

### v0.7.0

- Listing jobs
- Describing jobs
- Displaying logs with `sem logs <job-id>`
- Attaching to running jobs with `sem attach`
- Port forwarding for jobs with `sem port-forward`

### v0.6.1

- Dashboards: Get, List, Edit, Create, Update
- Removal of `sem describe` in favor of `sem get <kind> <name>`
- Projects are listing repository URL when `sem get projects`
- Test coverage for all commands

### v0.5.0

- Fix invalid .semaphore/semaphore.yml creation logic introduced in preview1
- Secret name presence is validates in create and update
- Both `sem create secret <name>` `sem create secrets <name>` is now available
  for consistency with other actions
- Prevent raising errors if `sem create invalidname` is executed. Now a list of
  subcommands is displayed.
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
