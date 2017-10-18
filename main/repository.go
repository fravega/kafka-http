package main

type Repository interface {
  Push(topic string, message []byte) error
  Stat() interface{}
  Close()
}
