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

var _ resource.Resource = &DatabricksInstancePoolResource{}
var _ resource.ResourceWithImportState = &DatabricksInstancePoolResource{}

func NewDatabricksInstancePoolResource() resource.Resource {
	return &DatabricksInstancePoolResource{}
}

type DatabricksInstancePoolResource struct {
	client *Config
}

type DatabricksInstancePoolResourceModel struct {
	ID               types.String `tfsdk:"id"`
	WorkspaceID      types.String `tfsdk:"workspace_id"`
	Name             types.String `tfsdk:"name"`
	NodeTypeID       types.String `tfsdk:"node_type_id"`
	MinIdleInstances types.Int64  `tfsdk:"min_idle_instances"`
	MaxCapacity      types.Int64  `tfsdk:"max_capacity"`
	PoolID           types.String `tfsdk:"pool_id"`
	Status           types.String `tfsdk:"status"`
	CreatedTime      types.String `tfsdk:"created_time"`
}

func (r *DatabricksInstancePoolResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance_pool"
}

func (r *DatabricksInstancePoolResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Databricks instance pool on OVH infrastructure",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Instance pool identifier",
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
				Description: "Instance pool name",
				Required:    true,
			},
			"node_type_id": schema.StringAttribute{
				Description: "Node type ID",
				Required:    true,
			},
			"min_idle_instances": schema.Int64Attribute{
				Description: "Minimum idle instances",
				Optional:    true,
			},
			"max_capacity": schema.Int64Attribute{
				Description: "Maximum capacity",
				Optional:    true,
			},
			"pool_id": schema.StringAttribute{
				Description: "Databricks pool ID",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Pool status",
				Computed:    true,
			},
			"created_time": schema.StringAttribute{
				Description: "Creation timestamp",
				Computed:    true,
			},
		},
	}
}

func (r *DatabricksInstancePoolResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DatabricksInstancePoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DatabricksInstancePoolResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "creating databricks instance pool resource")

	poolConfig := map[string]interface{}{
		"workspaceId":      data.WorkspaceID.ValueString(),
		"name":             data.Name.ValueString(),
		"nodeTypeId":       data.NodeTypeID.ValueString(),
		"minIdleInstances": data.MinIdleInstances.ValueInt64(),
		"maxCapacity":      data.MaxCapacity.ValueInt64(),
	}

	var result map[string]interface{}
	err := r.client.OVHClient.Post("/cloud/project/databricks/instance-pool", poolConfig, &result)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create instance pool, got error: %s", err))
		return
	}

	poolId := result["id"].(string)
	data.ID = types.StringValue(poolId)

	if databricksPoolId, ok := result["poolId"].(string); ok {
		data.PoolID = types.StringValue(databricksPoolId)
	}
	if status, ok := result["status"].(string); ok {
		data.Status = types.StringValue(status)
	}
	if createdTime, ok := result["createdTime"].(string); ok {
		data.CreatedTime = types.StringValue(createdTime)
	}

	tflog.Trace(ctx, "created databricks instance pool resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabricksInstancePoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DatabricksInstancePoolResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var pool map[string]interface{}
	err := r.client.OVHClient.Get(fmt.Sprintf("/cloud/project/databricks/instance-pool/%s", data.ID.ValueString()), &pool)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read instance pool, got error: %s", err))
		return
	}

	if workspaceId, ok := pool["workspaceId"].(string); ok {
		data.WorkspaceID = types.StringValue(workspaceId)
	}
	if name, ok := pool["name"].(string); ok {
		data.Name = types.StringValue(name)
	}
	if nodeTypeId, ok := pool["nodeTypeId"].(string); ok {
		data.NodeTypeID = types.StringValue(nodeTypeId)
	}
	if minIdleInstances, ok := pool["minIdleInstances"].(float64); ok {
		data.MinIdleInstances = types.Int64Value(int64(minIdleInstances))
	}
	if maxCapacity, ok := pool["maxCapacity"].(float64); ok {
		data.MaxCapacity = types.Int64Value(int64(maxCapacity))
	}
	if poolId, ok := pool["poolId"].(string); ok {
		data.PoolID = types.StringValue(poolId)
	}
	if status, ok := pool["status"].(string); ok {
		data.Status = types.StringValue(status)
	}
	if createdTime, ok := pool["createdTime"].(string); ok {
		data.CreatedTime = types.StringValue(createdTime)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabricksInstancePoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DatabricksInstancePoolResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateConfig := map[string]interface{}{
		"name":             data.Name.ValueString(),
		"nodeTypeId":       data.NodeTypeID.ValueString(),
		"minIdleInstances": data.MinIdleInstances.ValueInt64(),
		"maxCapacity":      data.MaxCapacity.ValueInt64(),
	}

	err := r.client.OVHClient.Put(fmt.Sprintf("/cloud/project/databricks/instance-pool/%s", data.ID.ValueString()), updateConfig, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update instance pool, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabricksInstancePoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DatabricksInstancePoolResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.OVHClient.Delete(fmt.Sprintf("/cloud/project/databricks/instance-pool/%s", data.ID.ValueString()), nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete instance pool, got error: %s", err))
		return
	}
}

func (r *DatabricksInstancePoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
