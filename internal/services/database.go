package services

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/dws-org/dws-event-service/configs"
	"github.com/dws-org/dws-event-service/internal/pkg/logger"
	"github.com/dws-org/dws-event-service/prisma/db"

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

// EnsureConnected ensures the Prisma client is connected, reconnecting if necessary
func (d *DatabaseService) EnsureConnected(ctx context.Context) error {
	if d.client == nil {
		return errors.New("prisma client is not initialized")
	}

	// Try to reconnect - Prisma Go client's Connect() can be called multiple times
	// If already connected, it will return an error, but we can ignore it
	// If not connected or connection was lost, it will reconnect
	if err := d.client.Prisma.Connect(); err != nil {
		// Check if error indicates already connected (this is fine)
		errStr := err.Error()
		if !strings.Contains(errStr, "already") && !strings.Contains(errStr, "connected") {
			d.logger.Errorf("failed to ensure database connection: %v", err)
			return err
		}
		// If already connected, that's fine - continue
		d.logger.Debugln("database connection already established")
	}

	return nil
}

// GetClient returns the Prisma client instance, ensuring it's connected first
func (d *DatabaseService) GetClient() *db.PrismaClient {
	// Ensure connection is active before returning client
	ctx := context.Background()
	if err := d.EnsureConnected(ctx); err != nil {
		d.logger.Warnf("warning: database connection check failed: %v", err)
	}
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
