resource "null_resource" "init_docker_images" {
  triggers = {
    project = google_project.project.number
  }

  provisioner "local-exec" {
    command = "make docker_images"
  }

  depends_on = [google_project_service.container_registry]
}
