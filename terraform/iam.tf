locals {
  github_actions_member = "serviceAccount:${google_service_account.github_actions.email}"
  compute_account       = "projects/${var.project}/serviceAccounts/${google_project.project.number}-compute@developer.gserviceaccount.com"
}

resource "google_project_iam_member" "firebase_admin" {
  project = google_project.project.project_id
  role    = "roles/firebase.admin"
  member  = local.github_actions_member
}

resource "google_project_iam_member" "api_keys_admin" {
  project = google_project.project.project_id
  role    = "roles/serviceusage.apiKeysViewer"
  member  = local.github_actions_member
}

resource "google_project_iam_member" "cloud_run_admin" {
  project = google_project.project.project_id
  role    = "roles/run.admin"
  member  = local.github_actions_member
}

resource "google_project_iam_member" "artifact_registry_writer" {
  project = google_project.project.project_id
  role    = "roles/artifactregistry.writer"
  member  = local.github_actions_member
}

resource "google_service_account_iam_member" "default-compute-account" {
  service_account_id = local.compute_account
  role               = "roles/iam.serviceAccountUser"
  member             = local.github_actions_member

  depends_on = [google_project_service.compute]
}
