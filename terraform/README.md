## Required software

* Terraform
* gcloud CLI
* Docker (with daemon running)

This setup was tested on the following versions:

```
Terraform v1.13.4
hashicorp/google v7.45.0
hashicorp/google-beta v7.45.0
```

## Setup

1. Authorize in gcloud CLI.

This projects aims for setup as easy as possible. Default application login is not recommended for production use.

```
gcloud auth login
gcloud auth application-default login
```

2. Run task. While terraform is running, you will be asked to confirm applying changes. Answer wih `yes`.

```bash
task
```

You will be asked to pick a region for Cloud Run and Firebase. If you want to use Cloud Run region different than
`europe-west1`, you need to **commit** changes in following files:

- `./Taskfile.yml` (the `deploy` task)
- `./web/firebase.json`

3. Make sure you enable `Email/Password` authentication provider in Firebase as described in the `task` output.

a. Open FireBase console: https://console.firebase.google.com
b. Choose `Wild Workouts` project
c. Go to `Authentication`
d. Choose `Sign-in method` tab
e. Click on `Email/Password`, switch to `Enabled` and click `Save`.

## Deployments

CI runs through the GitHub Actions workflow in `.github/workflows/ci.yml`: pull requests run lint,
tests, and builds only, and pushes to `master` additionally push images and deploy to Cloud Run and
Firebase Hosting. The workflow authenticates to Google Cloud with Workload Identity Federation (no service
account keys), so it needs three repository variables set. `task` prints the exact commands when the
setup finishes:

```bash
gh variable set GCP_PROJECT_ID --body "<project id>"
gh variable set GCP_WORKLOAD_IDENTITY_PROVIDER --body "$(terraform output -raw github_workload_identity_provider)"
gh variable set GCP_SERVICE_ACCOUNT --body "$(terraform output -raw github_actions_service_account)"
```

If you run this setup from a fork, set the `TF_VAR_github_repository` environment variable to your
`owner/name` before running `task`, so the Workload Identity provider trusts your repository.

## Destroy

If you want to tear down the project, run `task destroy`.

If you want to create it again, make sure to:
* Use different project name.
* Remove `terraform.tfstate` file.
