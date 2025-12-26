package controllers

import "github.com/dws-org/dws-event-service/internal/services"

func GetDatabaseSeviceInstance() *services.DatabaseService {
	return services.GetDatabaseSeviceInstance()
}
