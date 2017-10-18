# kafka-rest GO KAFKA client over HTTP

Simplifies access to Kafka exposing an HTTP(S) endpoint for posting messages.

# Endpoint

    POST /topic/<topicName>?single=false
    body

The `topicName` is the name of the topic to push the messages

The body are the current message / messages to push.

The `single` argument indicates if the body should be considered a single message (`true`) of multiple
(`false`, default).

If the Content-Type is "text/text" each line is taken as a message (`single=true`) or as the body of a single
message (`single=false`).

If the Content-Type is "application/json" the JSON values is the message (`single=true`) 
or an array is expected, if `single=false`, and each element is a message.

