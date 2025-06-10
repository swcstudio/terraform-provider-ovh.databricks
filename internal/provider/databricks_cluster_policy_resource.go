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

var _ resource.Resource = &DatabricksClusterPolicyResource{}
var _ resource.ResourceWithImportState = &DatabricksClusterPolicyResource{}

func NewDatabricksClusterPolicyResource() resource.Resource {
	return &DatabricksClusterPolicyResource{}
}

type DatabricksClusterPolicyResource struct {
	client *Config
}

type DatabricksClusterPolicyResourceModel struct {
	ID          types.String `tfsdk:"id"`
	WorkspaceID types.String `tfsdk:"workspace_id"`
	Name        types.String `tfsdk:"name"`
	Definition  types.String `tfsdk:"definition"`
	PolicyID    types.String `tfsdk:"policy_id"`
	CreatedTime types.String `tfsdk:"created_time"`
}

func (r *DatabricksClusterPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_policy"
}

func (r *DatabricksClusterPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Databricks cluster policy on OVH infrastructure",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Policy identifier",
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
				Description: "Policy name",
				Required:    true,
			},
			"definition": schema.StringAttribute{
				Description: "Policy definition JSON",
				Required:    true,
			},
			"policy_id": schema.StringAttribute{
				Description: "Databricks policy ID",
				Computed:    true,
			},
			"created_time": schema.StringAttribute{
				Description: "Creation timestamp",
				Computed:    true,
			},
		},
	}
}

func (r *DatabricksClusterPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DatabricksClusterPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DatabricksClusterPolicyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "creating databricks cluster policy resource")

	policyConfig := map[string]interface{}{
		"workspaceId": data.WorkspaceID.ValueString(),
		"name":        data.Name.ValueString(),
		"definition":  data.Definition.ValueString(),
	}

	var result map[string]interface{}
	err := r.client.OVHClient.Post("/cloud/project/databricks/cluster-policy", policyConfig, &result)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create cluster policy, got error: %s", err))
		return
	}

	policyId := result["id"].(string)
	data.ID = types.StringValue(policyId)

	if databricksPolicyId, ok := result["policyId"].(string); ok {
		data.PolicyID = types.StringValue(databricksPolicyId)
	}
	if createdTime, ok := result["createdTime"].(string); ok {
		data.CreatedTime = types.StringValue(createdTime)
	}

	tflog.Trace(ctx, "created databricks cluster policy resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabricksClusterPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DatabricksClusterPolicyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var policy map[string]interface{}
	err := r.client.OVHClient.Get(fmt.Sprintf("/cloud/project/databricks/cluster-policy/%s", data.ID.ValueString()), &policy)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cluster policy, got error: %s", err))
		return
	}

	if workspaceId, ok := policy["workspaceId"].(string); ok {
		data.WorkspaceID = types.StringValue(workspaceId)
	}
	if name, ok := policy["name"].(string); ok {
		data.Name = types.StringValue(name)
	}
	if definition, ok := policy["definition"].(string); ok {
		data.Definition = types.StringValue(definition)
	}
	if policyId, ok := policy["policyId"].(string); ok {
		data.PolicyID = types.StringValue(policyId)
	}
	if createdTime, ok := policy["createdTime"].(string); ok {
		data.CreatedTime = types.StringValue(createdTime)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabricksClusterPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DatabricksClusterPolicyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateConfig := map[string]interface{}{
		"name":       data.Name.ValueString(),
		"definition": data.Definition.ValueString(),
	}

	err := r.client.OVHClient.Put(fmt.Sprintf("/cloud/project/databricks/cluster-policy/%s", data.ID.ValueString()), updateConfig, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update cluster policy, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabricksClusterPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DatabricksClusterPolicyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.OVHClient.Delete(fmt.Sprintf("/cloud/project/databricks/cluster-policy/%s", data.ID.ValueString()), nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete cluster policy, got error: %s", err))
		return
	}
}

func (r *DatabricksClusterPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
