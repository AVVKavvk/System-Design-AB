## Configuration Files

### docker-compose.yml

```yaml
name: read-replicas
services:
  mongo1:
    image: mongo:latest
    container_name: mongo1
    hostname: mongo1
    command: ["mongod", "--replSet", "myReplicaSet", "--bind_ip_all"]
    ports:
      - "127.0.0.1:27017:27017"
    volumes:
      - mongo1_data:/data/db

  mongo2:
    image: mongo:latest
    container_name: mongo2
    hostname: mongo2
    command: ["mongod", "--replSet", "myReplicaSet", "--bind_ip_all"]
    ports:
      - "127.0.0.2:27017:27017"
    volumes:
      - mongo2_data:/data/db

  mongo3:
    image: mongo:latest
    container_name: mongo3
    hostname: mongo3
    command: ["mongod", "--replSet", "myReplicaSet", "--bind_ip_all"]
    ports:
      - "127.0.0.3:27017:27017"
    volumes:
      - mongo3_data:/data/db

volumes:
  mongo1_data:
  mongo2_data:
  mongo3_data:
```

### /etc/hosts

```
127.0.0.1 mongo1
127.0.0.2 mongo2
127.0.0.3 mongo3
```

### MongoDB Replica Set Initialization

```javascript
rs.initiate({
  _id: "myReplicaSet",
  members: [
    { _id: 0, host: "mongo1:27017" },
    { _id: 1, host: "mongo2:27017" },
    { _id: 2, host: "mongo3:27017" },
  ],
});
```

### Go Connection String

```go
mongodb://mongo1:27017,mongo2:27017,mongo3:27017/?replicaSet=myReplicaSet
```

## Setup Instructions

1. Update `/etc/hosts` with the hostname mappings
2. Start Docker containers: `docker-compose up -d`
3. Initialize
