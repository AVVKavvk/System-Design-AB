# Consistent Hashing

#### In distributed systems, consistent hashing is a method used to determine which node (server) "owns" or is responsible for a specific piece of data.

#### Unlike traditional hashing, where a change in the number of servers forces almost every key to be remapped, consistent hashing ensures that only a small fraction of keys are moved when nodes join or leave.

## How Ownership is Established

The system uses a conceptual Hash Ring, which represents the entire range of possible hash values (for example, from 3$0$ to 4$2^{32}-1$).5

1. `Mapping Nodes`: Every server is hashed (using its IP or ID) and placed at a specific point on the ring.

2. `Mapping Data`: Every piece of data (key) is hashed using the same function and placed on the same ring.

3. `Determining Ownership`: To find the "owner," you start at the data’s hash position and move clockwise around the ring. The first server you encounter is the one that "owns" that data.

## Key Characteristics of Ownership

- `Sticky Ownership`: Because the hash of a key doesn't change, it will always map to the same server as long as the cluster remains stable.

- `Minimal Disruption`: If a server (Node B) is removed, only the data it owned is transferred to its immediate clockwise neighbor (Node C). The rest of the cluster's ownership remains untouched.

- `Virtual Nodes`: To prevent one server from owning a much larger "slice" of the ring than others, systems often use Virtual Nodes. Each physical server is hashed multiple times to appear at several different locations on the ring, ensuring a more even distribution of ownership.

## Common Use Cases

- `Distributed Caching (e.g., Memcached)`: Ensures that when a cache server goes down, only the keys on that server are lost, preventing a "cache stampede" on the database.

- `Distributed Databases (e.g., Amazon DynamoDB, Apache Cassandra)`: Used to partition data across hundreds of nodes while maintaining high availability.

- `Load Balancing`: Directing specific user requests to the same backend server to maintain session state.
