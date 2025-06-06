package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDatabricksSecretScope() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a Databricks secret scope",

		CreateContext: resourceDatabricksSecretScopeCreate,
		ReadContext:   resourceDatabricksSecretScopeRead,
		UpdateContext: resourceDatabricksSecretScopeUpdate,
		DeleteContext: resourceDatabricksSecretScopeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Secret scope name",
			},
			"initial_manage_principal": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "users",
				Description: "Initial manage principal",
				ValidateFunc: validation.StringInSlice([]string{
					"users", "admins",
				}, false),
			},
			"keyvault_metadata": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Azure Key Vault metadata",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Azure Key Vault resource ID",
						},
						"dns_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Azure Key Vault DNS name",
						},
					},
				},
			},
			"backend_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Backend type",
			},
		},
	}
}

func resourceDatabricksSecretScopeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	scopeConfig := map[string]interface{}{
		"scope":                  d.Get("name").(string),
		"initialManagePrincipal": d.Get("initial_manage_principal").(string),
		"keyvaultMetadata":       d.Get("keyvault_metadata").([]interface{}),
	}

	var result map[string]interface{}
	err := config.OVHClient.Post("/cloud/project/databricks/secret-scope", scopeConfig, &result)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Databricks secret scope: %w", err))
	}

	scopeName := result["scope"].(string)
	d.SetId(scopeName)

	return resourceDatabricksSecretScopeRead(ctx, d, meta)
}

func resourceDatabricksSecretScopeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	scopeName := d.Id()

	var scope map[string]interface{}
	err := config.OVHClient.Get(fmt.Sprintf("/cloud/project/databricks/secret-scope/%s", scopeName), &scope)
	if err != nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("failed to read Databricks secret scope: %w", err))
	}

	d.Set("name", scope["name"])
	d.Set("backend_type", scope["backendType"])
	d.Set("keyvault_metadata", scope["keyvaultMetadata"])

	return nil
}

func resourceDatabricksSecretScopeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceDatabricksSecretScopeRead(ctx, d, meta)
}

func resourceDatabricksSecretScopeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	scopeName := d.Id()

	err := config.OVHClient.Delete(fmt.Sprintf("/cloud/project/databricks/secret-scope/%s", scopeName), nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Databricks secret scope: %w", err))
	}

	d.SetId("")
	return nil
}
