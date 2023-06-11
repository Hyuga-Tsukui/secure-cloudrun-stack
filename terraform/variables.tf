variable "project" {}

variable "region" {
  default = "asia-northeast1"
}

variable "zone" {
  default = "asia-northeast1-b"
}

variable "gcp_service_list" {
  description = "The list of apis necessary for the project"
  type        = list(string)
  default = [
    "run.googleapis.com",
  ]
}


variable "hello-service-image-uri" {}
variable "proxy-image-uri" {}

