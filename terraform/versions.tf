terraform {
  required_version = ">= 1.5"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 7.45"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 7.45"
    }
    null = {
      source  = "hashicorp/null"
      version = "~> 3.3"
    }
  }
}
