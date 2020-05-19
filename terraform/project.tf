provider "google" {
  project = var.project
  region  = var.region
}

data "google_billing_account" "account" {
  display_name = var.billing_account
}

resource "google_project" "project" {
  name            = "Wild Workouts"
  project_id      = var.project
  billing_account = data.google_billing_account.account.id
}

resource "google_project_iam_member" "owner" {
  role   = "roles/owner"
  member = "user:${var.user}"

  depends_on = [google_project.project]
}

resource "google_project_service" "compute" {
  service    = "compute.googleapis.com"
  depends_on = [google_project.project]
}

resource "google_project_service" "container_registry" {
  service    = "containerregistry.googleapis.com"
  depends_on = [google_project.project]

  disable_dependent_services = true
}

resource "google_project_service" "cloud_run" {
  service    = "run.googleapis.com"
  depends_on = [google_project.project]
}

resource "google_project_service" "cloud_build" {
  service    = "cloudbuild.googleapis.com"
  depends_on = [google_project.project]
}

resource "google_project_service" "firebase" {
  service    = "firebase.googleapis.com"
  depends_on = [google_project.project]

  disable_dependent_services = true
}

resource "google_project_service" "firestore" {
  service    = "firestore.googleapis.com"
  depends_on = [google_project.project]
}

resource "google_project_service" "source_repo" {
  service    = "sourcerepo.googleapis.com"
  depends_on = [google_project.project]
}
