# Terraform Provider for Databricks on OVHcloud

This Terraform provider enables you to manage Databricks resources on OVHcloud infrastructure, providing seamless integration between Databricks cloud-native services and OVHcloud's robust infrastructure platform.

## Features

- **Workspace Management**: Create and configure Databricks workspaces on OVH infrastructure
- **Cluster Provisioning**: Auto-scaling compute clusters with spot instance support
- **Job Orchestration**: Automated workflows for ETL and ML pipelines
- **Notebook Management**: Version-controlled notebook deployment
- **Unity Catalog**: Centralized data governance and lineage tracking
- **ML Operations**: Model serving and experiment tracking integration
- **Cost Optimization**: Leverage OVH's competitive pricing and resource scheduling

## Quick Start

```hcl
terraform {
  required_providers {
    databricks-ovh = {
      source  = "spectrumwebco/databricks-ovh"
      version = "~> 0.1.0"
    }
  }
}

provider "databricks-ovh" {
  ovh_endpoint           = "ovh-eu"
  ovh_application_key    = var.ovh_application_key
  ovh_application_secret = var.ovh_application_secret
  ovh_consumer_key       = var.ovh_consumer_key
  databricks_account_id  = var.databricks_account_id
  databricks_username    = var.databricks_username
  databricks_password    = var.databricks_password
}

resource "databricks_ovh_workspace" "main" {
  name            = "production-workspace"
  region          = "eu-west-1"
  tier            = "PREMIUM"
  deployment_name = "prod-analytics"
}
```

## Resources

- `databricks_ovh_workspace` - Databricks workspace management
- `databricks_ovh_cluster` - Compute cluster provisioning
- `databricks_ovh_job` - Job and workflow automation
- `databricks_ovh_notebook` - Notebook deployment and management
- `databricks_ovh_unity_catalog` - Data governance and cataloging
- `databricks_ovh_secret_scope` - Secret management integration
- `databricks_ovh_instance_pool` - Shared compute resource pools
- `databricks_ovh_cluster_policy` - Governance and compliance policies

## Data Sources

- `databricks_ovh_workspaces` - List available workspaces
- `databricks_ovh_clusters` - Query cluster information
- `databricks_ovh_jobs` - Job discovery and monitoring

## Authentication

The provider requires OVH API credentials and Databricks authentication:

```bash
export OVH_ENDPOINT="ovh-eu"
export OVH_APPLICATION_KEY="your-app-key"
export OVH_APPLICATION_SECRET="your-app-secret"
export OVH_CONSUMER_KEY="your-consumer-key"
export DATABRICKS_ACCOUNT_ID="your-account-id"
export DATABRICKS_USERNAME="your-username"
export DATABRICKS_PASSWORD="your-password"
```

## Examples

See the `examples/` directory for complete configuration examples including:

- Basic workspace and cluster setup
- Multi-environment ML pipelines
- Unity Catalog data governance
- Advanced job orchestration workflows

## Development

```bash
# Build the provider
make build

# Run tests
make test

# Run acceptance tests
make testacc

# Install locally
make install
```

## Requirements

- Terraform >= 1.0
- Go >= 1.18
- OVH Cloud account with API access
- Databricks account or trial

## License

This project is licensed under the Mozilla Public License 2.0 - see the [LICENSE](LICENSE) file for details.

## Support

For issues and feature requests, please use the [GitHub issue tracker](https://github.com/spectrumwebco/terraform-provider-databricks-ovh/issues).
