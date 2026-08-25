variable "project" {}
variable "name" {}
variable "location" {}
variable "repository" {
  description = "Artifact Registry repository holding the service images"
}
variable "protocol" {
  description = "grpc or http"
}
variable "envs" {
  type = list(object({
    name  = string
    value = string
  }))
  default = []
}
variable "auth" {
  type    = bool
  default = true
}
variable "dependency" {
  type = any
}
