locals {
  cloud_build_member = "serviceAccount:${google_project.project.number}@cloudbuild.gserviceaccount.com"
  compute_account    = "projects/${var.project}/serviceAccounts/${google_project.project.number}-compute@developer.gserviceaccount.com"
}

resource "google_project_iam_member" "firebase_admin" {
  role   = "roles/firebase.admin"
  member = local.cloud_build_member

  depends_on = [google_project_service.cloud_build]
}

resource "google_project_iam_member" "api_keys_admin" {
  role   = "roles/serviceusage.apiKeysViewer"
  member = local.cloud_build_member

  depends_on = [google_project_service.cloud_build]
}

resource "google_project_iam_member" "cloud_run_admin" {
  role   = "roles/run.admin"
  member = local.cloud_build_member

  depends_on = [google_project_service.cloud_build]
}

resource "google_service_account_iam_member" "default-compute-account" {
  service_account_id = local.compute_account
  role               = "roles/iam.serviceAccountUser"
  member             = local.cloud_build_member

  depends_on = [google_project_service.cloud_build]
}

