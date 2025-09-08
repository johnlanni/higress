-- Create MCP tools search table
-- Requires pgsearch extension: CREATE EXTENSION IF NOT EXISTS pgsearch;
-- Note: pgvector extension not needed for real[] type

-- Note: Table name can be customized via tableName configuration parameter, defaults to apig_mcp_tools
CREATE TABLE IF NOT EXISTS apig_mcp_tools (
    id VARCHAR(255) PRIMARY KEY,
    server_name VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    metadata TEXT,
    vector real[]
);

-- Create vector index using ann method
CREATE INDEX IF NOT EXISTS idx_tools_vector 
ON apig_mcp_tools USING ann(vector) WITH (dim = 1024, algorithm = hnswflat, distancemeasure = L2, vector_include = 0);

-- Create full-text search index using pgsearch BM25
CALL pgsearch.create_bm25(
    index_name => 'idx_tools_description_bm25',
    table_name => 'apig_mcp_tools',
    text_fields => '{description: {}}'
);

-- Create composite indexes
CREATE INDEX IF NOT EXISTS idx_tools_server_name 
ON apig_mcp_tools (server_name);

CREATE INDEX IF NOT EXISTS idx_tools_name 
ON apig_mcp_tools (name);

-- Insert sample data
INSERT INTO apig_mcp_tools (id, server_name, name, description, metadata) VALUES
('weather_server:get_weather', 'weather_server', 'get_weather', 'Get current weather data for a location', '{"name": "get_weather", "title": "Weather Data Retriever", "description": "Get current weather data for a location", "inputSchema": {"type": "object", "properties": {"location": {"type": "string", "description": "City name or zip code"}}, "required": ["location"]}, "outputSchema": {"type": "object", "properties": {"temperature": {"type": "number", "description": "Temperature in celsius"}, "conditions": {"type": "string", "description": "Weather conditions description"}, "humidity": {"type": "number", "description": "Humidity percentage"}}, "required": ["temperature", "conditions", "humidity"]}}'),
('database_server:query', 'database_server', 'query', 'Execute a read-only SQL query', '{"name": "query", "title": "SQL Query Executor", "description": "Execute a read-only SQL query", "inputSchema": {"type": "object", "properties": {"sql": {"type": "string", "description": "The SQL query to execute"}}, "required": ["sql"]}, "outputSchema": {"type": "object", "properties": {"results": {"type": "array", "description": "Query results"}}, "required": ["results"]}}'),
('file_server:read_file', 'file_server', 'read_file', 'Read contents of a file', '{"name": "read_file", "title": "File Reader", "description": "Read contents of a file", "inputSchema": {"type": "object", "properties": {"path": {"type": "string", "description": "Path to the file to read"}}, "required": ["path"]}, "outputSchema": {"type": "object", "properties": {"content": {"type": "string", "description": "File content"}}, "required": ["content"]}}'),
('api_server:http_request', 'api_server', 'http_request', 'Make an HTTP request to an API endpoint', '{"name": "http_request", "title": "HTTP Request Handler", "description": "Make an HTTP request to an API endpoint", "inputSchema": {"type": "object", "properties": {"url": {"type": "string", "description": "URL to make request to"}, "method": {"type": "string", "description": "HTTP method", "enum": ["GET", "POST", "PUT", "DELETE"]}, "headers": {"type": "object", "description": "Request headers"}}, "required": ["url"]}, "outputSchema": {"type": "object", "properties": {"status": {"type": "number", "description": "HTTP status code"}, "body": {"type": "string", "description": "Response body"}}, "required": ["status", "body"]}}')
ON CONFLICT DO NOTHING;
