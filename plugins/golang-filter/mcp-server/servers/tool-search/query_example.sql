-- Hybrid search query examples
-- This query demonstrates how to use vector similarity search and pgsearch BM25 full-text search for hybrid ranking

-- Create necessary indexes (run once during setup):
-- Vector index using ann method:
-- CREATE INDEX idx_tools_vector ON apig_mcp_tools USING ann(vector) WITH (dim = 1024, algorithm = hnswflat, distancemeasure = L2, vector_include = 0);
-- Full-text search index using pgsearch BM25:
-- CALL pgsearch.create_bm25(index_name => 'idx_tools_description_bm25', table_name => 'apig_mcp_tools', text_fields => '{description: {}}');
-- Note: In pgsearch.config('field_name:query'), field_name must match the indexed field

-- Example query: Search for tools related to "weather data"
WITH t1 AS (
    -- Full-text search part: Use pgsearch BM25 for text matching
    SELECT
        id,
        server_name,
        name,
        description,
        metadata,
        vector,
        description @@@ pgsearch.config('description:weather data') AS score,
        2 AS source
    FROM apig_mcp_tools
    ORDER BY score ASC
    LIMIT 10
),
t2 AS (
    -- Vector search part: Use cosine similarity for semantic matching
    -- Note: Query text needs to be converted to vector via OpenAI-compatible embedding API first
    SELECT
        id,
        server_name,
        name,
        description,
        metadata,
        vector,
        cosine_similarity(vector, ARRAY[0.1,0.2,0.3,0.4,0.5,0.6,0.7,0.8,0.9,1]::real[]) AS score,
        1 AS source
    FROM apig_mcp_tools
    WHERE vector IS NOT NULL
    ORDER BY vector <-> ARRAY[0.1,0.2,0.3,0.4,0.5,0.6,0.7,0.8,0.9,1]
    LIMIT 10
)
-- Hybrid results: Combine scores from full-text search and vector search
SELECT 
    COALESCE(t1.id, t2.id) as id,
    COALESCE(t1.server_name, t2.server_name) as server_name,
    COALESCE(t1.name, t2.name) as name,
    COALESCE(t1.description, t2.description) as description,
    COALESCE(t1.metadata, t2.metadata) as metadata,
    COALESCE(t1.vector, t2.vector) as vector,
    -- Hybrid scoring: text score weight * 0.2 + vector score weight * 0.8
    -- Weights can be adjusted based on business requirements
    COALESCE(ABS(t1.score), 0.0) * 0.2 + COALESCE(t2.score, 0.0) * 0.8 AS hybrid_score
FROM t1
FULL OUTER JOIN t2 ON t1.id = t2.id 
ORDER BY hybrid_score
LIMIT 10;

-- Test full-text search separately using BM25
SELECT 
    id,
    server_name,
    name,
    description,
    description @@@ pgsearch.config('description:weather data') AS bm25_score
FROM apig_mcp_tools
ORDER BY bm25_score ASC
LIMIT 5;

-- Test vector search separately using cosine similarity
SELECT 
    id,
    server_name,
    name,
    description,
    cosine_similarity(vector, ARRAY[0.1,0.2,0.3,0.4,0.5,0.6,0.7,0.8,0.9,1]::real[]) AS similarity_score
FROM apig_mcp_tools
WHERE vector IS NOT NULL
ORDER BY vector <-> ARRAY[0.1,0.2,0.3,0.4,0.5,0.6,0.7,0.8,0.9,1]
LIMIT 5;
