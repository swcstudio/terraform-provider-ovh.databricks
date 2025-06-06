terraform {
  required_providers {
    databricks-ovh = {
      source  = "spectrumwebco/databricks-ovh"
      version = "~> 0.1.0"
    }
  }
}

provider "databricks-ovh" {
  ovh_endpoint           = var.ovh_endpoint
  ovh_application_key    = var.ovh_application_key
  ovh_application_secret = var.ovh_application_secret
  ovh_consumer_key       = var.ovh_consumer_key
  databricks_account_id  = var.databricks_account_id
  databricks_username    = var.databricks_username
  databricks_password    = var.databricks_password
}

resource "databricks_ovh_workspace" "example" {
  name                = "example-workspace"
  region              = "eu-west-1"
  tier                = "PREMIUM"
  deployment_name     = "example-deployment"
  aws_region          = "eu-west-1"
  pricing_tier        = "PREMIUM"
  
  custom_tags = {
    Environment = "development"
    Team        = "ml"
  }
}

resource "databricks_ovh_cluster" "example" {
  workspace_id     = databricks_ovh_workspace.example.workspace_id
  cluster_name     = "example-cluster"
  spark_version    = "11.3.x-scala2.12"
  node_type_id     = "i3.xlarge"
  driver_node_type_id = "i3.xlarge"
  
  num_workers = 2
  
  autoscale {
    min_workers = 1
    max_workers = 4
  }
  
  autotermination_minutes = 60
  
  spark_conf = {
    "spark.databricks.cluster.profile" = "singleNode"
    "spark.master"                      = "local[*]"
  }
  
  custom_tags = {
    Environment = "development"
    Purpose     = "analytics"
  }
}

resource "databricks_ovh_job" "example" {
  workspace_id = databricks_ovh_workspace.example.workspace_id
  name         = "example-job"
  
  new_cluster {
    spark_version   = "11.3.x-scala2.12"
    node_type_id    = "i3.xlarge"
    num_workers     = 1
  }
  
  notebook_task {
    notebook_path = "/example-notebook"
  }
  
  email_notifications {
    on_start   = ["admin@example.com"]
    on_success = ["admin@example.com"]
    on_failure = ["admin@example.com"]
  }
  
  timeout_seconds = 3600
  max_retries     = 2
}

output "workspace_url" {
  value = databricks_ovh_workspace.example.workspace_url
}

output "cluster_id" {
  value = databricks_ovh_cluster.example.cluster_id
}
