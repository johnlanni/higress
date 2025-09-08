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
-- CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS pgsearch;

-- Note: Table name is configurable via the tableName parameter
CREATE TABLE apig_mcp_tools (
    id VARCHAR(255) PRIMARY KEY,
    server_name VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    metadata TEXT,
    vector real[]
);

-- Create vector index using ann method with HNSW algorithm
CREATE INDEX idx_tools_vector ON apig_mcp_tools USING ann(vector) 
WITH (dim = 1024, algorithm = hnswflat, distancemeasure = L2, vector_include = 0);

-- Create full-text search index using pgsearch BM25
CALL pgsearch.create_bm25(
    index_name => 'idx_tools_description_bm25',
    table_name => 'apig_mcp_tools',
    text_fields => '{description: {}}'
);
```

## Configuration Parameters

### Root Level Configuration

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| vector | object | Yes | - | Vector database configuration (see [Vector Configuration](#vector-configuration) below) |
| embedding | object | Yes | - | Embedding API configuration (see [Embedding Configuration](#embedding-configuration) below) |
| description | string | No | Tool search server for semantic similarity search | Server description |

### Vector Configuration

Configuration object for the `vector` field:

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| type | string | Yes | - | Vector database type (currently only "adb-postgres" supported) |
| dsn | string | Yes | - | PostgreSQL database connection string |
| vectorWeight | float64 | No | 0.5 | Vector search weight (0-1), textWeight = 1 - vectorWeight |
| tableName | string | No | apig_mcp_tools | Database table name |

### Embedding Configuration

Configuration object for the `embedding` field:

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| apiKey | string | Yes | - | API Key for the embedding service |
| baseURL | string | No | https://dashscope.aliyuncs.com/compatible-mode/v1 | Base URL for OpenAI-compatible embedding API |
| model | string | No | text-embedding-v4 | Vector model name |
| dimensions | int | No | 1024 | Vector dimensions |

## Configuration Example

```json
{
  "vector": {
    "type": "adb-postgres",
    "vectorWeight": 0.5,
    "tableName": "apig_mcp_tools",
    "dsn": "host=localhost user=postgres password=password dbname=mcp_tools port=5432 sslmode=disable"
  },
  "embedding": {
    "apiKey": "your-dashscope-api-key",
    "baseURL": "https://dashscope.aliyuncs.com/compatible-mode/v1",
    "model": "text-embedding-v4",
    "dimensions": 1024
  },
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

1. **Vector Search**: Uses cosine similarity with HNSW algorithm for efficient semantic similarity search
2. **Full-text Search**: Uses pgsearch BM25 algorithm for advanced text matching and ranking
3. **Hybrid Scoring**: Final score = BM25 score × (1 - vectorWeight) + Vector similarity score × vectorWeight

### Vector Search Implementation

Uses PostgreSQL's ann index with HNSW algorithm for efficient vector similarity search:

```sql
-- Vector similarity using cosine similarity
SELECT 
    cosine_similarity(vector, query_vector) AS similarity_score
FROM apig_mcp_tools
WHERE vector IS NOT NULL
ORDER BY vector <-> query_vector
```

### Full-text Search Implementation

Uses pgsearch extension's BM25 algorithm for advanced full-text search:

```sql
SELECT 
    description @@@ pgsearch.config('description:query_text') AS bm25_score
FROM apig_mcp_tools
ORDER BY bm25_score ASC
```

BM25 provides state-of-the-art text ranking based on term frequency and document frequency.

**Important for ADB PostgreSQL**: In ADB PostgreSQL, BM25 scores are negative values where smaller (more negative) values indicate higher relevance. Therefore, we use `ORDER BY score ASC` to get the most relevant results first. No WHERE clause is needed since all results have valid BM25 scores.

## Fallback Mechanism

When embedding generation fails, the system automatically falls back to text-only search to ensure service availability:

1. **Primary**: Hybrid search (vector + text)
2. **Fallback**: Text-only search when embedding API is unavailable
3. **Logging**: Comprehensive debug and error logging for troubleshooting

## Dependencies

- Alibaba Cloud ADB PostgreSQL database
- pgsearch extension (for advanced full-text search)
- OpenAI-compatible embedding API access (supports OpenAI, DashScope, Azure OpenAI, etc.)

**Note**: This implementation is specifically designed for Alibaba Cloud ADB PostgreSQL, which uses:
- Native `real[]` array type for vector storage
- BM25 scoring with negative values (smaller values indicate higher relevance)

## Supported Embedding Services

This server supports any OpenAI-compatible embedding API, including:

- **OpenAI**: `https://api.openai.com/v1`
- **DashScope (Alibaba Cloud)**: `https://dashscope.aliyuncs.com/compatible-mode/v1`
- **Azure OpenAI**: `https://your-resource.openai.azure.com/openai/deployments/your-deployment`
- **Other OpenAI-compatible services**: Any service that implements the OpenAI embedding API format

Simply configure the `baseURL` and `apiKey` parameters according to your chosen service.