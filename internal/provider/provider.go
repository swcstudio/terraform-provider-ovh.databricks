package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
)

type Config struct {
	OVHClient *ovh.Client
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"ovh_endpoint": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("OVH_ENDPOINT", nil),
					Description: "OVH API endpoint",
				},
				"ovh_application_key": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("OVH_APPLICATION_KEY", nil),
					Description: "OVH application key",
				},
				"ovh_application_secret": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("OVH_APPLICATION_SECRET", nil),
					Description: "OVH application secret",
				},
				"ovh_consumer_key": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("OVH_CONSUMER_KEY", nil),
					Description: "OVH consumer key",
				},
				"ovh_project_id": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("OVH_PROJECT_ID", nil),
					Description: "OVH project ID",
				},
				"databricks_account_id": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("DATABRICKS_ACCOUNT_ID", nil),
					Description: "Databricks account ID",
				},
				"databricks_username": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("DATABRICKS_USERNAME", nil),
					Description: "Databricks username",
				},
				"databricks_password": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("DATABRICKS_PASSWORD", nil),
					Description: "Databricks password",
				},
				"databricks_token": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("DATABRICKS_TOKEN", nil),
					Description: "Databricks personal access token",
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"databricks_ovh_workspace":      resourceDatabricksWorkspace(),
				"databricks_ovh_job":            resourceDatabricksJob(),
				"databricks_ovh_notebook":       resourceDatabricksNotebook(),
				"databricks_ovh_secret_scope":   resourceDatabricksSecretScope(),
				"databricks_ovh_instance_pool":  resourceDatabricksInstancePool(),
				"databricks_ovh_cluster_policy": resourceDatabricksClusterPolicy(),
			},
			DataSourcesMap: map[string]*schema.Resource{
				"databricks_ovh_workspaces": dataSourceDatabricksWorkspaces(),
			},
			ConfigureContextFunc: providerConfigure,
		}

		return p
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	_ = diag.Diagnostics{}

	endpoint := d.Get("ovh_endpoint").(string)
	applicationKey := d.Get("ovh_application_key").(string)
	applicationSecret := d.Get("ovh_application_secret").(string)
	consumerKey := d.Get("ovh_consumer_key").(string)

	client, err := ovh.NewClient(
		endpoint,
		applicationKey,
		applicationSecret,
		consumerKey,
	)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	config := &Config{
		OVHClient: client,
	}

	return config, nil
}
