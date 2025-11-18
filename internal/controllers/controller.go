package controllers

import "github.com/oskargbc/dws-event-service.git/internal/services"

func GetDatabaseSeviceInstance() *services.DatabaseService {
	return services.GetDatabaseSeviceInstance()
}
