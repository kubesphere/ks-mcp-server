package userrole

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	iamv1beta1 "kubesphere.io/api/iam/v1beta1"

	"kubesphere.io/ks-mcp-server/pkg/constants"
	"kubesphere.io/ks-mcp-server/pkg/kubesphere"
)

func ListUsers(ksconfig *kubesphere.KSConfig) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("list_users", mcp.WithDescription(`
Retrieve the paginated user list. The response will contain:
1. items: An array containing user data where:
  - username: Maps to metadata.name
  - specific metadata.annotations fields indicate:
    - iam.kubesphere.io/globalrole: The user's assigned cluster role
    - iam.kubesphere.io/granted-clusters: Clusters assigned to the user
2. totalItems: The total number of users in KubeSphere.
`),
			mcp.WithNumber("limit", mcp.Description("Number of users displayed at once. Default is "+constants.DefLimit)),
			mcp.WithNumber("page", mcp.Description("Page number of users to display. Default is "+constants.DefPage)),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// deal request params
			limit := constants.DefLimit
			if reqLimit, ok := request.Params.Arguments["limit"].(int64); ok && reqLimit != 0 {
				limit = fmt.Sprintf("%d", reqLimit)
			}
			page := constants.DefPage
			if reqPage, ok := request.Params.Arguments["page"].(int64); ok && reqPage != 0 {
				page = fmt.Sprintf("%d", reqPage)
			}
			// deal http request
			client, err := ksconfig.RestClient(iamv1beta1.SchemeGroupVersion, "")
			if err != nil {
				return nil, err
			}
			data, err := client.Get().Resource(iamv1beta1.ResourcesPluralUser).
				Param("sortBy", "createTime").Param("limit", limit).Param("page", page).Do(ctx).Raw()
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(string(data)), nil
		},
	}
}

func GetUser(ksconfig *kubesphere.KSConfig) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("get_user", mcp.WithDescription(`
Get user information by username. The response will contain:
- username: Maps to metadata.name
- specific metadata.annotations fields indicate:
  - iam.kubesphere.io/globalrole: The user's assigned platform role
  - iam.kubesphere.io/granted-clusters: Clusters assigned to the user
`),
			mcp.WithString("user", mcp.Description("the given username"), mcp.Required()),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// deal http request
			client, err := ksconfig.RestClient(iamv1beta1.SchemeGroupVersion, "")
			if err != nil {
				return nil, err
			}
			data, err := client.Get().Resource(iamv1beta1.ResourcesPluralUser).SubResource(request.Params.Arguments["user"].(string)).Do(ctx).Raw()
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(string(data)), nil
		},
	}
}
