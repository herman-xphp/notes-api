package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectMySQL connects to MySQL using environment variables
func ConnectMySQL() *gorm.DB {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// Validate required environment variables
	if host == "" || user == "" || name == "" || port == "" {
		log.Fatal("❌ Missing required database environment variables (DB_HOST, DB_USER, DB_NAME, DB_PORT)")
	}

	return ConnectMySQLWithParams(host, user, pass, name, port)
}

// ConnectMySQLWithParams connects to MySQL with explicit parameters
func ConnectMySQLWithParams(host, user, pass, name, port string) *gorm.DB {
	// Validate required parameters
	if host == "" || user == "" || name == "" || port == "" {
		log.Fatal("❌ Missing required database parameters")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name)

	// Retry logic with exponential backoff
	var db *gorm.DB
	var err error
	maxRetries := 5
	retryDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			break
		}

		if i < maxRetries-1 {
			log.Printf("⚠️ Failed to connect MySQL (attempt %d/%d): %v. Retrying in %v...", i+1, maxRetries, err, retryDelay)
			time.Sleep(retryDelay)
			retryDelay *= 2 // Exponential backoff
		}
	}

	if err != nil {
		log.Fatal("❌ Failed to connect MySQL after retries:", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("❌ Failed to get database instance:", err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection with retry
	pingRetryDelay := 2 * time.Second
	for i := 0; i < maxRetries; i++ {
		if err := sqlDB.Ping(); err == nil {
			break
		}
		if i < maxRetries-1 {
			log.Printf("⚠️ Failed to ping database (attempt %d/%d). Retrying...", i+1, maxRetries)
			time.Sleep(pingRetryDelay)
			pingRetryDelay *= 2
		} else {
			log.Fatal("❌ Failed to ping database after retries")
		}
	}

	log.Println("✅ MySQL connected")
	return db
}
