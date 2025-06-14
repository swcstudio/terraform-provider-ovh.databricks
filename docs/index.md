---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "databricks-ovh Provider"
subcategory: ""
description: |-
  The Databricks OVH provider enables management of Databricks resources on OVH cloud infrastructure.
---

# databricks-ovh Provider

The Databricks OVH provider enables management of Databricks resources on OVH cloud infrastructure.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `ovh_application_key` (String) OVH application key
- `ovh_application_secret` (String, Sensitive) OVH application secret
- `ovh_consumer_key` (String, Sensitive) OVH consumer key
- `ovh_endpoint` (String) OVH API endpoint

### Optional

- `databricks_account_id` (String) Databricks account ID
- `databricks_password` (String, Sensitive) Databricks password
- `databricks_token` (String, Sensitive) Databricks personal access token
- `databricks_username` (String) Databricks username
- `ovh_project_id` (String) OVH project ID
