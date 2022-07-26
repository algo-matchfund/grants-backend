package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-openapi/loads"
	"github.com/jessevdk/go-flags"
	"github.com/rs/cors"

	"github.com/algo-matchfund/grants-backend/api/handlers"
	"github.com/algo-matchfund/grants-backend/api/middlewares"
	"github.com/algo-matchfund/grants-backend/gen/restapi"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/config"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/algo-matchfund/grants-backend/internal/keycloak"
	"github.com/algo-matchfund/grants-backend/internal/service/watchdog"
	"github.com/algo-matchfund/grants-backend/internal/smartcontract"
)

func main() {
	defer func() {
		// an alternative to os.Exit, which respects defer calls
		err := recover()
		if err != nil {
			if code, ok := err.(int); ok {
				os.Exit(code)
			}

			panic(err) // not an exit code, bubble up panic
		}
	}()

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Println(err)
		panic(1)
	}

	api := operations.NewGrantsProgramAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	configFile := flag.String("config", "config.yml", "Path to the config file")

	flag.Parse()

	cfg := config.Config{}
	err = cfg.LoadConfig(*configFile)

	if err != nil {
		log.Printf("Error loading config: %s\n", err)
		panic(1)
	}

	keycloak := keycloak.NewKeycloakService(&cfg)

	db, err := database.NewGrantsDatabase(
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)

	if err != nil {
		log.Println(err)
		panic(1)
	}

	defer db.Close()

	err = db.Migrate()

	if err != nil {
		log.Println(err)
		panic(1)
	}

	watchdogFactory := watchdog.NewWatchdogFactory(&cfg, db)
	algorandWatchdog, err := watchdogFactory.GetWatchdog("algorand")
	if err != nil {
		log.Println(err)
		panic(1)
	}
	scc, err := smartcontract.NewSmartContractClient(&cfg)
	if err != nil {
		log.Println(err)
		panic(1)
	}

	// Project related endpoints
	api.GetProjectsHandler = handlers.NewGetProjectsHandler(db)
	api.GetProjectsIDHandler = handlers.NewGetProjectsIDHandler(db)
	api.GetProjectIDFundTxHandler = handlers.NewGetProjectIDFundTxHandler(db, watchdogFactory)
	api.GetProjectCategoriesHandler = handlers.NewGetProjectCategoriesHandler(db)
	api.GetProjectContributorsHandler = handlers.NewGetProjectContributorsHandler(db)
	api.GetProjectMatchCalculationHandler = handlers.NewGetProjectMatchCalculationHandler(db)
	api.GetProjectNewsHandler = handlers.NewGetProjectNewsHandler(db)
	api.GetProjectQAHandler = handlers.NewGetProjectQAHandler(db)
	api.PostProjectsHandler = handlers.NewPostProjectsHandler(db)
	api.PostProjectIDFundHandler = handlers.NewPostProjectIDFundHandler(db)
	api.PostProjectIDFundTxHandler = handlers.NewPostProjectIDFundTxHandler(db, watchdogFactory)
	api.PostProjectQuestionHandler = handlers.NewPostProjectQuestionHandler(db)
	api.DeleteQuestionHandler = handlers.NewDeleteQuestionHandler(db)
	api.PostProjectNewsHandler = handlers.NewPostProjectNewsHandler(db)
	api.UpdateProjectNewsItemHandler = handlers.NewUpdateProjectNewsItemHandler(db)
	api.DeleteProjectsIDNewsNewsIDHandler = handlers.NewDeleteProjectsIDNewsNewsIDHandler(db)
	api.PostProjectAnswerHandler = handlers.NewPostProjectAnswerHandler(db)
	api.UpdateAnswerHandler = handlers.NewUpdateAnswerHandler(db)
	api.DeleteAnswerHandler = handlers.NewDeleteAnswerHandler(db)
	api.PostProjectModerationByIDHandler = handlers.NewPostProjectModerationByIDHandler(db, scc)
	api.GetProjectsForModerationHandler = handlers.NewGetProjectsForModerationHandler(db)
	api.GetProjectModerationByIDHandler = handlers.NewGetProjectModerationByIDHandler(db)
	api.GetStatsHandler = handlers.NewGetStatsHandler(db)
	api.GetSmartContractTransactionsHandler = handlers.NewGetSmartContractHandler(scc)

	// // Public user related endpoints
	api.GetUsersUserIDHandler = handlers.NewGetUserUserIdHandler(db, keycloak)

	// // Authenticated user handlers
	api.GetUsersHandler = handlers.NewGetUsersHandler(db, keycloak)
	api.PostUsersHandler = handlers.NewPostUserHandler(db, keycloak)
	api.GetUsersSettingsHandler = handlers.NewGetUsersSettingsHandler(db, keycloak)
	api.PutUsersSettingsHandler = handlers.NewPutUsersSettingsHandler(db, keycloak)
	api.GetUsersNotificationsHandler = handlers.NewGetUsersNotificationsHandler(db)
	api.PutUsersNotificationsHandler = handlers.NewPutUsersNotificationsHandler(db)
	api.DeleteUsersNotificationsHandler = handlers.NewDeleteUsersNotificationsHandler(db)
	api.PutUsersNotificationsNotificationIDHandler = handlers.NewPutUsersNotificationsNotificationIDHandler(db)
	api.DeleteUsersNotificationsNotificationIDHandler = handlers.NewDeleteUsersNotificationsNotificationIDHandler(db)

	api.BearerAuth, err = middlewares.NewAuthenticator(cfg.Authentication.PublicKey)

	if err != nil {
		log.Println(err)
		panic(1)
	}

	server.ConfigureAPI()

	log.Printf("Setting allowed origins to %s\n", cfg.Server.AllowedOrigins)
	c := cors.New(cors.Options{
		AllowedOrigins: cfg.Server.AllowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"authorization", "content-type", "origin"},
	})

	server.SetHandler(c.Handler(server.GetHandler()))
	// server.SetHandler(middlewares.NewPanicMiddleware(c.Handler(server.GetHandler())))

	server.Host = cfg.Server.Host

	if cfg.Server.Ssl {
		server.TLSCertificate = flags.Filename(cfg.Server.Cert)
		server.TLSCertificateKey = flags.Filename(cfg.Server.Key)
		server.TLSPort = cfg.Server.Port
		server.EnabledListeners = []string{"https"}
	} else {
		server.Port = cfg.Server.Port
		server.EnabledListeners = []string{"http"}
	}

	// goroutine to update matches table periodically
	go func() {
		for {
			err := db.UpdateMatches()
			if err != nil {
				fmt.Println("Failed to update matches:", err)
			}
			time.Sleep(10 * time.Second)
		}
	}()

	cancelContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start watchdogs
	algorandWatchdog.StartWatch(cancelContext)

	if err := server.Serve(); err != nil {
		log.Println(err)
		panic(1)
	}
}
