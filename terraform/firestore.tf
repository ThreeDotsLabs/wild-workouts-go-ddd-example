resource "null_resource" "enable_firestore" {
  triggers = {
    project = google_project.project.number
  }

  provisioner "local-exec" {
    command = "make firestore"
  }

  depends_on = [
    google_firebase_project_location.default,
  ]
}

resource "google_firestore_index" "trainings_user_time" {
  collection = "trainings"

  fields {
    field_path = "UserUuid"
    order      = "ASCENDING"
  }

  fields {
    field_path = "Time"
    order      = "ASCENDING"
  }

  fields {
    field_path = "__name__"
    order      = "ASCENDING"
  }

  depends_on = [null_resource.enable_firestore]
}
