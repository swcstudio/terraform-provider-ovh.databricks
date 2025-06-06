package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDatabricksClusterPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a Databricks cluster policy",

		CreateContext: resourceDatabricksClusterPolicyCreate,
		ReadContext:   resourceDatabricksClusterPolicyRead,
		UpdateContext: resourceDatabricksClusterPolicyUpdate,
		DeleteContext: resourceDatabricksClusterPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster policy name",
			},
			"definition": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Policy definition in JSON format",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Policy description",
			},
			"max_clusters_per_user": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum clusters per user",
			},
			"policy_family_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Policy family ID",
			},
			"policy_family_definition_overrides": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Policy family definition overrides in JSON format",
			},
			"libraries": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Libraries",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"jar": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "JAR library",
						},
						"egg": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Egg library",
						},
						"whl": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Wheel library",
						},
						"pypi": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "PyPI library",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"package": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Package name",
									},
									"repo": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Repository",
									},
								},
							},
						},
						"maven": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Maven library",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"coordinates": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Maven coordinates",
									},
									"repo": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Repository",
									},
									"exclusions": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Exclusions",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"cran": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "CRAN library",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"package": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Package name",
									},
									"repo": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Repository",
									},
								},
							},
						},
					},
				},
			},
			"policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy ID",
			},
			"created_at_timestamp": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Created timestamp",
			},
		},
	}
}

func resourceDatabricksClusterPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	policyConfig := map[string]interface{}{
		"name":                             d.Get("name").(string),
		"definition":                       d.Get("definition").(string),
		"description":                      d.Get("description").(string),
		"maxClustersPerUser":               d.Get("max_clusters_per_user").(int),
		"policyFamilyId":                   d.Get("policy_family_id").(string),
		"policyFamilyDefinitionOverrides":  d.Get("policy_family_definition_overrides").(string),
		"libraries":                        d.Get("libraries").([]interface{}),
	}

	var result map[string]interface{}
	err := config.OVHClient.Post("/cloud/project/databricks/cluster-policy", policyConfig, &result)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Databricks cluster policy: %w", err))
	}

	policyId := result["policyId"].(string)
	d.SetId(policyId)

	return resourceDatabricksClusterPolicyRead(ctx, d, meta)
}

func resourceDatabricksClusterPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	policyId := d.Id()

	var policy map[string]interface{}
	err := config.OVHClient.Get(fmt.Sprintf("/cloud/project/databricks/cluster-policy/%s", policyId), &policy)
	if err != nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("failed to read Databricks cluster policy: %w", err))
	}

	d.Set("name", policy["name"])
	d.Set("definition", policy["definition"])
	d.Set("description", policy["description"])
	d.Set("max_clusters_per_user", policy["maxClustersPerUser"])
	d.Set("policy_family_id", policy["policyFamilyId"])
	d.Set("policy_family_definition_overrides", policy["policyFamilyDefinitionOverrides"])
	d.Set("libraries", policy["libraries"])
	d.Set("policy_id", policy["policyId"])
	d.Set("created_at_timestamp", policy["createdAtTimestamp"])

	return nil
}

func resourceDatabricksClusterPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	policyId := d.Id()

	if d.HasChanges("name", "definition", "description", "max_clusters_per_user", "policy_family_definition_overrides", "libraries") {
		updateConfig := map[string]interface{}{}

		if d.HasChange("name") {
			updateConfig["name"] = d.Get("name").(string)
		}
		if d.HasChange("definition") {
			updateConfig["definition"] = d.Get("definition").(string)
		}
		if d.HasChange("description") {
			updateConfig["description"] = d.Get("description").(string)
		}
		if d.HasChange("max_clusters_per_user") {
			updateConfig["maxClustersPerUser"] = d.Get("max_clusters_per_user").(int)
		}
		if d.HasChange("policy_family_definition_overrides") {
			updateConfig["policyFamilyDefinitionOverrides"] = d.Get("policy_family_definition_overrides").(string)
		}
		if d.HasChange("libraries") {
			updateConfig["libraries"] = d.Get("libraries").([]interface{})
		}

		err := config.OVHClient.Put(fmt.Sprintf("/cloud/project/databricks/cluster-policy/%s", policyId), updateConfig, nil)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to update Databricks cluster policy: %w", err))
		}
	}

	return resourceDatabricksClusterPolicyRead(ctx, d, meta)
}

func resourceDatabricksClusterPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	policyId := d.Id()

	err := config.OVHClient.Delete(fmt.Sprintf("/cloud/project/databricks/cluster-policy/%s", policyId), nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Databricks cluster policy: %w", err))
	}

	d.SetId("")
	return nil
}
