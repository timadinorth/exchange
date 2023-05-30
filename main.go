package main

import (
	"github.com/timadinorth/bet-exchange/api"
)

// @title BetPub exchange API
// @version 1.0
// @description Betting exchange API documentation
// @contact.name Tim Adi
// @contact.email timadinorth@gmail.com
// @BasePath /api/v1
func main() {
	s := &api.Server{}
	s.InitLogger()
	s.LoadConfig(".")
	s.ConnectDB()
	s.ConnectCache()
	s.InitWeb()
	s.RegisterRoutes()
	err := s.SetupModels()
	if err != nil {
		s.Log.Fatal("Failed to run db migrations")
	}
	s.Log.Info("starting...")
	s.Start()
}
