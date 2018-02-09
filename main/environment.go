package main

import (
  log "github.com/sirupsen/logrus"
  "strings"
  "os"
  "strconv"
)

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
