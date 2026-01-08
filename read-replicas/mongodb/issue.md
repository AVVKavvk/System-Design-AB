# MongoDB Replica Set - Split-Brain Networking Solution

## Problem Overview

Classic "Split-Brain" networking issue when running MongoDB Replica Sets inside Docker while the application runs on the host machine.

## The Technical Conflict

### 1. Internal Identity

- MongoDB nodes register themselves as: `mongo1:27017`, `mongo2:27017`, `mongo3:27017`
- These addresses work perfectly inside the Docker network
- MongoDB stores these specific strings in its internal configuration

### 2. Driver's Discovery Phase

- MongoDB Go Driver is a "Smart Client"
- Uses seed URIs to ask: "Who are all the members of this set?"
- Database responds with internal list: `mongo1:27017`, `mongo2:27017`, `mongo3:27017`
- Driver ignores initial localhost/mapped ports and connects to advertised names

### 3. Routing Mismatch

- **Hostname Issue**: Host machine doesn't recognize `mongo2`, `mongo3`
- **Port Issue**: Docker maps `mongo2` to `27018`, `mongo3` to `27019`, but database advertises `27017`
- Driver attempts `mongo2:27017` but nothing listens on that port on the host
- Result: EOF or context deadline exceeded errors

### 4. Why GET Worked but POST Failed

- **GET (Read)**: Used `SecondaryPreferred`, `mongo1` had 1:1 mapping (27017:27017) ✓
- **POST (Write)**: Writes go to Primary (`mongo2`), driver couldn't reach it on port 27017 ✗

## Problem vs Solution

| Problem             | Result                                           | Final Fix                                       |
| ------------------- | ------------------------------------------------ | ----------------------------------------------- |
| Hostname Resolution | Driver doesn't know what `mongo1` is             | Added entries to `/etc/hosts`                   |
| Port Mismatch       | Database advertises `:27017`, host uses `:27018` | Mapped host ports to unique IPs all on `:27017` |
| Connection Spam     | Logs filled with connection errors               | Implemented Singleton pattern with `sync.Once`  |

## Loopback IP Range Explanation

- `127.0.0.1` is the standard loopback address (localhost)
- Entire range `127.0.0.1` to `127.255.255.254` is reserved for loopback
- Allows multiple "virtual" local addresses on the same machine

### Why Multiple IPs (127.0.0.1, 127.0.0.2, 127.0.0.3)?

**The Conflict**: Three MongoDB containers all want port `27017`. On a single IP, only one program can own a port at a time.

**The Identity Issue**: MongoDB Replica Sets are strict about identity (Hostname + Port). If a node thinks it's `mongo2:27017`, it expects to be reached on port `27017`.

**The Solution**: Different local IPs give each container its own "private" port `27017`:

- `mongo1` → `127.0.0.1:27017`
- `mongo2` → `127.0.0.2:27017`
- `mongo3` → `127.0.0.3:27017`

## How It Fixed the 500 Error

1. **Discovery**: Go driver asks cluster for Primary → responds with `mongo2:27017`
2. **Resolution**: `/etc/hosts` tells driver that `mongo2` is at `127.0.0.2`
3. **Connection**: Driver connects to `127.0.0.2:27017` successfully (1:1 port mapping)
4. **Result**: POST (write) operations succeed ✓

## Current Working Setup

| Hostname | Host IP   | Host Port | Container Port |
| -------- | --------- | --------- | -------------- |
| mongo1   | 127.0.0.1 | 27017     | 27017          |
| mongo2   | 127.0.0.2 | 27017     | 27017          |
| mongo3   | 127.0.0.3 | 27017     | 27017          |

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

https://gemini.google.com/share/c7a5b47cbe73
