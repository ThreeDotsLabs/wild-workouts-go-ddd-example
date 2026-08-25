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

output "github_repository" {
  value = var.github_repository
}

output "github_actions_service_account" {
  value = google_service_account.github_actions.email
}

output "github_workload_identity_provider" {
  value = google_iam_workload_identity_pool_provider.github.name
}
