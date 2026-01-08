127.0.0.1 mongo1 mongo2 mongo3

rs.initiate({
\_id: "myReplicaSet",
members: [
{ _id: 0, host: "localhost:27017" },
{ _id: 1, host: "mongo2:27017" },
{ _id: 2, host: "mongo3:27017" }
]
});
