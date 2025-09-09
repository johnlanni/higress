package toolsearch

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/envoyproxy/envoy/contrib/golang/common/go/api"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBClient handles PostgreSQL database connections and operations
type DBClient struct {
	db         *gorm.DB
	dsn        string
	tableName  string
	userID     string
	reconnect  chan struct{}
	stop       chan struct{}
	panicCount int32
}

// ToolRecord represents a tool record in the database
type ToolRecord struct {
	ID          string `gorm:"primaryKey"`
	ServerName  string `gorm:"column:server_name"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
	Metadata    string `gorm:"column:metadata;type:text"`
	UserID      string `gorm:"column:user_id"`
}

// NewDBClient creates a new DBClient instance
func NewDBClient(dsn, tableName, userID string, stop chan struct{}) *DBClient {
	client := &DBClient{
		dsn:       dsn,
		tableName: tableName,
		userID:    userID,
		reconnect: make(chan struct{}, 1),
		stop:      stop,
	}

	// Start reconnection goroutine
	go client.reconnectLoop()

	// Try initial connection
	if err := client.connect(); err != nil {
		api.LogErrorf("Initial database connection failed: %v", err)
	}

	return client
}

func (c *DBClient) connect() error {
	gormConfig := gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	db, err := gorm.Open(postgres.Open(c.dsn), &gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	c.db = db
	return nil
}

func (c *DBClient) reconnectLoop() {
	defer func() {
		if r := recover(); r != nil {
			api.LogErrorf("Recovered from panic in reconnectLoop: %v", r)

			// Increment panic counter
			atomic.AddInt32(&c.panicCount, 1)

			// If panic count exceeds threshold, stop trying to reconnect
			if atomic.LoadInt32(&c.panicCount) > 3 {
				api.LogErrorf("Too many panics in reconnectLoop, stopping reconnection attempts")
				return
			}

			// Wait for a while before restarting
			time.Sleep(5 * time.Second)

			// Restart the reconnect loop
			go c.reconnectLoop()
		}
	}()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.stop:
			api.LogInfof("Database connection closed")
			return
		case <-ticker.C:
			if c.db == nil || c.Ping() != nil {
				if err := c.connect(); err != nil {
					api.LogErrorf("Database reconnection failed: %v", err)
				} else {
					api.LogInfof("Database reconnected successfully")
					atomic.StoreInt32(&c.panicCount, 0)
				}
			}
		case <-c.reconnect:
			if err := c.connect(); err != nil {
				api.LogErrorf("Database reconnection failed: %v", err)
			} else {
				api.LogInfof("Database reconnected successfully")
				atomic.StoreInt32(&c.panicCount, 0)
			}
		}
	}
}

func (c *DBClient) reconnectIfDbEmpty() error {
	if c.db == nil {
		select {
		case c.reconnect <- struct{}{}:
		default:
		}
		return fmt.Errorf("database is not connected, attempting to reconnect")
	}
	return nil
}

func (c *DBClient) handleSQLError(err error) error {
	if err != nil {
		select {
		case c.reconnect <- struct{}{}:
		default:
		}
		return fmt.Errorf("failed to execute SQL: %w", err)
	}
	return nil
}

// Ping checks database connectivity
func (c *DBClient) Ping() error {
	if c.db == nil {
		return fmt.Errorf("database connection is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	sqlDB, err := c.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying *sql.DB: %v", err)
	}

	return sqlDB.PingContext(ctx)
}

// SearchTools performs hybrid search (vector + full-text) on tools
func (c *DBClient) SearchTools(query string, vector []float32, topK int, vectorWeight, textWeight float64) ([]ToolRecord, error) {
	if err := c.reconnectIfDbEmpty(); err != nil {
		return nil, err
	}

	// Convert vector to string array format for PostgreSQL
	vectorParts := make([]string, len(vector))
	for i, v := range vector {
		vectorParts[i] = fmt.Sprintf("%f", v)
	}
	vectorStr := "{" + strings.Join(vectorParts, ",") + "}"

	// Build the hybrid search SQL query with parameterized queries
	var sql string
	var args []interface{}

	if c.userID != "" {
		sql = `
		WITH t1 AS (
			SELECT
				id,
				server_name,
				name,
				description,
				metadata,
				description @@@ pgsearch.config(CONCAT('description:', ?::text)) AS score,
				2 AS source
			FROM ` + c.tableName + `
			WHERE user_id = ?
			ORDER BY score ASC
			LIMIT ?
		),
		t2 AS (
			SELECT
				id,
				server_name,
				name,
				description,
				metadata,
				cosine_similarity(vector, ?::real[]) AS score,
				1 AS source
			FROM ` + c.tableName + `
			WHERE vector IS NOT NULL AND user_id = ?
			ORDER BY vector <-> ?
			LIMIT ?
		)
		SELECT 
			COALESCE(t1.id, t2.id) as id,
			COALESCE(t1.server_name, t2.server_name) as server_name,
			COALESCE(t1.name, t2.name) as name,
			COALESCE(t1.description, t2.description) as description,
			COALESCE(t1.metadata, t2.metadata) as metadata,
			COALESCE(ABS(t1.score), 0.0) * ? + COALESCE(t2.score, 0.0) * ? AS hybrid_score
		FROM t1
		FULL OUTER JOIN t2 ON t1.id = t2.id 
		ORDER BY hybrid_score DESC
		LIMIT ?`
		args = []interface{}{query, c.userID, topK, vectorStr, c.userID, vectorStr, topK, textWeight, vectorWeight, topK}
	} else {
		sql = `
		WITH t1 AS (
			SELECT
				id,
				server_name,
				name,
				description,
				metadata,
				description @@@ pgsearch.config(CONCAT('description:', ?::text)) AS score,
				2 AS source
			FROM ` + c.tableName + `
			WHERE 1=1
			ORDER BY score ASC
			LIMIT ?
		),
		t2 AS (
			SELECT
				id,
				server_name,
				name,
				description,
				metadata,
				cosine_similarity(vector, ?::real[]) AS score,
				1 AS source
			FROM ` + c.tableName + `
			WHERE vector IS NOT NULL
			ORDER BY vector <-> ?
			LIMIT ?
		)
		SELECT 
			COALESCE(t1.id, t2.id) as id,
			COALESCE(t1.server_name, t2.server_name) as server_name,
			COALESCE(t1.name, t2.name) as name,
			COALESCE(t1.description, t2.description) as description,
			COALESCE(t1.metadata, t2.metadata) as metadata,
			COALESCE(ABS(t1.score), 0.0) * ? + COALESCE(t2.score, 0.0) * ? AS hybrid_score
		FROM t1
		FULL OUTER JOIN t2 ON t1.id = t2.id 
		ORDER BY hybrid_score DESC
		LIMIT ?`
		args = []interface{}{query, topK, vectorStr, vectorStr, topK, textWeight, vectorWeight, topK}
	}

	api.LogInfof("Executing hybrid search SQL query for: '%s'", query)
	api.LogDebugf("SQL: %s", sql)
	api.LogDebugf("Args: %v", args)

	rows, err := c.db.Raw(sql, args...).Rows()
	if err := c.handleSQLError(err); err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ToolRecord
	for rows.Next() {
		var record ToolRecord
		var hybridScore float64

		err := rows.Scan(
			&record.ID,
			&record.ServerName,
			&record.Name,
			&record.Description,
			&record.Metadata,
			&hybridScore,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		results = append(results, record)
	}

	api.LogInfof("Hybrid search completed, found %d results", len(results))
	return results, nil
}

// SearchToolsTextOnly performs text-only search when embedding fails
func (c *DBClient) SearchToolsTextOnly(query string, topK int) ([]ToolRecord, error) {
	api.LogDebugf("Starting text-only search for query: '%s'", query)

	if err := c.reconnectIfDbEmpty(); err != nil {
		return nil, err
	}

	// Build text-only search SQL query with parameterized queries
	var sql string
	var args []interface{}

	if c.userID != "" {
		sql = `
		SELECT 
			id,
			server_name,
			name,
			description,
			metadata,
			description @@@ pgsearch.config(CONCAT('description:', ?)) AS score
		FROM ` + c.tableName + `
		WHERE user_id = ?
		ORDER BY score ASC
		LIMIT ?`
		args = []interface{}{query, c.userID, topK}
	} else {
		sql = `
		SELECT 
			id,
			server_name,
			name,
			description,
			metadata,
			description @@@ pgsearch.config(CONCAT('description:', ?)) AS score
		FROM ` + c.tableName + `
		ORDER BY score ASC
		LIMIT ?`
		args = []interface{}{query, topK}
	}

	api.LogInfof("Executing text-only search SQL query for: '%s'", query)
	api.LogDebugf("SQL: %s", sql)
	api.LogDebugf("Args: %v", args)

	rows, err := c.db.Raw(sql, args...).Rows()
	if err := c.handleSQLError(err); err != nil {
		api.LogErrorf("Text-only search SQL failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	var results []ToolRecord
	for rows.Next() {
		var record ToolRecord
		var score float64

		err := rows.Scan(
			&record.ID,
			&record.ServerName,
			&record.Name,
			&record.Description,
			&record.Metadata,
			&score,
		)
		if err != nil {
			api.LogErrorf("Failed to scan text-only search result: %v", err)
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		results = append(results, record)
	}

	api.LogInfof("Text-only search completed, found %d results", len(results))
	return results, nil
}

// GetAllTools retrieves all tools from the database
func (c *DBClient) GetAllTools() ([]ToolRecord, error) {
	if err := c.reconnectIfDbEmpty(); err != nil {
		return nil, err
	}

	api.LogInfof("Executing GetAllTools query from table: %s", c.tableName)

	var tools []ToolRecord
	query := c.db.Table(c.tableName)

	// Add user_id filter if userID is provided (for ADB PostgreSQL)
	if c.userID != "" {
		query = query.Where("user_id = ?", c.userID)
		api.LogDebugf("Added user_id filter: %s", c.userID)
	}

	err := query.Find(&tools).Error
	if err := c.handleSQLError(err); err != nil {
		return nil, err
	}

	api.LogInfof("GetAllTools query completed, found %d tools", len(tools))
	return tools, nil
}
