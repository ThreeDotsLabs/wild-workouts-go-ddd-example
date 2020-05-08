variable project {}
variable name {}
variable location {}
variable protocol {
  description = "grpc or http"
}
variable envs {
  type = list(object({
    name  = string
    value = string
  }))
  default = []
}
variable auth {
  type    = bool
  default = true
}
variable dependency {
  type = any
}
