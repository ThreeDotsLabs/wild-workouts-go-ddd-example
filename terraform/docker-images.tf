resource "google_artifact_registry_repository" "docker" {
  location      = var.region
  repository_id = "docker"
  format        = "DOCKER"

  depends_on = [google_project_service.artifact_registry]
}

resource "null_resource" "init_docker_images" {
  provisioner "local-exec" {
    command = "task docker-images"
  }

  depends_on = [google_artifact_registry_repository.docker]
}
