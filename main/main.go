package main

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/bshuster-repo/logrus-logstash-hook"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const AppName = "APP_NAME"
const EnvironmentEnv = "ENVIRONMENT"
const PortEnv = "PORT"
const LogstashServerEnv = "LOGSTASH_SERVER"

const BrokerEnv = "BROKER"

const Production = "production"

var Environment string

func initLogstash(appName string, address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
		return err
	}

	hook := logrustash.New(conn, logrustash.DefaultFormatter(log.Fields{"service_name": appName}))

	if err != nil {
		log.Fatal(err)
	}

	log.StandardLogger().Hooks.Add(hook)

	return nil
}

func initLogger() {
	if logstashServer := getStrEnv(LogstashServerEnv, ""); logstashServer != "" {
		initLogstash(getStrEnv(AppName, APP_NAME), logstashServer)
	}

	log.SetLevel(getLevelEnv("LOG_LEVEL", "INFO"))

	if Environment == Production {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		// The TextFormatter is default, you don't actually have to do this.
		log.SetFormatter(&log.TextFormatter{})
	}
}

func main() {
	defer func() {
		log.Info("Done everything")
		if x := recover(); x != nil {
			log.WithField("error", x).Error("Run time panic")
		}
	}()

	Environment = getStrEnv(EnvironmentEnv, "development")
	initLogger()

	log.Info("Starting ...")

	broker := getStrEnv(BrokerEnv, "localhost:9092")
	port := getIntEnv(PortEnv, 7075)

	var repository Repository
	repository, err := NewKafkaRepository(broker)

	if err != nil {
		log.Panic(err.Error())
		os.Exit(1)
	}

	controller := NewController(repository)

	statusMw := &rest.StatusMiddleware{}

	systemStat := NewSystemController(repository, statusMw)

	stack := []rest.Middleware{
		&AccessLogMiddleware{Logger: log.StandardLogger(), IgnoredPathPrefix: "/system"},
		//      Format: rest.CombinedLogFormat,
		&rest.TimerMiddleware{},
		&rest.RecorderMiddleware{},
		&rest.RecoverMiddleware{
			EnableResponseStackTrace: true,
		},
		&rest.JsonIndentMiddleware{},
		&rest.ContentTypeCheckerMiddleware{},
		&rest.GzipMiddleware{},
	}

	api := rest.NewApi()

	api.Use(statusMw)

	api.Use(stack...)

	routes := append(controller.Routes, systemStat.Routes...)

	router, err := rest.MakeRouter(routes...)

	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)

	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))

	log.WithFields(log.Fields{"app": APP_NAME, "version": VERSION, "port": port, "broker": broker}).
		Info(fmt.Sprintf("Started %v v%v on port %v with %v broker", APP_NAME, VERSION, port, broker))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func getIntEnv(name string, defValue int) int {
	v := strings.TrimSpace(os.Getenv(name))

	if v == "" {
		return defValue
	} else {
		value, err := strconv.Atoi(v)
		if err != nil {
			log.WithFields(log.Fields{"name": name, "value": v}).Error("The argument is not a number")
		}
		return value
	}
}

func getLevelEnv(name string, defValue string) log.Level {
	v := strings.TrimSpace(os.Getenv(name))

	if v == "" {
		v = defValue
	}

	level, err := log.ParseLevel(v)

	if err != nil {
		log.WithFields(log.Fields{"name": name, "value": v}).Error("The argument is not a log level")
		return log.InfoLevel
	}

	return level
}

func getStrEnv(name string, defValue string) string {
	v := strings.TrimSpace(os.Getenv(name))

	if v == "" {
		return defValue
	} else {
		return v
	}
}
