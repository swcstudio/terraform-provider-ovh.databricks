package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDatabricksNotebook() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a Databricks notebook",

		CreateContext: resourceDatabricksNotebookCreate,
		ReadContext:   resourceDatabricksNotebookRead,
		UpdateContext: resourceDatabricksNotebookUpdate,
		DeleteContext: resourceDatabricksNotebookDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Notebook path",
			},
			"language": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Notebook language",
				ValidateFunc: validation.StringInSlice([]string{
					"SCALA", "PYTHON", "SQL", "R",
				}, false),
			},
			"content_base64": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Base64 encoded notebook content",
			},
			"source": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Source file path",
			},
			"format": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "SOURCE",
				Description: "Notebook format",
				ValidateFunc: validation.StringInSlice([]string{
					"SOURCE", "HTML", "JUPYTER", "DBC",
				}, false),
			},
			"overwrite": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Overwrite existing notebook",
			},
			"exclude_hidden_files": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Exclude hidden files",
			},
			"object_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Object ID",
			},
			"object_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Object type",
			},
			"created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Created timestamp",
			},
			"modified_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Modified timestamp",
			},
		},
	}
}

func resourceDatabricksNotebookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	notebookConfig := map[string]interface{}{
		"path":               d.Get("path").(string),
		"language":           d.Get("language").(string),
		"contentBase64":      d.Get("content_base64").(string),
		"source":             d.Get("source").(string),
		"format":             d.Get("format").(string),
		"overwrite":          d.Get("overwrite").(bool),
		"excludeHiddenFiles": d.Get("exclude_hidden_files").(bool),
	}

	var result map[string]interface{}
	err := config.OVHClient.Post("/cloud/project/databricks/notebook", notebookConfig, &result)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Databricks notebook: %w", err))
	}

	notebookPath := result["path"].(string)
	d.SetId(notebookPath)

	return resourceDatabricksNotebookRead(ctx, d, meta)
}

func resourceDatabricksNotebookRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	notebookPath := d.Id()

	var notebook map[string]interface{}
	err := config.OVHClient.Get(fmt.Sprintf("/cloud/project/databricks/notebook?path=%s", notebookPath), &notebook)
	if err != nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("failed to read Databricks notebook: %w", err))
	}

	d.Set("path", notebook["path"])
	d.Set("language", notebook["language"])
	d.Set("format", notebook["format"])
	d.Set("object_id", notebook["objectId"])
	d.Set("object_type", notebook["objectType"])
	d.Set("created_at", notebook["createdAt"])
	d.Set("modified_at", notebook["modifiedAt"])

	return nil
}

func resourceDatabricksNotebookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	notebookPath := d.Id()

	if d.HasChanges("content_base64", "source", "format", "overwrite") {
		updateConfig := map[string]interface{}{
			"path": notebookPath,
		}

		if d.HasChange("content_base64") {
			updateConfig["contentBase64"] = d.Get("content_base64").(string)
		}
		if d.HasChange("source") {
			updateConfig["source"] = d.Get("source").(string)
		}
		if d.HasChange("format") {
			updateConfig["format"] = d.Get("format").(string)
		}
		if d.HasChange("overwrite") {
			updateConfig["overwrite"] = d.Get("overwrite").(bool)
		}

		err := config.OVHClient.Put(fmt.Sprintf("/cloud/project/databricks/notebook"), updateConfig, nil)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to update Databricks notebook: %w", err))
		}
	}

	return resourceDatabricksNotebookRead(ctx, d, meta)
}

func resourceDatabricksNotebookDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	notebookPath := d.Id()

	deleteConfig := map[string]interface{}{
		"path":      notebookPath,
		"recursive": false,
	}

	err := config.OVHClient.Post("/cloud/project/databricks/notebook/delete", deleteConfig, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Databricks notebook: %w", err))
	}

	d.SetId("")
	return nil
}
