package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &DatabricksJobResource{}
var _ resource.ResourceWithImportState = &DatabricksJobResource{}

func NewDatabricksJobResource() resource.Resource {
	return &DatabricksJobResource{}
}

type DatabricksJobResource struct {
	client *Config
}

type DatabricksJobResourceModel struct {
	ID          types.String `tfsdk:"id"`
	WorkspaceID types.String `tfsdk:"workspace_id"`
	Name        types.String `tfsdk:"name"`
	JobID       types.String `tfsdk:"job_id"`
	Status      types.String `tfsdk:"status"`
	CreatedTime types.String `tfsdk:"created_time"`
}

func (r *DatabricksJobResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_job"
}

func (r *DatabricksJobResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Databricks job on OVH infrastructure",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Job identifier",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workspace_id": schema.StringAttribute{
				Description: "Workspace ID",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Job name",
				Required:    true,
			},
			"job_id": schema.StringAttribute{
				Description: "Databricks job ID",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Job status",
				Computed:    true,
			},
			"created_time": schema.StringAttribute{
				Description: "Creation timestamp",
				Computed:    true,
			},
		},
	}
}

func (r *DatabricksJobResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *DatabricksJobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DatabricksJobResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "creating databricks job resource")

	jobConfig := map[string]interface{}{
		"workspaceId": data.WorkspaceID.ValueString(),
		"name":        data.Name.ValueString(),
	}

	var result map[string]interface{}
	err := r.client.OVHClient.Post("/cloud/project/databricks/job", jobConfig, &result)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create job, got error: %s", err))
		return
	}

	jobId := result["id"].(string)
	data.ID = types.StringValue(jobId)

	if databricksJobId, ok := result["jobId"].(string); ok {
		data.JobID = types.StringValue(databricksJobId)
	}
	if status, ok := result["status"].(string); ok {
		data.Status = types.StringValue(status)
	}
	if createdTime, ok := result["createdTime"].(string); ok {
		data.CreatedTime = types.StringValue(createdTime)
	}

	tflog.Trace(ctx, "created databricks job resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabricksJobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DatabricksJobResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var job map[string]interface{}
	err := r.client.OVHClient.Get(fmt.Sprintf("/cloud/project/databricks/job/%s", data.ID.ValueString()), &job)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read job, got error: %s", err))
		return
	}

	if workspaceId, ok := job["workspaceId"].(string); ok {
		data.WorkspaceID = types.StringValue(workspaceId)
	}
	if name, ok := job["name"].(string); ok {
		data.Name = types.StringValue(name)
	}
	if jobId, ok := job["jobId"].(string); ok {
		data.JobID = types.StringValue(jobId)
	}
	if status, ok := job["status"].(string); ok {
		data.Status = types.StringValue(status)
	}
	if createdTime, ok := job["createdTime"].(string); ok {
		data.CreatedTime = types.StringValue(createdTime)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabricksJobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DatabricksJobResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateConfig := map[string]interface{}{
		"name": data.Name.ValueString(),
	}

	err := r.client.OVHClient.Put(fmt.Sprintf("/cloud/project/databricks/job/%s", data.ID.ValueString()), updateConfig, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update job, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabricksJobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DatabricksJobResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.OVHClient.Delete(fmt.Sprintf("/cloud/project/databricks/job/%s", data.ID.ValueString()), nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete job, got error: %s", err))
		return
	}
}

func (r *DatabricksJobResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
