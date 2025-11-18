package services

import (
	"context"
	"errors"
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
	if err := d.client.Prisma.Connect(); err != nil {
		d.logger.Errorln("cound't connect to prisma, ", err)
	}
	d.logger.Infoln("connected to database")
}
func (d *DatabaseService) DbDisconnect() {
	if err := d.client.Prisma.Disconnect(); err != nil {
		d.logger.Errorln("cound't close connection to prisma, ", err)
	}
}

// HealthCheck validates that the Prisma client is initialized. Extend this method
// with domain-specific checks (for example, SELECT 1 queries) once your schema
// is in place.
func (d *DatabaseService) HealthCheck(ctx context.Context) error {
	if d.client == nil {
		return errors.New("prisma client is not initialized")
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
