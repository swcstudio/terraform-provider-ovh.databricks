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

var _ resource.Resource = &DatabricksNotebookResource{}
var _ resource.ResourceWithImportState = &DatabricksNotebookResource{}

func NewDatabricksNotebookResource() resource.Resource {
	return &DatabricksNotebookResource{}
}

type DatabricksNotebookResource struct {
	client *Config
}

type DatabricksNotebookResourceModel struct {
	ID          types.String `tfsdk:"id"`
	WorkspaceID types.String `tfsdk:"workspace_id"`
	Path        types.String `tfsdk:"path"`
	Language    types.String `tfsdk:"language"`
	Content     types.String `tfsdk:"content"`
	Format      types.String `tfsdk:"format"`
	NotebookID  types.String `tfsdk:"notebook_id"`
	CreatedTime types.String `tfsdk:"created_time"`
}

func (r *DatabricksNotebookResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notebook"
}

func (r *DatabricksNotebookResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Databricks notebook on OVH infrastructure",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Notebook identifier",
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
			"path": schema.StringAttribute{
				Description: "Notebook path",
				Required:    true,
			},
			"language": schema.StringAttribute{
				Description: "Notebook language",
				Required:    true,
			},
			"content": schema.StringAttribute{
				Description: "Notebook content",
				Optional:    true,
			},
			"format": schema.StringAttribute{
				Description: "Notebook format",
				Optional:    true,
			},
			"notebook_id": schema.StringAttribute{
				Description: "Databricks notebook ID",
				Computed:    true,
			},
			"created_time": schema.StringAttribute{
				Description: "Creation timestamp",
				Computed:    true,
			},
		},
	}
}

func (r *DatabricksNotebookResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DatabricksNotebookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DatabricksNotebookResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "creating databricks notebook resource")

	notebookConfig := map[string]interface{}{
		"workspaceId": data.WorkspaceID.ValueString(),
		"path":        data.Path.ValueString(),
		"language":    data.Language.ValueString(),
		"content":     data.Content.ValueString(),
		"format":      data.Format.ValueString(),
	}

	var result map[string]interface{}
	err := r.client.OVHClient.Post("/cloud/project/databricks/notebook", notebookConfig, &result)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create notebook, got error: %s", err))
		return
	}

	notebookId := result["id"].(string)
	data.ID = types.StringValue(notebookId)

	if databricksNotebookId, ok := result["notebookId"].(string); ok {
		data.NotebookID = types.StringValue(databricksNotebookId)
	}
	if createdTime, ok := result["createdTime"].(string); ok {
		data.CreatedTime = types.StringValue(createdTime)
	}

	tflog.Trace(ctx, "created databricks notebook resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabricksNotebookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DatabricksNotebookResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var notebook map[string]interface{}
	err := r.client.OVHClient.Get(fmt.Sprintf("/cloud/project/databricks/notebook/%s", data.ID.ValueString()), &notebook)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read notebook, got error: %s", err))
		return
	}

	if workspaceId, ok := notebook["workspaceId"].(string); ok {
		data.WorkspaceID = types.StringValue(workspaceId)
	}
	if path, ok := notebook["path"].(string); ok {
		data.Path = types.StringValue(path)
	}
	if language, ok := notebook["language"].(string); ok {
		data.Language = types.StringValue(language)
	}
	if content, ok := notebook["content"].(string); ok {
		data.Content = types.StringValue(content)
	}
	if format, ok := notebook["format"].(string); ok {
		data.Format = types.StringValue(format)
	}
	if notebookId, ok := notebook["notebookId"].(string); ok {
		data.NotebookID = types.StringValue(notebookId)
	}
	if createdTime, ok := notebook["createdTime"].(string); ok {
		data.CreatedTime = types.StringValue(createdTime)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabricksNotebookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DatabricksNotebookResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateConfig := map[string]interface{}{
		"path":     data.Path.ValueString(),
		"language": data.Language.ValueString(),
		"content":  data.Content.ValueString(),
		"format":   data.Format.ValueString(),
	}

	err := r.client.OVHClient.Put(fmt.Sprintf("/cloud/project/databricks/notebook/%s", data.ID.ValueString()), updateConfig, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update notebook, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabricksNotebookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DatabricksNotebookResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.OVHClient.Delete(fmt.Sprintf("/cloud/project/databricks/notebook/%s", data.ID.ValueString()), nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete notebook, got error: %s", err))
		return
	}
}

func (r *DatabricksNotebookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
