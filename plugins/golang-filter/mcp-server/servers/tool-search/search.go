package toolsearch

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/envoyproxy/envoy/contrib/golang/common/go/api"
)

// SearchService handles tool search operations
type SearchService struct {
	dbClient        *DBClient
	embeddingClient *EmbeddingClient
	vectorWeight    float64
	textWeight      float64
}

// NewSearchService creates a new SearchService instance
func NewSearchService(dbClient *DBClient, embeddingClient *EmbeddingClient, vectorWeight, textWeight float64) *SearchService {
	return &SearchService{
		dbClient:        dbClient,
		embeddingClient: embeddingClient,
		vectorWeight:    vectorWeight,
		textWeight:      textWeight,
	}
}

// ToolSearchResult represents the result of a tool search
type ToolSearchResult struct {
	Tools []ToolDefinition `json:"tools"`
}

// ToolDefinition represents a tool definition in the search result
type ToolDefinition map[string]interface{}

// SearchTools performs semantic search for tools
func (s *SearchService) SearchTools(ctx context.Context, query string, topK int) (*ToolSearchResult, error) {
	api.LogDebugf("Starting tool search for query: '%s', topK: %d", query, topK)
	api.LogDebugf("Vector weight: %f, Text weight: %f", s.vectorWeight, s.textWeight)

	// Generate vector embedding for the query
	vector, err := s.embeddingClient.GetEmbedding(ctx, query)
	if err != nil {
		// Log embedding error and fallback to text-only search
		api.LogErrorf("Failed to generate embedding for query '%s': %v", query, err)
		api.LogInfof("Falling back to text-only search")

		// Perform text-only search (set vector weight to 0, text weight to 1)
		records, searchErr := s.dbClient.SearchToolsTextOnly(query, topK)
		if searchErr != nil {
			return nil, fmt.Errorf("failed to perform text-only search after embedding failure: %w", searchErr)
		}

		api.LogDebugf("Text-only search completed, found %d records", len(records))
		return s.convertRecordsToResult(records), nil
	}

	api.LogDebugf("Embedding generated successfully, vector dimension: %d", len(vector))

	// Perform hybrid search
	records, err := s.dbClient.SearchTools(query, vector, topK, s.vectorWeight, s.textWeight)
	if err != nil {
		return nil, fmt.Errorf("failed to search tools: %w", err)
	}

	api.LogDebugf("Hybrid search completed, found %d records", len(records))

	return s.convertRecordsToResult(records), nil
}

// convertRecordsToResult converts database records to tool search result
func (s *SearchService) convertRecordsToResult(records []ToolRecord) *ToolSearchResult {
	api.LogDebugf("Converting %d records to tool definitions", len(records))

	tools := make([]ToolDefinition, 0, len(records))
	for i, record := range records {
		var tool ToolDefinition

		// Parse metadata if available
		if record.Metadata != "" {
			if err := json.Unmarshal([]byte(record.Metadata), &tool); err != nil {
				api.LogWarnf("Failed to parse metadata for tool %s___%s: %v", record.ServerName, record.Name, err)
				// If metadata parsing fails, create a basic tool definition
				tool = ToolDefinition{
					"name":        record.Name,
					"description": record.Description,
				}
			} else {
				api.LogDebugf("Successfully parsed metadata for tool %s___%s", record.ServerName, record.Name)
			}
		} else {
			api.LogDebugf("No metadata found for tool %s___%s, using basic definition", record.ServerName, record.Name)
			// If no metadata, create a basic tool definition
			tool = ToolDefinition{
				"name":        record.Name,
				"description": record.Description,
			}
		}

		// Update the name to include server name
		tool["name"] = fmt.Sprintf("%s___%s", record.ServerName, record.Name)

		tools = append(tools, tool)

		if i < 3 { // Log first 3 tools for debugging
			api.LogDebugf("Tool %d: %s - %s", i+1, tool["name"], record.Description)
		}
	}

	api.LogDebugf("Successfully converted %d tools", len(tools))
	return &ToolSearchResult{Tools: tools}
}

// GetAllTools retrieves all available tools
func (s *SearchService) GetAllTools() (*ToolSearchResult, error) {
	records, err := s.dbClient.GetAllTools()
	if err != nil {
		return nil, fmt.Errorf("failed to get all tools: %w", err)
	}

	// Convert records to tool definitions
	tools := make([]ToolDefinition, 0, len(records))
	for _, record := range records {
		var tool ToolDefinition

		// Parse metadata if available
		if record.Metadata != "" {
			if err := json.Unmarshal([]byte(record.Metadata), &tool); err != nil {
				// If metadata parsing fails, create a basic tool definition
				tool = ToolDefinition{
					"name":        record.Name,
					"description": record.Description,
				}
			}
		} else {
			// If no metadata, create a basic tool definition
			tool = ToolDefinition{
				"name":        record.Name,
				"description": record.Description,
			}
		}

		// Update the name to include server name
		tool["name"] = fmt.Sprintf("%s___%s", record.ServerName, record.Name)

		tools = append(tools, tool)
	}

	return &ToolSearchResult{Tools: tools}, nil
}
