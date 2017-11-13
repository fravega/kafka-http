package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"strings"
	"time"
)

type AccessLogMiddleware struct {

	// Logger points to the logger object used by this middleware, it defaults to
	// log.New(os.Stderr, "", 0).
	Logger *log.Logger

	// This path should not be logged. Usually the `health` request is hidden
	IgnoredPathPrefix string
}

// MiddlewareFunc makes AccessLogApacheMiddleware implement the Middleware interface.
func (mw *AccessLogMiddleware) MiddlewareFunc(h rest.HandlerFunc) rest.HandlerFunc {

	// set the default Logger
	if mw.Logger == nil {
		mw.Logger = log.StandardLogger()
	}

	if mw.IgnoredPathPrefix == "" {
		mw.IgnoredPathPrefix = "*" // Path does start with this
	}

	return func(w rest.ResponseWriter, r *rest.Request) {

		// call the handler
		h(w, r)

		path := r.URL.Path

		if !strings.HasPrefix(path, mw.IgnoredPathPrefix) {
			util := &accessLogUtil{w, r}

			fields := log.Fields{
				"bytes":        util.BytesWritten(),
				"responseTime": util.ResponseTime(),
				"remoteAddr":   util.ApacheRemoteAddr(),
				"protocol":     util.R.Proto,
				"method":       util.R.Method,
				"processId":    util.Pid(),
				"queryStr":     util.ApacheQueryString(),
				// "request": util.R.Method + " " + util.R.URL.RequestURI + " " + util.R.Proto,
				"statusCode": util.StatusCode(),
				"startTime":  util.StartTime().Format("2006-01-02T15:04:05.999-0700"),
				"remoteUser": util.RemoteUser(),
				"User-Agent": util.R.UserAgent(),
				"Referer":    util.R.Referer(),
			}

			mw.Logger.WithFields(fields).Info(path)
		}
	}
}

// Note: Extracted from  "github.com/ant0ine/go-json-rest/rest" to adapt the Logger.

// accessLogUtil provides a collection of utility functions that derive data from the Request object.
// This object is used to provide data to the Apache Style template and the the JSON log record.
type accessLogUtil struct {
	W rest.ResponseWriter
	R *rest.Request
}

// As stored by the auth middlewares.
func (u *accessLogUtil) RemoteUser() string {
	if u.R.Env["REMOTE_USER"] != nil {
		return u.R.Env["REMOTE_USER"].(string)
	}
	return ""
}

// If qs exists then return it with a leadin "?", apache log style.
func (u *accessLogUtil) ApacheQueryString() string {
	if u.R.URL.RawQuery != "" {
		return "?" + u.R.URL.RawQuery
	}
	return ""
}

// When the request entered the timer middleware.
func (u *accessLogUtil) StartTime() *time.Time {
	if u.R.Env["START_TIME"] != nil {
		return u.R.Env["START_TIME"].(*time.Time)
	}
	return nil
}

// If remoteAddr is set then return is without the port number, apache log style.
func (u *accessLogUtil) ApacheRemoteAddr() string {
	remoteAddr := u.R.RemoteAddr
	if remoteAddr != "" {
		if ip, _, err := net.SplitHostPort(remoteAddr); err == nil {
			return ip
		}
	}
	return ""
}

// As recorded by the recorder middleware.
func (u *accessLogUtil) StatusCode() int {
	if u.R.Env["STATUS_CODE"] != nil {
		return u.R.Env["STATUS_CODE"].(int)
	}
	return 0
}

// As mesured by the timer middleware.
func (u *accessLogUtil) ResponseTime() *time.Duration {
	if u.R.Env["ELAPSED_TIME"] != nil {
		return u.R.Env["ELAPSED_TIME"].(*time.Duration)
	}
	return nil
}

// Process id.
func (u *accessLogUtil) Pid() int {
	return os.Getpid()
}

// As recorded by the recorder middleware.
func (u *accessLogUtil) BytesWritten() int64 {
	if u.R.Env["BYTES_WRITTEN"] != nil {
		return u.R.Env["BYTES_WRITTEN"].(int64)
	}
	return 0
}
