package database

import (
	"fmt"
	"service/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSslMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %w", err)
	}

	return db, nil
}

// func connectMongoDB() {
// 	godotenv.Load("../.env")
// 	mongoUri := os.Getenv("MONGO")
// 	// Establish a connection to the MongoDB database
// 	clientOptions := options.Client().ApplyURI(mongoUri).SetMaxPoolSize(5).SetMinPoolSize(1).SetServerSelectionTimeout(5 * time.Second)
// 	// Connect to MongoDB
// 	client, err := mongo.Connect(context.Background(), clientOptions)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Check the connection
// 	err = client.Ping(context.Background(), nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Connected to MongoDB!")
// }
