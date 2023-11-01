module "ncpaste-tf" {
  source = "./ncpaste-tf"
  gcp_project_id = var.gcp_project_id
  gcp_region = var.gcp_region
}

variable "gcp_project_id" {
  type      = string
}

variable "gcp_region" {
  type      = string
}