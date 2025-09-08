package toolsearch

import (
	"errors"
	"fmt"

	"github.com/alibaba/higress/plugins/golang-filter/mcp-session/common"
	"github.com/envoyproxy/envoy/contrib/golang/common/go/api"
	"github.com/mark3labs/mcp-go/mcp"
)

const Version = "1.0.0"

func init() {
	common.GlobalRegistry.RegisterServer("tool-search", &ToolSearchConfig{})
}

type ToolSearchConfig struct {
	dsn          string
	apiKey       string
	baseURL      string
	model        string
	dimensions   int
	vectorWeight float64
	tableName    string
	description  string
}

func (c *ToolSearchConfig) ParseConfig(config map[string]any) error {
	dsn, ok := config["dsn"].(string)
	if !ok {
		return errors.New("missing dsn")
	}
	c.dsn = dsn

	apiKey, ok := config["apiKey"].(string)
	if !ok {
		return errors.New("missing apiKey")
	}
	c.apiKey = apiKey

	// Optional configurations with defaults
	if baseURL, ok := config["baseURL"].(string); ok {
		c.baseURL = baseURL
	} else {
		c.baseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	}

	if model, ok := config["model"].(string); ok {
		c.model = model
	} else {
		c.model = "text-embedding-v4"
	}

	if dimensions, ok := config["dimensions"].(float64); ok {
		c.dimensions = int(dimensions)
	} else {
		c.dimensions = 1024
	}

	if vectorWeight, ok := config["vectorWeight"].(float64); ok {
		c.vectorWeight = vectorWeight
	} else {
		c.vectorWeight = 0.5
	}

	if tableName, ok := config["tableName"].(string); ok {
		c.tableName = tableName
	} else {
		c.tableName = "apig_mcp_tools"
	}

	c.description, ok = config["description"].(string)
	if !ok {
		c.description = "Tool search server for semantic similarity search"
	}

	api.LogDebugf("ToolSearchConfig ParseConfig: %+v", config)
	return nil
}

func (c *ToolSearchConfig) NewServer(serverName string) (*common.MCPServer, error) {
	mcpServer := common.NewMCPServer(
		serverName,
		Version,
		common.WithInstructions(fmt.Sprintf("This is a tool search server: %s", c.description)),
	)

	// Create database client
	dbClient := NewDBClient(c.dsn, c.tableName, mcpServer.GetDestoryChannel())

	// Create embedding client
	embeddingClient := NewEmbeddingClient(c.apiKey, c.baseURL, c.model, c.dimensions)

	// Create search service
	textWeight := 1.0 - c.vectorWeight
	searchService := NewSearchService(dbClient, embeddingClient, c.vectorWeight, textWeight)

	// Add tool search tool
	mcpServer.AddTool(
		mcp.NewToolWithRawSchema("x_higress_tool_search", "Higress MCP Tools Searcher", GetToolSearchSchema()),
		HandleToolSearch(searchService),
	)

	return mcpServer, nil
}
