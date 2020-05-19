resource "google_sourcerepo_repository" "wild_workouts" {
  name = var.repository_name

  depends_on = [
    google_project_service.source_repo,
  ]
}

resource "google_cloudbuild_trigger" "trigger" {
  trigger_template {
    branch_name = "master"
    repo_name   = google_sourcerepo_repository.wild_workouts.name
  }

  filename = "cloudbuild.yaml"

  depends_on = [google_sourcerepo_repository.wild_workouts]
}

resource "null_resource" "firebase_builder" {
  provisioner "local-exec" {
    command = "make firebase_builder"
  }

  depends_on = [google_project_service.container_registry]
}
