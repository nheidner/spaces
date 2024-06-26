variable "google_project_id" {
  type        = string
  description = "description"
}

variable "credentials_file" {
  type        = string
  description = "description"
}

variable "linode_token" {
  type = string
}

variable "kubeconfig_path" {
  type = string
}

variable "db_user_name" {
  type = string
}

variable "db_user_password" {
  type      = string
  sensitive = true
}

variable "db_name" {
  type = string
}

variable "authorized_networks" {
  type = list(map(string))
}