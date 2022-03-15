# scalingo\_deployer

## Description

Dockerized app to deploy to scalingo PaaS provider.

## Usage

On semaphore CI for example

```
docker run
  -e SCALINGO_APP
  -e GITHUB_OWNER_REPO=$SEMAPHORE_GIT_REPO_SLUG
  -e GIT_REF=$SEMAPHORE_GIT_SHA
  -e GITHUB_API_TOKEN=$GITHUB_API_TOKEN
  -e SCALINGO_API_TOKEN=$SCALINGO_API_TOKEN
  gcr.io/betterplace-183212/scalingo_deployer
```
