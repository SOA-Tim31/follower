package data

import (
	"Rest/domain"
	"context"
	"log"
	"strconv"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type UserRepository struct {
	// Thread-safe instance which maintains a database connection pool
	driver neo4j.DriverWithContext
	logger *log.Logger
}

// NoSQL: Constructor which reads db configuration from environment and creates a keyspace
func New(logger *log.Logger) (*UserRepository, error) {
	// Local instance
	uri := "bolt://neo4j:7687"
	user := "neo4j"
	pass := "nekaSifra"
	auth := neo4j.BasicAuth(user, pass, "")

	driver, err := neo4j.NewDriverWithContext(uri, auth)
	if err != nil {
		logger.Panic(err)
		return nil, err
	}

	// Return repository with logger and DB session
	return &UserRepository{
		driver: driver,
		logger: logger,
	}, nil
}

// Check if connection is established
func (mr *UserRepository) CheckConnection() {
	ctx := context.Background()
	err := mr.driver.VerifyConnectivity(ctx)
	if err != nil {
		mr.logger.Panic(err)
		return
	}
	// Print Neo4J server address
	mr.logger.Printf(`Neo4J server address: %s`, mr.driver.Target().Host)
}

// Disconnect from database
func (mr *UserRepository) CloseDriverConnection(ctx context.Context) {
	mr.driver.Close(ctx)
}

func (mr *UserRepository) WriteUser(user *domain.User) error {
	// Neo4J Sessions are lightweight so we create one for each transaction (Cassandra sessions are not lightweight!)
	// Sessions are NOT thread safe
	ctx := context.Background()
	session := mr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	// ExecuteWrite for write transactions (Create/Update/Delete)
	savedPerson, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"create (u:User) SET u.Id = $id, u.Username = $username  return u.Username + ', from node ' + id(u)",
				map[string]any{"id": user.Id, "username": user.Username})
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				return result.Record().Values[0], nil
			}

			return nil, result.Err()
		})
	if err != nil {
		mr.logger.Println("Error inserting Person:", err)
		return err
	}
	mr.logger.Println(savedPerson.(string))
	return nil
}

func (mr *UserRepository) GetRecommendations(userId string) ([]*domain.User, error) {
	ctx := context.Background()
	session := mr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	intUserId, err := strconv.Atoi(userId) // Convert string userID to integer
	if err != nil {
		mr.logger.Println("Error converting userId to integer:", err)
		return nil, err
	}

	userResults, err := session.ExecuteRead(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				`MATCH (u:User)-[:IS_FOLLOWING]->(f:User)-[:IS_FOLLOWING]->(r:User)
				WHERE u.Id = $userId AND NOT (u)-[:IS_FOLLOWING]->(r) AND r.Username <> u.Username
				RETURN DISTINCT r.Id AS id, r.Username AS username
				`,
				map[string]any{"userId": intUserId}) // Use converted integer ID
			if err != nil {
				return nil, err
			}

			var users []*domain.User
			for result.Next(ctx) {
				record := result.Record()
				id, _ := record.Get("id")
				username, _ := record.Get("username")
				users = append(users, &domain.User{
					Id:       int32(id.(int64)), // Corrected line
					Username: username.(string),
				})
			}
			if result.Err() != nil {
				return nil, result.Err()
			}
			return users, nil
		})
	if err != nil {
		mr.logger.Println("Error querying search:", err)
		return nil, err
	}

	return userResults.([]*domain.User), nil // Type assertion to a slice of pointers to domain.User
}

func (mr *UserRepository) FollowUser(user *domain.User, userToFollow *domain.User) error {
	ctx := context.Background()
	session := mr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"match (user1:User), (user2:User) where user1.Username = $username AND user2.Username = $followUsername create (user1)-[r: IS_FOLLOWING]->(user2) return type(r)",
				map[string]any{"username": user.Username, "followUsername": userToFollow.Username})
			if err != nil {
				return nil, err
			}
			if result.Next(ctx) {
				return result.Record().Values[0], nil
			}
			return nil, result.Err()
		})
	if err != nil {
		mr.logger.Println("Error inserting following:", err)
		return err
	}

	return nil
}
