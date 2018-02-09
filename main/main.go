package main

import (
  log "github.com/sirupsen/logrus"
  "net/http"
  "github.com/ant0ine/go-json-rest/rest"
  "strconv"
  "os"
  "fmt"
)

const AppName = "APP_NAME"
const EnvironmentEnv = "ENVIRONMENT"
const PortEnv = "PORT"

const BrokerEnv = "BROKER"

var Environment  string

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

  broker := getStrEnv(BrokerEnv, "localhost:1111")
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

  stack := [] rest.Middleware{
    &AccessLogMiddleware{ Logger: log.StandardLogger(), IgnoredPathPrefix: "/system" },
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

  log.
    WithFields(log.Fields{"app": APP_NAME, "version": VERSION, "port": port, "broker": broker}).
    Info(fmt.Sprintf("Started %v v%v on port %v with %v broker", APP_NAME, VERSION, port, broker))
  log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), nil))
}


