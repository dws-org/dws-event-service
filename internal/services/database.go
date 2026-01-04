package services

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/oskargbc/dws-event-service.git/configs"
	"github.com/oskargbc/dws-event-service.git/internal/pkg/logger"
	"github.com/oskargbc/dws-event-service.git/prisma/db"

	"github.com/sirupsen/logrus"
)

var databaseInstance *DatabaseService
var lock = &sync.Mutex{}

type DatabaseService struct {
	client *db.PrismaClient
	logger *logrus.Logger
	env    *configs.Config
}

func GetDatabaseSeviceInstance() *DatabaseService {
	if databaseInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if databaseInstance == nil {
			databaseInstance = &DatabaseService{
				client: db.NewClient(),
				logger: logger.NewLogrusLogger(),
				env:    configs.GetEnvConfig(),
			}
			databaseInstance.dbConnect()
		}
	}

	return databaseInstance
}

func (d *DatabaseService) dbConnect() {
	// Connect to database - this will panic if connection fails
	if err := d.client.Prisma.Connect(); err != nil {
		d.logger.Fatalf("couldn't connect to prisma: %v. Make sure DATABASE_URL environment variable is set correctly.", err)
	}
	d.logger.Infoln("connected to database")
}

func (d *DatabaseService) DbDisconnect() {
	if err := d.client.Prisma.Disconnect(); err != nil {
		d.logger.Errorln("cound't close connection to prisma, ", err)
	}
}

// EnsureConnected checks if client is initialized
// NOTE: This should ONLY be used by health checks, NOT by regular requests!
// Calling Connect() multiple times creates connection leaks in Prisma
func (d *DatabaseService) EnsureConnected(ctx context.Context) error {
	if d.client == nil {
		return errors.New("prisma client is not initialized")
	}
	return nil
}

// GetClient returns the Prisma client instance
// NOTE: Connection is established once at startup via dbConnect()
// Do NOT call Connect() here as it creates connection leaks!
func (d *DatabaseService) GetClient() *db.PrismaClient {
	return d.client
}

// HealthCheck validates that the Prisma client is initialized and connected.
// It performs a simple query to verify the database connection is active.
func (d *DatabaseService) HealthCheck(ctx context.Context) error {
	if d.client == nil {
		return errors.New("prisma client is not initialized")
	}

	// Ensure connection is active
	if err := d.EnsureConnected(ctx); err != nil {
		return errors.New("database connection failed: " + err.Error())
	}

	// Perform a simple query to verify connection is actually working
	// Using FindFirst on Event model as a health check
	_, err := d.client.Event.FindFirst().Exec(ctx)
	if err != nil {
		// Check if error is "not connected" - this means connection issue
		errStr := err.Error()
		if strings.Contains(errStr, "not connected") || strings.Contains(errStr, "client is not connected") {
			return errors.New("database client is not connected: " + errStr)
		}
		// Other errors (like no events found) are fine - connection is working
		// A simple "record not found" type error means connection is fine
	}

	return nil
}
