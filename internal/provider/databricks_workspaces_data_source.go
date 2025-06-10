package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &DatabricksWorkspacesDataSource{}

func NewDatabricksWorkspacesDataSource() datasource.DataSource {
	return &DatabricksWorkspacesDataSource{}
}

type DatabricksWorkspacesDataSource struct {
	client *Config
}

type DatabricksWorkspacesDataSourceModel struct {
	ID         types.String                         `tfsdk:"id"`
	Region     types.String                         `tfsdk:"region"`
	Status     types.String                         `tfsdk:"status"`
	Workspaces []DatabricksWorkspaceDataSourceModel `tfsdk:"workspaces"`
}

type DatabricksWorkspaceDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Region       types.String `tfsdk:"region"`
	Tier         types.String `tfsdk:"tier"`
	WorkspaceID  types.String `tfsdk:"workspace_id"`
	WorkspaceURL types.String `tfsdk:"workspace_url"`
	Status       types.String `tfsdk:"status"`
	CreatedTime  types.String `tfsdk:"created_time"`
}

func (d *DatabricksWorkspacesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspaces"
}

func (d *DatabricksWorkspacesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about Databricks workspaces on OVH infrastructure.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Data source identifier",
				Computed:    true,
			},
			"region": schema.StringAttribute{
				Description: "Filter workspaces by OVH region",
				Optional:    true,
			},
			"status": schema.StringAttribute{
				Description: "Filter workspaces by status",
				Optional:    true,
			},
			"workspaces": schema.ListNestedAttribute{
				Description: "List of Databricks workspaces",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Workspace identifier",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Workspace name",
							Computed:    true,
						},
						"region": schema.StringAttribute{
							Description: "OVH region",
							Computed:    true,
						},
						"tier": schema.StringAttribute{
							Description: "Databricks tier",
							Computed:    true,
						},
						"workspace_id": schema.StringAttribute{
							Description: "Databricks workspace ID",
							Computed:    true,
						},
						"workspace_url": schema.StringAttribute{
							Description: "Workspace URL",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "Workspace status",
							Computed:    true,
						},
						"created_time": schema.StringAttribute{
							Description: "Creation timestamp",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *DatabricksWorkspacesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *DatabricksWorkspacesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DatabricksWorkspacesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading Databricks workspaces")

	var workspaces []map[string]interface{}
	err := d.client.OVHClient.Get("/cloud/project/databricks/workspace", &workspaces)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workspaces, got error: %s", err))
		return
	}

	var filteredWorkspaces []DatabricksWorkspaceDataSourceModel
	for _, workspace := range workspaces {
		workspaceModel := DatabricksWorkspaceDataSourceModel{}

		if id, ok := workspace["id"].(string); ok {
			workspaceModel.ID = types.StringValue(id)
		}
		if name, ok := workspace["name"].(string); ok {
			workspaceModel.Name = types.StringValue(name)
		}
		if region, ok := workspace["region"].(string); ok {
			workspaceModel.Region = types.StringValue(region)
		}
		if tier, ok := workspace["tier"].(string); ok {
			workspaceModel.Tier = types.StringValue(tier)
		}
		if workspaceId, ok := workspace["workspaceId"].(string); ok {
			workspaceModel.WorkspaceID = types.StringValue(workspaceId)
		}
		if workspaceUrl, ok := workspace["workspaceUrl"].(string); ok {
			workspaceModel.WorkspaceURL = types.StringValue(workspaceUrl)
		}
		if status, ok := workspace["status"].(string); ok {
			workspaceModel.Status = types.StringValue(status)
		}
		if createdTime, ok := workspace["createdTime"].(string); ok {
			workspaceModel.CreatedTime = types.StringValue(createdTime)
		}

		if !data.Region.IsNull() && !data.Region.IsUnknown() {
			if workspaceModel.Region.ValueString() != data.Region.ValueString() {
				continue
			}
		}

		if !data.Status.IsNull() && !data.Status.IsUnknown() {
			if workspaceModel.Status.ValueString() != data.Status.ValueString() {
				continue
			}
		}

		filteredWorkspaces = append(filteredWorkspaces, workspaceModel)
	}

	data.Workspaces = filteredWorkspaces
	data.ID = types.StringValue("workspaces")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
