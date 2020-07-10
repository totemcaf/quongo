# quongo
Simple Queue implemented over MongoDB in GO, using a REST delivery.

![Mongo](images/quongo_640.png)

# tl;dr

* message queue system
* runs stand-alone (download)
* simple RESTFull delivery
* fully asynchronous implementation, no blocking calls
* Selectable backend: Memory, MongoDB, Redis

# Configuration

For Redis backend:

* REDIS_URI: The URI in Redis URI format of the provided server(s)
* REDIS_PREFIX: The prefix for all keys used by this implementation. If the prefix
does not end in ':' it will addded.

For MongoDB backend:

* MONGO_URL: List of URLs separated by command of mongodb hosts. defaults to: "localhost:27017"
* MONGO_USER: User name to connect as, defaults to "quongo"
* MONGO_PWD: User password to connect as, defaults to "quongosecret"
* MONGO_DB: Name of database to use, default to "quongo"
* QUONGO_PORT: Port where the Quongo app is going to listen

# Sample use

Send a message:

    curl
# API

### Create a queue.
This call is not necessary, the queue will be created automatically with the first push of a message.
Use this message if you want to tune parameters.

    POST /v1/queues
    {
      "name": "queue name"
      "visibilityWindow": "TBD",
      "duplicationWindow": "TBD",
      "duplicate": "ignore",
      "target": {
        "url": "http://api/v1/event",
        "maxEventsPerRequest": 1,
        "method": "POST",
        "timeoutMs": 200
      },
      "additionalHeaderes": [
        "X-HeaderName1",
        "X-HeaderName2"
      ]
    }

* name: Name of the queue.
* visibilityWindow: Time to keep requested message until it is put in the queue again. Expressed in ms. This is the mechanism to ensure if the consumer fails to process the message and dit not send the acknowlede or free the messagge, it will not starve for ever in the queue.
* duplicationWindow: Time the deleted message is recorded to ignore new message with same id. Used to avoid duplicated messages when the processing time is less than the time taken to send the same message again.
* duplicate: In case same message is stored twice, do *ignore* or *update* or send *error*. *Ignore* just discard the new message keeping the oldest in the queue. *Udapte* will replace the original message with the new one, keeping the newest. *Error* will report a duplicate message.
* target: If configured, the Quongo will spawn a consumer for the queue that will pop messages from the queue and send (Using the method indicated, defaults to HTTP PUT) to the provided URL. The message ID will be added to the URL where the patter "{mid}" is found. If the target does not respond, an exponential backoff will be used.
* additionalHeaderes: These headers, if present in message push request, will be send with the message on respons or push.

Responds with queue description.

### Get queue data

    GET /v1/queues/:queueId

Sample answer:

    {
      "name": "myFirstQueue"
      "visibilityWindow": 2000
    }

### Push new message (idempotent)

    PUT /v1/queues/:queueId/messages/:messageId

The request body is the message payload.

Why I need a message ID? This way if you send the same message twice, it is updated.

Note: If the message was consumed before the first post, it will added again. Si (de)duplication window.

Response:

  200 Acepted

Header: X-MId will have the message id

### Push new message (no idempotent)

    POST /v1/queues/:queueId/messages

The request body is the message payload.

Note: This format allows to have duplicated messaged. It is recommended to use the format
with Id.

### Get specific message (without consuming)

    GET /v1/queues/:queueId/messages/:mId

### Get next message from queue and temporarily remove it

    GET /v1/queues/:queueId/messages/pop

The ##X-Ack## header value should be used to confirm that message was processed.

### Get next count messages from queue and temporarily remove them

    GET /v1/queues/:queueId/messages/pop-many?limit=:count

The ##X-Ack## header value should be used to confirm that messages were processed.

### Keep alive

    PATCH /v1/queues/:queueId/messages/:messageId/ack/:ackId

If the message still is pending of acknoledeegement, delays its visibility window.
This request is used by a process that takes to long to process the message so it avoids the
message be stolen by another consumer.

### Acknoledge message was processed

    DELETE /v1/queues/:queueId/messages/:messageId/ack/:ackId

### Get next message from queue with auto acknoledge

    DELETE /v1/queues/:queueId/messages/pop

### Get several messages from queue with auto acknoledge

    DELETE /v1/queues/:queueId/messages/pop-many?limit=:count

## Responses

### Queue response

This is the response from the queue endpoints:

    {
      "id": "58c85bcc3c00003f00ead573",
      "name": "queue name",
      "total": 2343,
      "hidden": 87989,
      "inProcess": 345,
      "created": "2017-02-01T12:00:45.004+0003"
    }

### Message response

These are the response headers from messages endpoints:

    X-MId: 58c85bcc3c00003f00ead573
    X-CId: 58c85bcc3c00003c00ead589
    X-GId: 767bbc767767d7ae6f67b77a
    X-Created: 2017-03-14T21:08:28+0000
    Retry-After: 2017-03-14T21:08:28+0000
    Expires: 2017-03-14T24:08:28+0000
    X-Ack: 58c85bcc3c00004200ead57d
    Last-Modified: 2017-03-14T21:08:28+0000
    X-Holder: 58c85bcc3c00003f00ead57a
    X-Deleted:

Where:

* **Retry-After**: The date time since the message will be available to pick up. The visible is set to a time in the future
  in case of message delay, and where a message is being processed by a client (in its visibility window).
* **X-CId**: A correlation id provided by message producer to associate the message to the business transaction.
* **X-GId**: A group id. Used to keep only one message per group. (See [Message group](#message-groups))
* **X-Holder**: An id provided by the client that picked up the message
* **X-Ack**: An id generated at the moment of picking a message and used to acknowledge it,
  or widening the visibility window (in case the client processing requires more time).
* **X-Deleted**: If present, the message was already processed, but not still deleted.
* **Response body**: The message data.

### Message groups

A message group allows to remove related messages from the queue when a new message is pushed. For
example if a queue holds events on entities, the entity Id can be used as the **grouping id** so
only the last event is stored in the queue.

As always, already processed messages does not longer exist in the queue, so they are not taken into account.

# Note

I'm using [Git Flow](https://danielkummer.github.io/git-flow-cheatsheet/) to handle the repository operations.

Also following Clear Architecture as describe [here](https://manuel.kiessling.net/2012/09/28/applying-the-clean-architecture-to-go-applications/).

https://godoc.org/github.com/mongodb/mongo-go-driver/mongo
https://godoc.org/github.com/mongodb/mongo-go-driver/mongo/clientopt

# TODO

## v0.1

[ ] Allow to send a message with ID
[ ] Allow to consume a message
[ ] Allow to acknowledeg a message
[ ] Memory repository

## v0.2

[ ] Allow to create a queue
[ ] JSONLog

## v0.3

[ ] Redis repository

## v??
[ ] Allow auto acknowledge. The message is aknowledge in the moment it is requested.
[ ] Extended deduplication. Allow to deduplicate consummed messaged (at least for a windowed time)
[ ] Allow PUSH mode with target
[ ] Add target no responding configuration


https://www.reddit.com/r/golang/comments/h7o2lt/fiber_vs_echo_vs_gin_pros_and_cons/
https://github.com/labstack/echo
https://github.com/gofiber/fiber
https://github.com/avelino/awesome-go#web-frameworks
