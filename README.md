# Terraform GHA Controller
Run Terraform within controlled private environments.

**Not currently suitable for use in production environments**

## Why does this exist?
To encourage invoking your IaC pipeline from a secure network.

## Features
- Tightly integrated with GitHub Checks API
- GitHub webhook secret verification
- GitHub checks API
- Supports local or remote Terraform backends

Tooling out of the box:
- helm, kubectl, git, kustomize

## Roadmap to v1

- [ ] [**Security**] Temporary persistence
- [ ] [**Feature**] Tooling extensibility (without Docker build'ing)
- [ ] [**Security**] Rootless container (or fakeroot?)
- [ ] [**Feature**] Multiple GitHub repository support
- [ ] [**Tests**] More please

#### v2 ideas
- [ ] Convert git adapters into adapter pattern (GitHub being the first) 


## Setting up
1. GitHub Application 
   - [Contents] Private key
   - [Meta] Installation ID
   - [Meta] Application ID
   - [Permission] Checks (Write)
   - [Permission] Contents (Read)
   
### Setup environment
```bash
export APP_ID=123456
export INSTALLATION_ID=99999
export GH_SECRET=12345
export PRIVATE_KEY=github-app.pem
export GH_OWNER=github-user-or-org
```

Set Terraform variables
```bash
TFVARS="/config/good.auto.tfvars"
```

### Running server
```
go run ./cmd/server
```
or Docker
```
docker run -p 8080:8080 \
    -e APP_ID=123456 \
    -e INSTALLATION_ID=12345678 \ 
    -e GH_SECRET=65 \ 
    -e PRIVATE_KEY=/config/github-app.key \
    -e GH_OWNER=organisation
    -e TF_CLI_CONFIG_FILE=/.terraform/.terraformrc \
    -v "/home/user/github/github-app.key:/config/github-app.key" \
    -v "/home/user/.:/config/github-app.key" \
    ghcr.io/dev-this/tf-gha-orch:latest
```

### Submit Request
`POST /webhook` GitHub webhook payload
- Ensure `x-hub-signature` request header is set

