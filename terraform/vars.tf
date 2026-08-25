variable "project" {}
variable "user" {}
variable "region" {}
variable "firebase_location" {}

variable "billing_account" {
  description = "Billing account display name"
}

variable "github_repository" {
  description = "GitHub repository (owner/name) allowed to deploy via GitHub Actions"
  default     = "ThreeDotsLabs/wild-workouts-go-ddd-example"
}
