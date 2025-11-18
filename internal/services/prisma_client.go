package services

import (
	"context"

	"github.com/oskargbc/dws-event-service.git/prisma/db"
)

type PrismaClient struct {
	client *db.PrismaClient
}

// NewPrismaClient creates a new Prisma client instance
func NewPrismaClient() (*PrismaClient, error) {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return nil, err
	}

	return &PrismaClient{client: client}, nil
}

// Close closes the Prisma client connection
func (p *PrismaClient) Close() error {
	return p.client.Prisma.Disconnect()
}

// GetClient returns the underlying Prisma client
func (p *PrismaClient) GetClient() *db.PrismaClient {
	return p.client
}

// HealthCheck performs a health check on the database
// Note: Implement based on your schema once models are defined
func (p *PrismaClient) HealthCheck(ctx context.Context) error {
	// Basic health check - connection is verified during Connect()
	// Add specific model queries here once you have models defined
	return nil
}
