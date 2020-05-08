output url {
  value = google_cloud_run_service.service.status[0].url
}

output endpoint {
  value = "${trimprefix(google_cloud_run_service.service.status[0].url, "https://")}:443"
}