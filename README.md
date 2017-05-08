# quongo
Simple Queue implemented over MongoDB in GO, using a REST interface.

![Mongo](images/quongo_640.png)

# tl;dr

* message queue system
* runs stand-alone (download) or embedded
* simple RESTFull interface
* fully asynchronous implementation, no blocking calls

# API

### Create a queue.
This call is not necessary, the queue will be created automatically with the first push of a message. 
Use this message if you want to tune parameters.
 
    POST /v1/queue
    {
      "name": "queue name"
      "visibilityWindow": "TBD"
    }

Responds with queue description.

### Get queue data

    GET /v1/queue/:queueId

### Push new message (idempotent)

    PUT /v1/queue/:queueId/message/:messageId

Why I need a message ID? This way if you send the same message twice, it is updated.

Note: If the message was consumed before the first post, it will added again.

### Push new message (no idempotent)

    POST /v1/queue/:queueId/message

## Responses

### Queue response

This is the response from the queue services:

    {
      "id": "58c85bcc3c00003f00ead573",
      "name": "queue name",
      "total": 2343,
      "hidden": 87989,
      "inProcess": 345,
      "created": "2017-02-01T12:00:45.004+0003"
    }

### Message response

This is the response from messages services

    {
        "id": "58c85bcc3c00003f00ead573",
        "cid": "58c85bcc3c00003c00ead589",
        "gid": null,
        "created": "2017-03-14T21:08:28+0000",
        "visible": "2017-03-14T21:08:28+0000",
        "holder": "58c85bcc3c00003f00ead57a",
        "ack": "58c85bcc3c00004200ead57d",
        "deleted": null,
        "payLoad": {...}
    }

Where:

* **visible**: The date time since the message will be available to pick up. The visible is set to a future time
  in case of message delay, and where a message is being processed by a client (in its visibility window).
* **cid**: A correlation id provided by message producer to associate the message to the business transaction.
* **gid**: A group id. Used to keep only one message per group. (See [Message group](#message-groups))
* **holder**: An id provided by the client that picked up the message
* **ack**: An id generated at the moment of picking a message and used to acknowledge it, 
  or widening the visibility window (in case the client processing requires more time).
* **deleted**: If present, the message was already processed
* **payload**: The message data. It is a valid JSON value.

### Message groups

A message group allows to remove related messages from the queue when a new message is pushed. For
example if a queue holds events on entities, the entity Id can be used as the **grouping id** so only the
last event is stored in the queue.

As always, already processed messages does not longer exist in the queue, so they are not taken into account.

# Note

I'm using [Git Flow](https://danielkummer.github.io/git-flow-cheatsheet/) to handle the repository operations.
