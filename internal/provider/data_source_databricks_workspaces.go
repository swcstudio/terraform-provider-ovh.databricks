package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDatabricksWorkspaces() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves information about Databricks workspaces on OVH infrastructure",

		ReadContext: dataSourceDatabricksWorkspacesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter workspaces by OVH region",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter workspaces by status",
			},
			"workspaces": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of Databricks workspaces",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"workspace_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Workspace ID",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Workspace name",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "OVH region",
						},
						"tier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Databricks tier",
						},
						"workspace_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Workspace URL",
						},
						"workspace_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Workspace status",
						},
						"creation_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation timestamp",
						},
						"deployment_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment name",
						},
						"aws_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AWS region",
						},
						"pricing_tier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Pricing tier",
						},
						"custom_tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Custom tags",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceDatabricksWorkspacesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	var workspaces []map[string]interface{}
	err := config.OVHClient.Get("/cloud/project/databricks/workspace", &workspaces)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Databricks workspaces: %w", err))
	}

	region := d.Get("region").(string)
	status := d.Get("status").(string)

	var filteredWorkspaces []map[string]interface{}
	for _, workspace := range workspaces {
		if region != "" && workspace["region"].(string) != region {
			continue
		}
		if status != "" && workspace["workspaceStatus"].(string) != status {
			continue
		}
		filteredWorkspaces = append(filteredWorkspaces, workspace)
	}

	workspaceList := make([]interface{}, len(filteredWorkspaces))
	for i, workspace := range filteredWorkspaces {
		workspaceMap := map[string]interface{}{
			"workspace_id":     workspace["workspaceId"],
			"name":             workspace["name"],
			"region":           workspace["region"],
			"tier":             workspace["tier"],
			"workspace_url":    workspace["workspaceUrl"],
			"workspace_status": workspace["workspaceStatus"],
			"creation_time":    workspace["creationTime"],
			"deployment_name":  workspace["deploymentName"],
			"aws_region":       workspace["awsRegion"],
			"pricing_tier":     workspace["pricingTier"],
		}

		if tags, ok := workspace["customTags"].(map[string]interface{}); ok {
			workspaceMap["custom_tags"] = tags
		}

		workspaceList[i] = workspaceMap
	}

	d.Set("workspaces", workspaceList)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}
