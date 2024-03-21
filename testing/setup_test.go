package testing

import (
	"context"
	"database/sql"
	"excel-receiver/config"
	"excel-receiver/http/api"
	"excel-receiver/provider"
	"excel-receiver/repository"
	"excel-receiver/service"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-stomp/stomp/v3"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type testSvc struct {
	suite.Suite
	ctx     context.Context
	log     provider.ILogger
	artemis *stomp.Conn
	app     *api.App
	db      *sqlx.DB
}

func (suite *testSvc) SetupSuite() {
	suite.ctx = context.Background()

	suite.loadConfig()

	logger := provider.NewLogger()
	suite.log = logger

	suite.initArtemis()
	suite.initMysql()
	suite.initSchema()
	suite.initAPP()
}

func (suite *testSvc) loadConfig() {
	t := suite.T()

	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	assert.NoError(t, err)

	conf := config.Config{}
	err = viper.Unmarshal(&conf)
	assert.NoError(t, err)
	config.Configuration = conf
}

func (suite *testSvc) initAPP() {
	tokenRepo := repository.NewToken(suite.db, suite.log)
	requestRepo := repository.NewRequest(suite.db, suite.log)
	queueRepo := repository.NewQueueArtemis(suite.artemis, suite.log)

	svc := service.NewSendRequestService(suite.log, queueRepo, requestRepo)

	app := api.NewApp(suite.log, svc, tokenRepo)
	suite.app = app

	t := suite.T()

	server, err := suite.app.CreateServer("localhost:5050")
	assert.NoError(t, err)

	go func() {
		assert.NoError(t, server.ListenAndServe())
	}()

	time.Sleep(2 * time.Second)
}

func (suite *testSvc) initArtemis() {
	ctx := context.Background()

	t := suite.T()

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "docker.io/apache/activemq-artemis:2.30.0-alpine",
			Env: map[string]string{
				"ARTEMIS_USER":     "artemis",
				"ARTEMIS_PASSWORD": "artemis",
				"AMQ_USER":         "artemis",
				"AMQ_PASSWORD":     "artemis",
			},
			ExposedPorts: []string{"61616/tcp", "8161/tcp"},
			WaitingFor: wait.ForAll(
				wait.ForLog("Server is now live"),
				wait.ForLog("REST API available"),
			),
		},
		Started: true,
	}

	artemisContainer, err := testcontainers.GenericContainer(ctx, req)
	require.NoError(t, err)

	host, err := artemisContainer.Host(ctx)
	if err != nil {
		require.NoError(t, err)
	}

	port, err := artemisContainer.MappedPort(suite.ctx, "61616")
	assert.NoError(t, err)

	host = fmt.Sprintf("%s:%s", host, port.Port())

	conn, err := stomp.Dial("tcp", host, stomp.ConnOpt.Login("artemis", "artemis"))
	if err != nil {
		require.NoError(t, err)
	}

	suite.artemis = conn

	fmt.Println("ADDRESS:", config.Configuration.Artemis.Address)
}

func (suite *testSvc) initMysql() {
	t := suite.T()

	mysqlC, err := testcontainers.GenericContainer(suite.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mysql:8.0",
			ExposedPorts: []string{"3306/tcp"},
			Env: map[string]string{
				"MYSQL_ROOT_PASSWORD": "root",
				"MYSQL_DATABASE":      "excel_rec_db",
				// "MYSQL_USER":          "root",
				// "MYSQL_PASSWORD": "root",
			},
			WaitingFor: wait.ForLog("mysqld: ready for connections."),
		},
		Started: true,
	})
	require.NoError(t, err)

	if mysqlC == nil {
		t.Fatalf("MySQL container is nil")
	}

	logOutput, _ := mysqlC.Logs(suite.ctx)
	t.Logf("MySQL container logs: %s", logOutput)

	host, err := mysqlC.Host(suite.ctx)
	require.NoError(t, err)
	port, err := mysqlC.MappedPort(suite.ctx, "3306")
	require.NoError(t, err)

	url := (&url.URL{
		User:     url.UserPassword("root", "root"),
		Host:     fmt.Sprintf("tcp(%s:%s)", host, port.Port()),
		Path:     "excel_rec_db",
		RawQuery: strings.Join([]string{"multiStatements=true", "parseTime=true"}, "&"),
	}).String()

	url = strings.TrimLeft(url, "/")

	dbConn, err := sql.Open("mysql", url)
	require.Nil(t, err)

	db := sqlx.NewDb(dbConn, "mysql")
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		time.Sleep(5 * time.Second)
		err = db.Ping()
		if err == nil {
			break
		}
	}

	suite.db = db
}

func (suite *testSvc) initSchema() {
	t := suite.T()

	sqlFilePath := "sql/schema.sql"
	sqlContent, err := os.ReadFile(sqlFilePath)
	if err != nil {
		t.Fatal(err)
	}

	_, err = suite.db.ExecContext(suite.ctx, string(sqlContent))
	if err != nil {
		t.Fatal(err)
	}
}
