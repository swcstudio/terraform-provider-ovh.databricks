---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "databricks-ovh_notebook Resource - terraform-provider-databricks-ovh"
subcategory: ""
description: |-
  Manages a Databricks notebook on OVH infrastructure
---

# databricks-ovh_notebook (Resource)

Manages a Databricks notebook on OVH infrastructure



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `language` (String) Notebook language
- `path` (String) Notebook path
- `workspace_id` (String) Workspace ID

### Optional

- `content` (String) Notebook content
- `format` (String) Notebook format

### Read-Only

- `created_time` (String) Creation timestamp
- `id` (String) Notebook identifier
- `notebook_id` (String) Databricks notebook ID
