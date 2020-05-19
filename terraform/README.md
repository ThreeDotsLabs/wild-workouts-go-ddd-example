## Required software

* Terraform
* gcloud CLI
* Docker (with daemon running)

This setup was tested on the following versions:

```
Terraform v0.12.24

Google Cloud SDK 290.0.1
alpha 2019.05.17
beta 2019.05.17
core 2020.04.24
```

## Setup

1. Authorize in gcloud CLI.

This projects aims for setup as easy as possible. Default application login is not recommended for production use.

```
gcloud auth login
gcloud auth application-default login
```

2. Run make. While terraform is running, you will be asked to confirm applying changes. Answer wih `yes`.

```bash
make
```

You will be asked to pick a region for Cloud Run and Firebase. If you want to use Cloud Run region different than
`europe-west1`, you need to **commit** changes in following files:

- `./scripts/deploy.sh`
- `./web/firebase.json`

3. Make sure you enable `Email/Password` authentication provider in Firebase as described in the `make` output.

a. Open FireBase console: https://console.firebase.google.com
b. Choose `Wild Workouts` project
c. Go to `Authentication`
d. Choose `Sign-in method` tab
e. Click on `Email/Password`, switch to `Enabled` and click `Save`.

## Cloud builds

Go to https://console.cloud.google.com/cloud-build/builds to see your recent builds.

## Destroy

If you want to tear down the project, run `make destroy`.

If you want to create it again, make sure to:
* Use different project name.
* Remove `terraform.tfstate` file.
