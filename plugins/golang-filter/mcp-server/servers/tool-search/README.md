# Tool Search MCP Server

This is an MCP Server based on Higress golang-filter for implementing tool search functionality. It supports hybrid search combining vector similarity and full-text search for semantic retrieval of the most relevant topK tools based on user queries.

## Features

- **Hybrid Search**: Combines vector similarity search and full-text search
- **Semantic Search**: Uses OpenAI-compatible embedding APIs to generate text vectors for semantic similarity comparison
- **Configurable Weights**: Supports configurable weight ratios for vector search and full-text search
- **PostgreSQL Integration**: Utilizes PostgreSQL vector extensions and full-text search capabilities

## Database Schema

Create the following table in PostgreSQL:

```sql
-- Install required extensions
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS pgsearch;

-- Note: Table name is configurable via the tableName parameter
CREATE TABLE apig_mcp_tools (
    id SERIAL PRIMARY KEY,
    server_name VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    metadata JSONB,
    vector VECTOR(1024)  -- Requires pgvector extension
);

-- Create full-text search index (using pgsearch)
CREATE INDEX idx_tools_description_pgsearch ON apig_mcp_tools USING gin(description gin_pgsearch_ops);

-- Create vector index
CREATE INDEX idx_tools_vector ON apig_mcp_tools USING ivfflat (vector vector_cosine_ops) WITH (lists = 100);
```

## Configuration Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| dsn | string | Yes | - | PostgreSQL database connection string |
| apiKey | string | Yes | - | API Key for the embedding service |
| baseURL | string | No | https://dashscope.aliyuncs.com/compatible-mode/v1 | Base URL for OpenAI-compatible embedding API |
| model | string | No | text-embedding-v4 | Vector model name |
| dimensions | int | No | 1024 | Vector dimensions |
| vectorWeight | float64 | No | 0.5 | Vector search weight (0-1), textWeight = 1 - vectorWeight |
| tableName | string | No | apig_mcp_tools | Database table name |
| description | string | No | Tool search server for semantic similarity search | Server description |

## Configuration Example

```json
{
  "dsn": "host=localhost user=postgres password=password dbname=mcp_tools port=5432 sslmode=disable",
  "apiKey": "your-api-key",
  "baseURL": "https://api.openai.com/v1",
  "model": "text-embedding-3-small",
  "dimensions": 1536,
  "vectorWeight": 0.6,
  "tableName": "apig_mcp_tools",
  "description": "MCP Tools Search Server"
}
```

## MCP Tool Interface

### x_higress_tool_search

**Description**: Performs semantic similarity comparison based on query statements to search for the most relevant tools

**Input Parameters**:
- `query` (string, required): Query statement for semantic similarity comparison with tool descriptions
- `topK` (integer, optional): Number of tools to select, defaults to 10

**Output**:
```json
{
  "tools": [
    {
      "name": "server_name___tool_name",
      "title": "Tool Title", 
      "description": "Tool description",
      "inputSchema": {...},
      "outputSchema": {...}
    }
  ]
}
```

Note: Each tool definition in the `tools` array directly uses the complete JSON content from the database `metadata` field, with only the `name` field modified to the format `${serverName}___${toolName}`.

## Hybrid Search Algorithm

This server uses the following hybrid search algorithm:

1. **Vector Search**: Uses cosine similarity to calculate similarity between query vector and tool description vectors
2. **Full-text Search**: Uses pgsearch extension's `@@@` operator for advanced full-text search
3. **Hybrid Scoring**: Final score = Full-text score × (1 - vectorWeight) + Vector similarity score × vectorWeight

### Full-text Search Implementation

Uses pgsearch extension's advanced full-text search capabilities:

```sql
SELECT 
    description @@@ pgsearch.config('text:query_text') AS score
FROM apig_mcp_tools
WHERE description @@@ pgsearch.config('text:query_text')
ORDER BY score
```

This approach provides more precise text matching and scoring mechanisms.

## Fallback Mechanism

When embedding generation fails, the system automatically falls back to text-only search to ensure service availability:

1. **Primary**: Hybrid search (vector + text)
2. **Fallback**: Text-only search when embedding API is unavailable
3. **Logging**: Comprehensive debug and error logging for troubleshooting

## Dependencies

- PostgreSQL database
- pgvector extension (for vector storage and similarity computation)
- pgsearch extension (for advanced full-text search)
- OpenAI-compatible embedding API access (supports OpenAI, DashScope, Azure OpenAI, etc.)

## Supported Embedding Services

This server supports any OpenAI-compatible embedding API, including:

- **OpenAI**: `https://api.openai.com/v1`
- **DashScope (Alibaba Cloud)**: `https://dashscope.aliyuncs.com/compatible-mode/v1`
- **Azure OpenAI**: `https://your-resource.openai.azure.com/openai/deployments/your-deployment`
- **Other OpenAI-compatible services**: Any service that implements the OpenAI embedding API format

Simply configure the `baseURL` and `apiKey` parameters according to your chosen service.