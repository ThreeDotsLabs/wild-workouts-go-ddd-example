module cloud_run_trainer_grpc {
  source = "./service"

  project    = var.project
  location   = var.region
  dependency = null_resource.init_docker_images

  name     = "trainer"
  protocol = "grpc"
}

module cloud_run_trainer_http {
  source = "./service"

  project    = var.project
  location   = var.region
  dependency = null_resource.init_docker_images

  name     = "trainer"
  protocol = "http"
  auth     = false

  envs = [
    {
      name  = "TRAINER_GRPC_ADDR"
      value = module.cloud_run_trainer_grpc.endpoint
    }
  ]
}

module cloud_run_trainings_http {
  source = "./service"

  project    = var.project
  location   = var.region
  dependency = null_resource.init_docker_images

  name     = "trainings"
  protocol = "http"
  auth     = false

  envs = [
    {
      name  = "TRAINER_GRPC_ADDR"
      value = module.cloud_run_trainer_grpc.endpoint
    },
    {
      name  = "USERS_GRPC_ADDR"
      value = module.cloud_run_users_grpc.endpoint
    }
  ]
}

module cloud_run_users_grpc {
  source = "./service"

  project    = var.project
  location   = var.region
  dependency = null_resource.init_docker_images

  name     = "users"
  protocol = "grpc"
}

module cloud_run_users_http {
  source = "./service"

  project    = var.project
  location   = var.region
  dependency = null_resource.init_docker_images

  name     = "users"
  protocol = "http"
  auth     = false

  envs = [
    {
      name  = "USERS_GRPC_ADDR"
      value = module.cloud_run_users_grpc.endpoint
    }
  ]
}
