package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDatabricksJob() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a Databricks job",

		CreateContext: resourceDatabricksJobCreate,
		ReadContext:   resourceDatabricksJobRead,
		UpdateContext: resourceDatabricksJobUpdate,
		DeleteContext: resourceDatabricksJobDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Job name",
			},
			"new_cluster": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "New cluster configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"spark_version": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Spark version",
						},
						"node_type_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Node type ID",
						},
						"num_workers": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of workers",
						},
						"autoscale": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Autoscale configuration",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min_workers": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Minimum workers",
									},
									"max_workers": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Maximum workers",
									},
								},
							},
						},
					},
				},
			},
			"existing_cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Existing cluster ID",
			},
			"notebook_task": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Notebook task",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"notebook_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Notebook path",
						},
						"base_parameters": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Base parameters",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"spark_jar_task": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Spark JAR task",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"main_class_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Main class name",
						},
						"parameters": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Parameters",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"spark_python_task": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Spark Python task",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"python_file": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Python file path",
						},
						"parameters": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Parameters",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"spark_submit_task": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Spark submit task",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameters": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Parameters",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"pipeline_task": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Pipeline task",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pipeline_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Pipeline ID",
						},
					},
				},
			},
			"python_wheel_task": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Python wheel task",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"package_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Package name",
						},
						"entry_point": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Entry point",
						},
						"parameters": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Parameters",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
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
			"email_notifications": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Email notifications",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"on_start": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "On start emails",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"on_success": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "On success emails",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"on_failure": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "On failure emails",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"no_alert_for_skipped_runs": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "No alert for skipped runs",
						},
					},
				},
			},
			"timeout_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Timeout in seconds",
			},
			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum retries",
			},
			"min_retry_interval_millis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Minimum retry interval in milliseconds",
			},
			"retry_on_timeout": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Retry on timeout",
			},
			"schedule": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Schedule configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"quartz_cron_expression": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Quartz cron expression",
						},
						"timezone_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Timezone ID",
						},
						"pause_status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Pause status",
						},
					},
				},
			},
			"max_concurrent_runs": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Maximum concurrent runs",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Job tags",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"job_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Job ID",
			},
			"creator_user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creator username",
			},
			"created_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Created time",
			},
		},
	}
}

func resourceDatabricksJobCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	jobConfig := map[string]interface{}{
		"name":                    d.Get("name").(string),
		"newCluster":              d.Get("new_cluster").([]interface{}),
		"existingClusterId":       d.Get("existing_cluster_id").(string),
		"notebookTask":            d.Get("notebook_task").([]interface{}),
		"sparkJarTask":            d.Get("spark_jar_task").([]interface{}),
		"sparkPythonTask":         d.Get("spark_python_task").([]interface{}),
		"sparkSubmitTask":         d.Get("spark_submit_task").([]interface{}),
		"pipelineTask":            d.Get("pipeline_task").([]interface{}),
		"pythonWheelTask":         d.Get("python_wheel_task").([]interface{}),
		"libraries":               d.Get("libraries").([]interface{}),
		"emailNotifications":      d.Get("email_notifications").([]interface{}),
		"timeoutSeconds":          d.Get("timeout_seconds").(int),
		"maxRetries":              d.Get("max_retries").(int),
		"minRetryIntervalMillis":  d.Get("min_retry_interval_millis").(int),
		"retryOnTimeout":          d.Get("retry_on_timeout").(bool),
		"schedule":                d.Get("schedule").([]interface{}),
		"maxConcurrentRuns":       d.Get("max_concurrent_runs").(int),
		"tags":                    d.Get("tags"),
	}

	var result map[string]interface{}
	err := config.OVHClient.Post("/cloud/project/databricks/job", jobConfig, &result)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Databricks job: %w", err))
	}

	jobId := fmt.Sprintf("%v", result["jobId"])
	d.SetId(jobId)

	return resourceDatabricksJobRead(ctx, d, meta)
}

func resourceDatabricksJobRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	jobId := d.Id()

	var job map[string]interface{}
	err := config.OVHClient.Get(fmt.Sprintf("/cloud/project/databricks/job/%s", jobId), &job)
	if err != nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("failed to read Databricks job: %w", err))
	}

	d.Set("name", job["name"])
	d.Set("new_cluster", job["newCluster"])
	d.Set("existing_cluster_id", job["existingClusterId"])
	d.Set("notebook_task", job["notebookTask"])
	d.Set("spark_jar_task", job["sparkJarTask"])
	d.Set("spark_python_task", job["sparkPythonTask"])
	d.Set("spark_submit_task", job["sparkSubmitTask"])
	d.Set("pipeline_task", job["pipelineTask"])
	d.Set("python_wheel_task", job["pythonWheelTask"])
	d.Set("libraries", job["libraries"])
	d.Set("email_notifications", job["emailNotifications"])
	d.Set("timeout_seconds", job["timeoutSeconds"])
	d.Set("max_retries", job["maxRetries"])
	d.Set("min_retry_interval_millis", job["minRetryIntervalMillis"])
	d.Set("retry_on_timeout", job["retryOnTimeout"])
	d.Set("schedule", job["schedule"])
	d.Set("max_concurrent_runs", job["maxConcurrentRuns"])
	d.Set("job_id", job["jobId"])
	d.Set("creator_user_name", job["creatorUserName"])
	d.Set("created_time", job["createdTime"])

	if tags, ok := job["tags"].(map[string]interface{}); ok {
		d.Set("tags", tags)
	}

	return nil
}

func resourceDatabricksJobUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	jobId := d.Id()

	if d.HasChanges("name", "new_cluster", "existing_cluster_id", "notebook_task", "spark_jar_task", "spark_python_task", "spark_submit_task", "pipeline_task", "python_wheel_task", "libraries", "email_notifications", "timeout_seconds", "max_retries", "min_retry_interval_millis", "retry_on_timeout", "schedule", "max_concurrent_runs", "tags") {
		updateConfig := map[string]interface{}{}

		if d.HasChange("name") {
			updateConfig["name"] = d.Get("name").(string)
		}
		if d.HasChange("new_cluster") {
			updateConfig["newCluster"] = d.Get("new_cluster").([]interface{})
		}
		if d.HasChange("existing_cluster_id") {
			updateConfig["existingClusterId"] = d.Get("existing_cluster_id").(string)
		}
		if d.HasChange("notebook_task") {
			updateConfig["notebookTask"] = d.Get("notebook_task").([]interface{})
		}
		if d.HasChange("spark_jar_task") {
			updateConfig["sparkJarTask"] = d.Get("spark_jar_task").([]interface{})
		}
		if d.HasChange("spark_python_task") {
			updateConfig["sparkPythonTask"] = d.Get("spark_python_task").([]interface{})
		}
		if d.HasChange("spark_submit_task") {
			updateConfig["sparkSubmitTask"] = d.Get("spark_submit_task").([]interface{})
		}
		if d.HasChange("pipeline_task") {
			updateConfig["pipelineTask"] = d.Get("pipeline_task").([]interface{})
		}
		if d.HasChange("python_wheel_task") {
			updateConfig["pythonWheelTask"] = d.Get("python_wheel_task").([]interface{})
		}
		if d.HasChange("libraries") {
			updateConfig["libraries"] = d.Get("libraries").([]interface{})
		}
		if d.HasChange("email_notifications") {
			updateConfig["emailNotifications"] = d.Get("email_notifications").([]interface{})
		}
		if d.HasChange("timeout_seconds") {
			updateConfig["timeoutSeconds"] = d.Get("timeout_seconds").(int)
		}
		if d.HasChange("max_retries") {
			updateConfig["maxRetries"] = d.Get("max_retries").(int)
		}
		if d.HasChange("min_retry_interval_millis") {
			updateConfig["minRetryIntervalMillis"] = d.Get("min_retry_interval_millis").(int)
		}
		if d.HasChange("retry_on_timeout") {
			updateConfig["retryOnTimeout"] = d.Get("retry_on_timeout").(bool)
		}
		if d.HasChange("schedule") {
			updateConfig["schedule"] = d.Get("schedule").([]interface{})
		}
		if d.HasChange("max_concurrent_runs") {
			updateConfig["maxConcurrentRuns"] = d.Get("max_concurrent_runs").(int)
		}
		if d.HasChange("tags") {
			updateConfig["tags"] = d.Get("tags")
		}

		err := config.OVHClient.Put(fmt.Sprintf("/cloud/project/databricks/job/%s", jobId), updateConfig, nil)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to update Databricks job: %w", err))
		}
	}

	return resourceDatabricksJobRead(ctx, d, meta)
}

func resourceDatabricksJobDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	jobId := d.Id()

	err := config.OVHClient.Delete(fmt.Sprintf("/cloud/project/databricks/job/%s", jobId), nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Databricks job: %w", err))
	}

	d.SetId("")
	return nil
}
