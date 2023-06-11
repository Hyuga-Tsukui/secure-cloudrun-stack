variable "project" {}

variable "region" {
  default = "asia-northeast1"
}

variable "zone" {
  default = "asia-northeast1-b"
}

variable "hello-service-image-uri" {}
variable "proxy-image-uri" {}
variable "service_account_email" {
  type    = string
  default = "76919647250-compute@developer.gserviceaccount.com"
}
