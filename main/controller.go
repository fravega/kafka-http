package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/buger/jsonparser"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

type Controller struct {
	repository Repository
	Routes     []*rest.Route
}

func NewController(repository Repository) *Controller {
	ctrl := Controller{repository: repository}

	routes := []*rest.Route{
		//    rest.Get("/v1/job", ctrl.QueueAll),
		rest.Post("/v1/topics/#topicName", ctrl.ProduceMessages),
	}

	ctrl.Routes = routes

	return &ctrl
}

// POST /api/v1/topics/
func (c *Controller) ProduceMessages(w rest.ResponseWriter, r *rest.Request) {
	query := r.URL.Query()

	var err error
	single, err := parseBoolean(query.Get("single"), false)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	topic := r.PathParam("topicName")

	contentType := r.Header.Get("Content-Type")
	var messages [][]byte

	var mtype string

	switch contentType {
	case "text/text":
		messages, err = extractTextMessages(r, single)
		mtype = "text"

	case "application/json":
		messages, err = extractJsonMessages(r, single)
		mtype = "json"

	default:
		err = errors.New("Unhandled content type: " + contentType)
	}

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send messages
	for _, message := range messages {
		c.repository.Push(topic, message) // TODO Check result !!!!!!!
	}

	// Report result
	var plural = "s"

	if len(messages) == 1 {
		plural = ""
	}

	w.WriteJson(fmt.Sprintf("Sent %v %v message%v to %v", len(messages), mtype, plural, topic))
	log.WithFields(log.Fields{"type": mtype, "count": len(messages), "topic": topic}).Info("messages sent")
}

func extractTextMessages(r *rest.Request, single bool) ([][]byte, error) {
	body, err := parseBodyAsTextBytes(r)

	if err != nil {
		return nil, err
	}

	if single {
		return [][]byte{body}, nil
	} else {
		return bytes.Split(body, []byte("\n")), nil
	}
}

func extractJsonMessages(r *rest.Request, single bool) ([][]byte, error) {
	body, err := parseBodyAsTextBytes(r)

	if err != nil {
		return nil, err
	}

	if single {
		if json.Valid(body) {
			return [][]byte{body}, nil
		} else {
			return nil, errors.New("invalid JSON")
		}
	} else {

		values := [][]byte{}

		jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, innerErr error) {
			if innerErr != nil {
				err = innerErr
			} else {
				values = append(values, value)
			}
		})

		return values, nil
	}
}

func parseBoolean(paramName string, defaultValue bool) (bool, error) {
	paramName = strings.ToLower(strings.TrimSpace(paramName))

	switch paramName {
	case "":
		return defaultValue, nil
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return defaultValue, errors.New("invalid boolean value " + paramName)
	}
}

func parseBodyAsTextBytes(r *rest.Request) ([]byte, error) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return []byte{}, err
	}

	return body, err
}
