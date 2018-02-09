package main

import (
  log "github.com/sirupsen/logrus"
  "github.com/bshuster-repo/logrus-logstash-hook"
  "net"
)

const LogstashServerEnv = "LOGSTASH_SERVER"
const Production = "production"

func initLogstash(appName string, address string) error {
  conn, err := net.Dial("tcp", address)
  if err != nil {
    log.Fatal(err)
    return err
  }

  hook := logrustash.New(conn, logrustash.DefaultFormatter(log.Fields{"type": appName}))

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

