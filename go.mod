module github.com/algo-matchfund/grants-backend

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/Masterminds/squirrel v1.5.1
	github.com/algorand/go-algorand-sdk v1.14.1
	github.com/go-openapi/analysis v0.21.3 // indirect
	github.com/go-openapi/errors v0.20.2
	github.com/go-openapi/loads v0.21.1
	github.com/go-openapi/runtime v0.23.3
	github.com/go-openapi/spec v0.20.5
	github.com/go-openapi/strfmt v0.21.2
	github.com/go-openapi/swag v0.21.1
	github.com/go-openapi/validate v0.21.0
	github.com/golang-migrate/migrate/v4 v4.15.1
	github.com/jessevdk/go-flags v1.5.0
	github.com/lib/pq v1.10.4
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/rs/cors v1.8.0
	go.mongodb.org/mongo-driver v1.9.0 // indirect
	golang.org/x/net v0.0.0-20220425223048-2871e0cb64e4
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20220422013727-9388b58f7150 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/algo-matchfund/grants-backend/api/handlers => ./api/handlers
	github.com/algo-matchfund/grants-backend/api/middlewares => ./grants-backend/api/middlewares
	github.com/algo-matchfund/grants-backend/gen/restapi => ./grants-backend/gen/restapi
	github.com/algo-matchfund/grants-backend/gen/restapi/operations => ./grants-backend/gen/restapi/operations
	github.com/algo-matchfund/grants-backend/internal/config => ./grants-backend/internal/config
	github.com/algo-matchfund/grants-backend/internal/database => ./grants-backend/internal/database
	github.com/algo-matchfund/grants-backend/internal/roles => ./grants-backend/internal/roles
	github.com/algo-matchfund/grants-backend/internal/service => ./grants-backend/internal/service
	github.com/algo-matchfund/grants-backend/internal/smartcontract => ./grants-backend/internal/smartcontract
)
