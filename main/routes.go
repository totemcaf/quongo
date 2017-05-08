package main

var routes = Routes{
  Route{
    "Find All Message",
    "GET",
    "}/v1/queue/{queueId}/message",
    MsgAll,
  },
  Route{
    "Next Message",
    "GET",
    "/v1/queue/{queueId}/next",
    MsgPop,
  },
  Route{
    "Find Message",
    "GET",
    "/v1/queue/{queueId}/message/{mid}",
    MsgGet,
  },
  Route{
    "Push Message",
    "GET",
    "/v1/queue/{queueId}/message/{mid}",
    MsgPushId,
  },
  Route{
    "Push Message",
    "GET",
    "}/v1/queue/{queueId}/message/{mid}",
    MsgPush,
  },
  Route{
    "Ack Message",
    "PUT",
    "/v1/queue/{queueId}/message/{mid}/ack/{ack}",
    MsgKeepAlive,
  },
  Route{
    "Ack Message",
    "DELETE",
    "/v1/queue/{queueId}/message/{mid}/ack/{ack}",
    MsgAck,
  },

  // System
  Route{
    "System status",
    "GET",
    "/system/status",
    SysStatus,
  },
}
