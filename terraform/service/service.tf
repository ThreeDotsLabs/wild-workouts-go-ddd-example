locals {
  default_envs = [
    {
      name  = "GCP_PROJECT"
      value = var.project
    },
    {
      name  = "SERVER_TO_RUN"
      value = var.protocol
    }
  ]
}

data "google_container_registry_image" "image" {
  name = var.name
}

resource "google_cloud_run_service" "service" {
  name     = "${var.name}-${var.protocol}"
  location = var.location

  template {
    spec {
      containers {
        image = data.google_container_registry_image.image.image_url

        dynamic "env" {
          for_each = concat(local.default_envs, var.envs)
          content {
            name  = env.value.name
            value = env.value.value
          }
        }
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "3"
      }
    }
  }

  autogenerate_revision_name = true

  depends_on = [var.dependency]
}

data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth_policy" {
  count = var.auth ? 0 : 1

  location = google_cloud_run_service.service.location
  service  = google_cloud_run_service.service.name

  policy_data = data.google_iam_policy.noauth.policy_data
}
