// +heroku goVersion 1.15

module github.com/bradpurchase/grocerytime-backend

go 1.15

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/denisenkom/go-mssqldb v0.0.0-20200910202707-1e08a3fab204 // indirect
	github.com/go-gormigrate/gormigrate/v2 v2.0.0
	github.com/gofrs/uuid v3.3.0+incompatible // indirect
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/graphql-go/graphql v0.7.9
	github.com/graphql-go/handler v0.2.3
	github.com/joho/godotenv v1.3.0
	github.com/kr/text v0.2.0 // indirect
	github.com/lib/pq v1.8.0 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.8.2-0.20190227000051-27936f6d90f9 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/sendgrid/rest v2.6.2+incompatible // indirect
	github.com/sendgrid/sendgrid-go v3.7.0+incompatible
	github.com/stretchr/testify v1.6.1
	github.com/tidwall/gjson v1.6.1
	github.com/trevex/graphql-go-subscription v0.0.0-20170731204342-4a0a4158754b
	golang.org/x/crypto v0.0.0-20201012173705-84dcc777aaee
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	gorm.io/driver/postgres v1.0.2
	gorm.io/gorm v1.20.2
)
