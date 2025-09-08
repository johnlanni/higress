package toolsearch

import (
	"context"
	"fmt"
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
	reconnect  chan struct{}
	stop       chan struct{}
	panicCount int32
}

// ToolRecord represents a tool record in the database
type ToolRecord struct {
	ID          int       `gorm:"primaryKey"`
	ServerName  string    `gorm:"column:server_name"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
	Metadata    string    `gorm:"column:metadata;type:jsonb"`
	Vector      []float32 `gorm:"column:vector;type:vector(1024)"`
}

// NewDBClient creates a new DBClient instance
func NewDBClient(dsn, tableName string, stop chan struct{}) *DBClient {
	client := &DBClient{
		dsn:       dsn,
		tableName: tableName,
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

	// Convert vector to PostgreSQL array format
	vectorStr := "ARRAY["
	for i, v := range vector {
		if i > 0 {
			vectorStr += ","
		}
		vectorStr += fmt.Sprintf("%f", v)
	}
	vectorStr += "]::real[]"

	// Build the hybrid search SQL query
	sql := fmt.Sprintf(`
		WITH t1 AS (
			SELECT
				id,
				server_name,
				name,
				description,
				metadata,
				vector,
				description @@@ pgsearch.config('text:%s') AS score,
				2 AS source
			FROM %s
			WHERE description @@@ pgsearch.config('text:%s')
			ORDER BY score
			LIMIT $2
		),
		t2 AS (
			SELECT
				id,
				server_name,
				name,
				description,
				metadata,
				vector,
				cosine_similarity(vector, %s) AS score,
				1 AS source
			FROM %s
			WHERE vector IS NOT NULL
			ORDER BY vector <-> %s
			LIMIT $2
		)
		SELECT 
			COALESCE(t1.id, t2.id) as id,
			COALESCE(t1.server_name, t2.server_name) as server_name,
			COALESCE(t1.name, t2.name) as name,
			COALESCE(t1.description, t2.description) as description,
			COALESCE(t1.metadata, t2.metadata) as metadata,
			COALESCE(t1.vector, t2.vector) as vector,
			COALESCE(ABS(t1.score), 0.0) * $3 + COALESCE(t2.score, 0.0) * $4 AS hybrid_score
		FROM t1
		FULL OUTER JOIN t2 ON t1.id = t2.id 
		ORDER BY hybrid_score DESC
		LIMIT $2
	`, query, c.tableName, query, vectorStr, c.tableName, vectorStr)

	rows, err := c.db.Raw(sql, query, topK, textWeight, vectorWeight).Rows()
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
			&record.Vector,
			&hybridScore,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		results = append(results, record)
	}

	return results, nil
}

// SearchToolsTextOnly performs text-only search when embedding fails
func (c *DBClient) SearchToolsTextOnly(query string, topK int) ([]ToolRecord, error) {
	api.LogDebugf("Starting text-only search for query: '%s'", query)

	if err := c.reconnectIfDbEmpty(); err != nil {
		return nil, err
	}

	// Build text-only search SQL query
	sql := fmt.Sprintf(`
		SELECT 
			id,
			server_name,
			name,
			description,
			metadata,
			vector,
			description @@@ pgsearch.config('text:%s') AS score
		FROM %s
		WHERE description @@@ pgsearch.config('text:%s')
		ORDER BY score DESC
		LIMIT $1
	`, query, c.tableName, query)

	api.LogDebugf("Executing text-only search SQL")

	rows, err := c.db.Raw(sql, topK).Rows()
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
			&record.Vector,
			&score,
		)
		if err != nil {
			api.LogErrorf("Failed to scan text-only search result: %v", err)
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		results = append(results, record)
	}

	api.LogDebugf("Text-only search found %d results", len(results))
	return results, nil
}

// GetAllTools retrieves all tools from the database
func (c *DBClient) GetAllTools() ([]ToolRecord, error) {
	if err := c.reconnectIfDbEmpty(); err != nil {
		return nil, err
	}

	var tools []ToolRecord
	err := c.db.Table(c.tableName).Find(&tools).Error
	if err := c.handleSQLError(err); err != nil {
		return nil, err
	}

	return tools, nil
}
