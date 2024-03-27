package options

import (
	"crypto/tls"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/robfig/cron/v3"
	"net/http"
	"v1/pkg/apiserver"
	"v1/pkg/client/cache"
	"v1/pkg/client/mysql"
	"v1/pkg/logger"
	"v1/pkg/model"
	genericoptions "v1/pkg/server/options"
	"v1/pkg/token"

	cliflag "k8s.io/component-base/cli/flag"
)

var jwtSecret = []byte("6322ecab234de98bf206aff5b216d1e4") // MD5 ("cspm")

type ServerRunOptions struct {
	GenericServerRunOptions *genericoptions.ServerRunOptions
	RDBOptions              *mysql.Options
	LoggerOptions           *logger.Options

	DebugMode bool
}

func NewServerRunOptions() *ServerRunOptions {
	s := &ServerRunOptions{
		GenericServerRunOptions: genericoptions.NewServerRunOptions(),
		RDBOptions:              mysql.NewMysqlOptions(mysql.SetDefaultRdbDbname("graduation_project")),
		LoggerOptions:           logger.NewLoggerOptions(),
	}

	return s
}

func (s *ServerRunOptions) Flags() (fss cliflag.NamedFlagSets) {
	fs := fss.FlagSet("generic")
	fs.BoolVar(&s.DebugMode, "debug", s.DebugMode, "Don't enable this if you don't know what it means.")
	s.GenericServerRunOptions.AddFlags(fs)
	s.RDBOptions.AddFlags(fss.FlagSet("rdb"))
	s.LoggerOptions.AddFlags(fss.FlagSet("log"))

	return fss
}

// NewAPIServer creates an APIServer instance using given options
func (s *ServerRunOptions) NewAPIServer(stopCh <-chan struct{}) (*apiserver.APIServer, error) {

	var (
		apiServer = &apiserver.APIServer{
			TokenManager: token.NewJWTTokenManager(jwtSecret, jwt.SigningMethodHS256),
			CacheClient:  cache.NewSimpleCache(),
			Crontab:      cron.New(),
			// Sched:        scan.NewScheduler(),
		}
	)

	logger.InitLogger(s.LoggerOptions)

	// connect to mysql
	if s.RDBOptions != nil {
		apiServer.RDBClient = mysql.NewMysqlClient(s.RDBOptions)

		// AutoMigrate tables
		_ = apiServer.RDBClient.AutoMigrate(
			new(model.User),
			new(model.Class),
			new(model.Profession),
			new(model.Project),
			new(model.Config),
			new(model.Resume),
			new(model.Company),
			new(model.AuditLog),
			new(model.College),
			new(model.Interview),
		)
	}

	// apiServer.Sched = scan.NewScheduler()

	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d", s.GenericServerRunOptions.BindAddress, s.GenericServerRunOptions.Port),
	}

	if s.GenericServerRunOptions.TlsPrivateKey != "" && s.GenericServerRunOptions.TlsCertFile != "" {
		certificate, err := tls.LoadX509KeyPair(s.GenericServerRunOptions.TlsCertFile, s.GenericServerRunOptions.TlsPrivateKey)
		if err != nil {
			return nil, err
		}

		server.TLSConfig = &tls.Config{Certificates: []tls.Certificate{certificate}}
	}

	apiServer.Server = server
	return apiServer, nil
}
