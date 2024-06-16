package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/awhdesmond/revolut-user-service/pkg/api"
	"github.com/awhdesmond/revolut-user-service/pkg/common"
	"github.com/awhdesmond/revolut-user-service/pkg/users"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	cfgFlagHost        = "host"
	cfgFlagPort        = "port"
	cfgFlagMetricsPort = "metrics-port"
	cfgFlagLogLevel    = "log-level"
	cfgFlagCORSOrigin  = "cors-origin"

	cfgFlagPostgresHost     = "postgres-host"
	cfgFlagPostgresPort     = "postgres-port"
	cfgFlagPostgresDatabase = "postgres-database"
	cfgFlagPostgresUsername = "postgres-username"
	cfgFlagPostgresPassword = "postgres-password"

	envVarPrefix = "REVOLUT_USERS_SVC"

	defaultApiPort     = "8080"
	defaultMetricsPort = "8081"
	defaultLogLevel    = "info"
	defaultCORSOrigin  = "*"
)

type ServerConfig struct {
	common.PostgresSQLConfig `mapstructure:",squash"`

	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	MetricsPort string `mapstructure:"metrics-port"`
	CORSOrigin  string `mapstructure:"cors-origin"`
}

func (cfg ServerConfig) HTTPBindAddress() string {
	return fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
}

func (cfg ServerConfig) HTTPMetricsBindAddress() string {
	return fmt.Sprintf("%s:%s", cfg.Host, cfg.MetricsPort)
}

func main() {
	viper.SetDefault(cfgFlagHost, "")
	viper.SetDefault(cfgFlagPort, defaultApiPort)
	viper.SetDefault(cfgFlagMetricsPort, defaultMetricsPort)
	viper.SetDefault(cfgFlagLogLevel, defaultLogLevel)
	viper.SetDefault(cfgFlagCORSOrigin, defaultCORSOrigin)

	viper.SetDefault(cfgFlagPostgresHost, "")
	viper.SetDefault(cfgFlagPostgresPort, "")
	viper.SetDefault(cfgFlagPostgresDatabase, "")
	viper.SetDefault(cfgFlagPostgresUsername, "")
	viper.SetDefault(cfgFlagPostgresPassword, "")

	viper.SetEnvPrefix(envVarPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// Logger

	logger, _ := common.InitZap(viper.GetString(cfgFlagLogLevel))
	defer logger.Sync()
	stdLog := zap.RedirectStdLog(logger)
	defer stdLog()

	var srvCfg ServerConfig
	if err := viper.Unmarshal(&srvCfg); err != nil {
		logger.Panic("config unmarshal failed", zap.Error(err))
	}

	logger.Info("server configuration", zap.String("config", fmt.Sprintf("%+v", srvCfg)))

	// Make Server
	apiSrv, err := makeAPIServer(srvCfg, logger)
	if err != nil {
		logger.Panic("error initialising api server", zap.Error(err))
	}

	// Run Servers
	go func() {
		logger.Info("starting api server",
			zap.String("host", srvCfg.Host),
			zap.String("port", srvCfg.Port),
			zap.String("commit", common.GitCommit),
			zap.String("version", common.Version),
		)
		if err := http.ListenAndServe(srvCfg.HTTPBindAddress(), apiSrv.Handler); err != nil {
			logger.Panic("error starting api server", zap.Error(err))
			os.Exit(1)
		}
	}()

	go func() {
		logger.Info("starting metrics server")
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(srvCfg.HTTPMetricsBindAddress(), nil)
	}()

	// graceful shutdown
	stopCh := api.SetupSignalHandler()
	sd, _ := api.NewShutdown(logger)
	sd.Graceful(stopCh, apiSrv)
}

func makeAPIServer(cfg ServerConfig, logger *zap.Logger) (*http.Server, error) {
	pgSess, err := common.NewPostgresDBSession(cfg.PostgresSQLConfig)
	if err != nil {
		return nil, err
	}

	store := users.NewStore(pgSess, logger)
	svc := users.NewDefaultService(store)
	handler := users.MakeHandler(svc)

	r := mux.NewRouter()
	securityMW := api.NewSecureHeadersMiddleware(cfg.CORSOrigin)
	wrwMW := api.NewWrappedReponseWriterMiddleware()
	loggingMW := api.NewMuxLoggingMiddleware(logger)
	metricsMW := api.NewMetricsMiddleware()

	r.Use(securityMW.Handler)
	r.Use(wrwMW.Handler)
	r.Use(loggingMW.Handler)
	r.Use(metricsMW.Handler)

	r.HandleFunc("/healthz", api.HealthzHandler)
	r.PathPrefix("/hello").Handler(handler)

	return &http.Server{Handler: r, Addr: cfg.HTTPBindAddress()}, nil
}
