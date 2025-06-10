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

var _ resource.Resource = &DatabricksSecretScopeResource{}
var _ resource.ResourceWithImportState = &DatabricksSecretScopeResource{}

func NewDatabricksSecretScopeResource() resource.Resource {
	return &DatabricksSecretScopeResource{}
}

type DatabricksSecretScopeResource struct {
	client *Config
}

type DatabricksSecretScopeResourceModel struct {
	ID          types.String `tfsdk:"id"`
	WorkspaceID types.String `tfsdk:"workspace_id"`
	Name        types.String `tfsdk:"name"`
	ScopeID     types.String `tfsdk:"scope_id"`
	CreatedTime types.String `tfsdk:"created_time"`
}

func (r *DatabricksSecretScopeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secret_scope"
}

func (r *DatabricksSecretScopeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Databricks secret scope on OVH infrastructure",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Secret scope identifier",
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
				Description: "Secret scope name",
				Required:    true,
			},
			"scope_id": schema.StringAttribute{
				Description: "Databricks scope ID",
				Computed:    true,
			},
			"created_time": schema.StringAttribute{
				Description: "Creation timestamp",
				Computed:    true,
			},
		},
	}
}

func (r *DatabricksSecretScopeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DatabricksSecretScopeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DatabricksSecretScopeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "creating databricks secret scope resource")

	scopeConfig := map[string]interface{}{
		"workspaceId": data.WorkspaceID.ValueString(),
		"name":        data.Name.ValueString(),
	}

	var result map[string]interface{}
	err := r.client.OVHClient.Post("/cloud/project/databricks/secret-scope", scopeConfig, &result)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create secret scope, got error: %s", err))
		return
	}

	scopeId := result["id"].(string)
	data.ID = types.StringValue(scopeId)

	if databricksScopeId, ok := result["scopeId"].(string); ok {
		data.ScopeID = types.StringValue(databricksScopeId)
	}
	if createdTime, ok := result["createdTime"].(string); ok {
		data.CreatedTime = types.StringValue(createdTime)
	}

	tflog.Trace(ctx, "created databricks secret scope resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabricksSecretScopeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DatabricksSecretScopeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var scope map[string]interface{}
	err := r.client.OVHClient.Get(fmt.Sprintf("/cloud/project/databricks/secret-scope/%s", data.ID.ValueString()), &scope)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read secret scope, got error: %s", err))
		return
	}

	if workspaceId, ok := scope["workspaceId"].(string); ok {
		data.WorkspaceID = types.StringValue(workspaceId)
	}
	if name, ok := scope["name"].(string); ok {
		data.Name = types.StringValue(name)
	}
	if scopeId, ok := scope["scopeId"].(string); ok {
		data.ScopeID = types.StringValue(scopeId)
	}
	if createdTime, ok := scope["createdTime"].(string); ok {
		data.CreatedTime = types.StringValue(createdTime)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabricksSecretScopeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DatabricksSecretScopeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateConfig := map[string]interface{}{
		"name": data.Name.ValueString(),
	}

	err := r.client.OVHClient.Put(fmt.Sprintf("/cloud/project/databricks/secret-scope/%s", data.ID.ValueString()), updateConfig, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update secret scope, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabricksSecretScopeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DatabricksSecretScopeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.OVHClient.Delete(fmt.Sprintf("/cloud/project/databricks/secret-scope/%s", data.ID.ValueString()), nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete secret scope, got error: %s", err))
		return
	}
}

func (r *DatabricksSecretScopeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
