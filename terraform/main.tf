terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.51.0"
    }
  }
}

provider "google" {
  project = var.project
  region  = var.region
  zone    = var.zone
}

module "project-services" {
  source  = "terraform-google-modules/project-factory/google//modules/project_services"
  version = "~> 14.2"

  project_id                 = var.project
  disable_dependent_services = false

  activate_apis = [
    "cloudresourcemanager.googleapis.com",
    "run.googleapis.com"
  ]
}


resource "google_cloud_run_v2_service" "hello-service" {
  name     = "hello-service"
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {

    containers {
      image = var.hello-service-image-uri
    }

    scaling {
      min_instance_count = 0
      max_instance_count = 1
    }
  }
}

resource "google_cloud_run_v2_service_iam_binding" "hello-service" {
  name     = google_cloud_run_v2_service.hello-service.name
  location = google_cloud_run_v2_service.hello-service.location
  project  = google_cloud_run_v2_service.hello-service.project
  role     = "roles/run.invoker"
  members = [
    "serviceAccount:${var.service_account_email}"
  ]
  depends_on = [google_cloud_run_v2_service.hello-service]
}


resource "google_cloud_run_v2_service" "proxy" {
  name     = "proxy"
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {

    containers {
      image = var.proxy-image-uri
      env {
        name  = "REMOTE_URL"
        value = google_cloud_run_v2_service.hello-service.uri
      }
    }

    scaling {
      min_instance_count = 0
      max_instance_count = 1
    }
  }
  depends_on = [google_cloud_run_v2_service.hello-service]
}

resource "google_cloud_run_v2_service_iam_binding" "proxy" {
  name     = google_cloud_run_v2_service.proxy.name
  location = google_cloud_run_v2_service.proxy.location
  project  = google_cloud_run_v2_service.proxy.project
  role     = "roles/run.invoker"
  members = [
    "allUsers"
  ]
  depends_on = [google_cloud_run_v2_service.proxy]
}
