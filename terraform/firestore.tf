resource "google_firestore_database" "default" {
  name        = "(default)"
  location_id = var.firebase_location
  type        = "FIRESTORE_NATIVE"

  depends_on = [
    google_project_service.firestore,
    google_firebase_project.default,
  ]
}

resource "google_firestore_index" "trainings_user_time" {
  database   = google_firestore_database.default.name
  collection = "trainings"

  fields {
    field_path = "UserUuid"
    order      = "ASCENDING"
  }

  fields {
    field_path = "Canceled"
    order      = "ASCENDING"
  }

  fields {
    field_path = "Time"
    order      = "ASCENDING"
  }
}
