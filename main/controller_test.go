package main

import (
  "testing"
  "github.com/ant0ine/go-json-rest/rest"
  "github.com/stretchr/testify/assert"
)

func TestSendSingleTextMessage(*testing.T) {

}

func TestSendMultipleTextMessage(*testing.T) {

}

func TestSendSingleJsonMessage(*testing.T) {

}

func TestSendMultipleJsonMessage(*testing.T) {

}

func TestSendInvalidJsonMessage(*testing.T) {

}

func TestSendInvalidSingleParameter(*testing.T) {

}

func TestSendInvalidContentType(t *testing.T) {
  c, repo := makeController()
  w, r := mockHttpCall()

  c.ProduceMessages(w, r)

  if len(repo.Sent) > 0 {
    assert.Fail(t, "Messages should not be sent")
  }
}

func makeController() (*Controller, RepositoryMock) {
  r := RepositoryMock{}

  return NewController(r), r
}

func mockHttpCall() (w rest.ResponseWriter, r *rest.Request) {
  return nil, nil
}

type Event struct {
  Topic   string
  Msg     []byte
}

type RepositoryMock struct {
  Sent    []Event
}

func (r RepositoryMock) Push(topic string, message []byte) error {
  r.Sent = append(r.Sent, Event{topic, message })

  return nil  // No error
}

func (r RepositoryMock) Stat() interface{} {
  panic("Not implemented")
}

func (r RepositoryMock) Close() {
  // do nothing
}
