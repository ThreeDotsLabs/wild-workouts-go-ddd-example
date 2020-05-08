output "trainer_grpc_url" {
  value = module.cloud_run_trainer_grpc.url
}

output "trainer_http_url" {
  value = module.cloud_run_trainer_http.url
}

output "trainings_http_url" {
  value = module.cloud_run_trainings_http.url
}

output "users_grpc_url" {
  value = module.cloud_run_users_grpc.url
}

output "users_http_url" {
  value = module.cloud_run_users_http.url
}

output "repo_url" {
  value = google_sourcerepo_repository.wild_workouts.url
}
