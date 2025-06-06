package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDatabricksInstancePool() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a Databricks instance pool",

		CreateContext: resourceDatabricksInstancePoolCreate,
		ReadContext:   resourceDatabricksInstancePoolRead,
		UpdateContext: resourceDatabricksInstancePoolUpdate,
		DeleteContext: resourceDatabricksInstancePoolDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"instance_pool_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance pool name",
			},
			"node_type_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Node type ID",
			},
			"min_idle_instances": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				Description:  "Minimum idle instances",
				ValidateFunc: validation.IntAtLeast(0),
			},
			"max_capacity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "Maximum capacity",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"idle_instance_autotermination_minutes": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      60,
				Description:  "Idle instance auto termination in minutes",
				ValidateFunc: validation.IntAtLeast(10),
			},
			"enable_elastic_disk": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable elastic disk",
			},
			"disk_spec": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Disk specification",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Disk type",
						},
						"disk_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Disk count",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Disk size in GB",
						},
					},
				},
			},
			"preloaded_spark_versions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Preloaded Spark versions",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"aws_attributes": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "AWS attributes",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "AWS zone ID",
						},
						"instance_profile_arn": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Instance profile ARN",
						},
						"spot_bid_price_percent": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     100,
							Description: "Spot bid price percent",
						},
						"availability": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "SPOT_WITH_FALLBACK",
							Description: "Availability type",
						},
					},
				},
			},
			"custom_tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Custom tags",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"instance_pool_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance pool ID",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance pool state",
			},
			"stats": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Instance pool statistics",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"used_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Used instances count",
						},
						"idle_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Idle instances count",
						},
						"pending_used_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Pending used instances count",
						},
						"pending_idle_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Pending idle instances count",
						},
					},
				},
			},
			"default_tags": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Default tags",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceDatabricksInstancePoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	poolConfig := map[string]interface{}{
		"instancePoolName":                    d.Get("instance_pool_name").(string),
		"nodeTypeId":                          d.Get("node_type_id").(string),
		"minIdleInstances":                    d.Get("min_idle_instances").(int),
		"maxCapacity":                         d.Get("max_capacity").(int),
		"idleInstanceAutoterminationMinutes":  d.Get("idle_instance_autotermination_minutes").(int),
		"enableElasticDisk":                   d.Get("enable_elastic_disk").(bool),
		"diskSpec":                            d.Get("disk_spec").([]interface{}),
		"preloadedSparkVersions":              d.Get("preloaded_spark_versions").([]interface{}),
		"awsAttributes":                       d.Get("aws_attributes").([]interface{}),
		"customTags":                          d.Get("custom_tags"),
	}

	var result map[string]interface{}
	err := config.OVHClient.Post("/cloud/project/databricks/instance-pool", poolConfig, &result)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Databricks instance pool: %w", err))
	}

	poolId := result["instancePoolId"].(string)
	d.SetId(poolId)

	return resourceDatabricksInstancePoolRead(ctx, d, meta)
}

func resourceDatabricksInstancePoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	poolId := d.Id()

	var pool map[string]interface{}
	err := config.OVHClient.Get(fmt.Sprintf("/cloud/project/databricks/instance-pool/%s", poolId), &pool)
	if err != nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("failed to read Databricks instance pool: %w", err))
	}

	d.Set("instance_pool_name", pool["instancePoolName"])
	d.Set("node_type_id", pool["nodeTypeId"])
	d.Set("min_idle_instances", pool["minIdleInstances"])
	d.Set("max_capacity", pool["maxCapacity"])
	d.Set("idle_instance_autotermination_minutes", pool["idleInstanceAutoterminationMinutes"])
	d.Set("enable_elastic_disk", pool["enableElasticDisk"])
	d.Set("disk_spec", pool["diskSpec"])
	d.Set("preloaded_spark_versions", pool["preloadedSparkVersions"])
	d.Set("aws_attributes", pool["awsAttributes"])
	d.Set("instance_pool_id", pool["instancePoolId"])
	d.Set("state", pool["state"])
	d.Set("stats", pool["stats"])

	if customTags, ok := pool["customTags"].(map[string]interface{}); ok {
		d.Set("custom_tags", customTags)
	}
	if defaultTags, ok := pool["defaultTags"].(map[string]interface{}); ok {
		d.Set("default_tags", defaultTags)
	}

	return nil
}

func resourceDatabricksInstancePoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	poolId := d.Id()

	if d.HasChanges("instance_pool_name", "min_idle_instances", "max_capacity", "idle_instance_autotermination_minutes", "custom_tags") {
		updateConfig := map[string]interface{}{}

		if d.HasChange("instance_pool_name") {
			updateConfig["instancePoolName"] = d.Get("instance_pool_name").(string)
		}
		if d.HasChange("min_idle_instances") {
			updateConfig["minIdleInstances"] = d.Get("min_idle_instances").(int)
		}
		if d.HasChange("max_capacity") {
			updateConfig["maxCapacity"] = d.Get("max_capacity").(int)
		}
		if d.HasChange("idle_instance_autotermination_minutes") {
			updateConfig["idleInstanceAutoterminationMinutes"] = d.Get("idle_instance_autotermination_minutes").(int)
		}
		if d.HasChange("custom_tags") {
			updateConfig["customTags"] = d.Get("custom_tags")
		}

		err := config.OVHClient.Put(fmt.Sprintf("/cloud/project/databricks/instance-pool/%s", poolId), updateConfig, nil)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to update Databricks instance pool: %w", err))
		}
	}

	return resourceDatabricksInstancePoolRead(ctx, d, meta)
}

func resourceDatabricksInstancePoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	poolId := d.Id()

	err := config.OVHClient.Delete(fmt.Sprintf("/cloud/project/databricks/instance-pool/%s", poolId), nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Databricks instance pool: %w", err))
	}

	d.SetId("")
	return nil
}
