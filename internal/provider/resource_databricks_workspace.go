package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDatabricksWorkspace() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a Databricks workspace on OVH infrastructure",

		CreateContext: resourceDatabricksWorkspaceCreate,
		ReadContext:   resourceDatabricksWorkspaceRead,
		UpdateContext: resourceDatabricksWorkspaceUpdate,
		DeleteContext: resourceDatabricksWorkspaceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Workspace name",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "OVH region",
				ValidateFunc: validation.StringInSlice([]string{
					"eu-west-1", "eu-central-1", "us-east-1", "us-west-2", "ap-southeast-1",
				}, false),
			},
			"tier": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "STANDARD",
				Description: "Databricks tier",
				ValidateFunc: validation.StringInSlice([]string{
					"STANDARD", "PREMIUM", "ENTERPRISE",
				}, false),
			},
			"deployment_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Deployment name",
			},
			"aws_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "AWS region for workspace",
			},
			"credentials_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Credentials ID",
			},
			"storage_configuration_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Storage configuration ID",
			},
			"network_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Network configuration ID",
			},
			"customer_managed_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Customer managed key ID",
			},
			"pricing_tier": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "STANDARD",
				Description: "Pricing tier",
			},
			"custom_tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Custom tags",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ovh_optimization": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable OVH infrastructure optimization",
			},
			"cost_tracking": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable cost tracking",
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Workspace ID",
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
		},
	}
}

func resourceDatabricksWorkspaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	workspaceConfig := map[string]interface{}{
		"name":                     d.Get("name").(string),
		"region":                   d.Get("region").(string),
		"tier":                     d.Get("tier").(string),
		"deploymentName":           d.Get("deployment_name").(string),
		"awsRegion":                d.Get("aws_region").(string),
		"credentialsId":            d.Get("credentials_id").(string),
		"storageConfigurationId":   d.Get("storage_configuration_id").(string),
		"networkId":                d.Get("network_id").(string),
		"customerManagedKeyId":     d.Get("customer_managed_key_id").(string),
		"pricingTier":              d.Get("pricing_tier").(string),
		"customTags":               d.Get("custom_tags"),
		"ovhOptimization":          d.Get("ovh_optimization").(bool),
		"costTracking":             d.Get("cost_tracking").(bool),
	}

	var result map[string]interface{}
	err := config.OVHClient.Post("/cloud/project/databricks/workspace", workspaceConfig, &result)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Databricks workspace: %w", err))
	}

	workspaceId := result["id"].(string)
	d.SetId(workspaceId)

	return resourceDatabricksWorkspaceRead(ctx, d, meta)
}

func resourceDatabricksWorkspaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	workspaceId := d.Id()

	var workspace map[string]interface{}
	err := config.OVHClient.Get(fmt.Sprintf("/cloud/project/databricks/workspace/%s", workspaceId), &workspace)
	if err != nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("failed to read Databricks workspace: %w", err))
	}

	d.Set("name", workspace["name"])
	d.Set("region", workspace["region"])
	d.Set("tier", workspace["tier"])
	d.Set("deployment_name", workspace["deploymentName"])
	d.Set("aws_region", workspace["awsRegion"])
	d.Set("credentials_id", workspace["credentialsId"])
	d.Set("storage_configuration_id", workspace["storageConfigurationId"])
	d.Set("network_id", workspace["networkId"])
	d.Set("customer_managed_key_id", workspace["customerManagedKeyId"])
	d.Set("pricing_tier", workspace["pricingTier"])
	d.Set("ovh_optimization", workspace["ovhOptimization"])
	d.Set("cost_tracking", workspace["costTracking"])
	d.Set("workspace_id", workspace["workspaceId"])
	d.Set("workspace_url", workspace["workspaceUrl"])
	d.Set("workspace_status", workspace["workspaceStatus"])
	d.Set("creation_time", workspace["creationTime"])

	if tags, ok := workspace["customTags"].(map[string]interface{}); ok {
		d.Set("custom_tags", tags)
	}

	return nil
}

func resourceDatabricksWorkspaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	workspaceId := d.Id()

	if d.HasChanges("name", "tier", "pricing_tier", "custom_tags") {
		updateConfig := map[string]interface{}{}

		if d.HasChange("name") {
			updateConfig["name"] = d.Get("name").(string)
		}
		if d.HasChange("tier") {
			updateConfig["tier"] = d.Get("tier").(string)
		}
		if d.HasChange("pricing_tier") {
			updateConfig["pricingTier"] = d.Get("pricing_tier").(string)
		}
		if d.HasChange("custom_tags") {
			updateConfig["customTags"] = d.Get("custom_tags")
		}

		err := config.OVHClient.Put(fmt.Sprintf("/cloud/project/databricks/workspace/%s", workspaceId), updateConfig, nil)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to update Databricks workspace: %w", err))
		}
	}

	return resourceDatabricksWorkspaceRead(ctx, d, meta)
}

func resourceDatabricksWorkspaceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	workspaceId := d.Id()

	err := config.OVHClient.Delete(fmt.Sprintf("/cloud/project/databricks/workspace/%s", workspaceId), nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Databricks workspace: %w", err))
	}

	d.SetId("")
	return nil
}
