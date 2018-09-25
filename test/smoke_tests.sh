#!/usr/bin/env bash

# Get status
curl http://localhost:7070/api/system/status

# Get stats
curl http://localhost:7070/api/system/stats

# Get existent queues
curl 'http://localhost:7070/api/v1/queue'

# Create Pepe queue
curl -XPOST -H "Content-type: application/json" -d '{
  "name": "pepe",
  "visibilityWindow": 10000
}' 'http://localhost:7070/api/v1/queue'

# Update time window duration
curl -XPUT -H "Content-type: application/json" -d '{
  "name": "pepe",
  "visibilityWindow": 20000000000
}' 'http://localhost:7070/api/v1/queue/pepe'

curl -XPOST -H "Content-type: application/json" -d '{
  "visibilityWindow": "10s"
}' 'http://localhost:7070/api/v1/queue/pepe/message'

# Get existent queue
curl 'http://localhost:7070/api/v1/queue/pepe'

# Create Invalid name queue
curl -XPOST -H "Content-type: application/json" -d '{
  "name": "invalid &% name",
  "visibilityWindow": 10000000
}' 'http://localhost:7070/api/v1/queue'

# Too long name
curl -XPOST -H "Content-type: application/json" -d '{
  "name": "name123456789012345678901234567890",
  "visibilityWindow": 10000000
}' 'http://localhost:7070/api/v1/queue'


# Add a message
curl -XPOST -H "Content-type: application/json" -d '{
  "anId": "aaaasasdoipI",
  "name": "Something cute",
  "aValue": 32432
}' 'http://localhost:7070/api/v1/queue/pepe/message'

# Add many messages
for i in {1..50}
do
    curl -XPOST -H "Content-type: application/json" -d "{
      \"anId\": \"aaaasasdoipI-${i}\",
      \"name\": \"Something cute\",
      \"aValue\": ${i}
    }" 'http://localhost:7070/api/v1/queue/pepe/message'
done

# Find next message
curl 'http://localhost:7070/api/v1/queue/pepe/message/next'


curl -XPOST -H "Content-type: application/json" -d '{
  "anId": "otra cosa",
  "name": "muy bonito",
  "aValue": 111
}' 'http://localhost:7070/api/v1/queue/pepe1/message'


curl -XPOST -H "Content-type: application/json" -d '{
  "name": "pepe1",
  "visibilityWindow": 10000000000
}' 'http://localhost:7070/api/v1/queue'
