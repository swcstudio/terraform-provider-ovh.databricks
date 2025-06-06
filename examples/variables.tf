variable "ovh_endpoint" {
  description = "OVH API endpoint"
  type        = string
  default     = "ovh-eu"
}

variable "ovh_application_key" {
  description = "OVH application key"
  type        = string
  sensitive   = true
}

variable "ovh_application_secret" {
  description = "OVH application secret"
  type        = string
  sensitive   = true
}

variable "ovh_consumer_key" {
  description = "OVH consumer key"
  type        = string
  sensitive   = true
}

variable "databricks_account_id" {
  description = "Databricks account ID"
  type        = string
}

variable "databricks_username" {
  description = "Databricks username"
  type        = string
}

variable "databricks_password" {
  description = "Databricks password"
  type        = string
  sensitive   = true
}
