package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func main() {
	etcdEndpoint := os.Getenv("ETCD_ENDPOINT")
	if etcdEndpoint == "" {
		etcdEndpoint = "localhost:2379"
	}

	podName := os.Getenv("POD_NAME")
	if podName == "" {
		hostname, _ := os.Hostname()
		podName = hostname
	}

	fmt.Printf("[%s] Starting leader election app...\n", podName)
	fmt.Printf("[%s] Connecting to etcd at %s\n", podName, etcdEndpoint)

	// Retry connection to etcd
	var cli *clientv3.Client
	var err error
	for i := 0; i < 10; i++ {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   []string{etcdEndpoint},
			DialTimeout: 5 * time.Second,
		})
		if err == nil {
			break
		}
		fmt.Printf("[%s] Failed to connect to etcd, retrying... (%d/10)\n", podName, i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("[%s] Could not connect to etcd: %v\n", podName, err)
	}
	defer cli.Close()

	fmt.Printf("[%s] ✅ Connected to etcd successfully\n", podName)

	// Continuous leader election loop
	for {
		runLeaderElection(cli, podName)
		time.Sleep(2 * time.Second)
	}
}

func runLeaderElection(cli *clientv3.Client, nodeID string) {
	// SERVICE DISCOVERY (using Leases)
	// Create a lease that expires in 10 seconds
	session, err := concurrency.NewSession(cli, concurrency.WithTTL(10)) //Lease with ttl of 10 seconds
	if err != nil {
		log.Printf("[%s] Failed to create session: %v\n", nodeID, err)
		return
	}
	defer session.Close()

	election := concurrency.NewElection(session, "/k8s-demo/leader")
	ctx := context.Background()

	fmt.Printf("[%s] 🗳️  Campaigning for leadership...\n", nodeID)

	if err := election.Campaign(ctx, nodeID); err != nil { // other node will be blocked, only one can win
		log.Printf("[%s] Campaign failed: %v\n", nodeID, err)
		return
	}

	fmt.Printf("\n")
	fmt.Printf("╔════════════════════════════════════════╗\n")
	fmt.Printf("║  [%s] 👑 I AM THE LEADER! 👑  ║\n", nodeID)
	fmt.Printf("╚════════════════════════════════════════╝\n")
	fmt.Printf("\n")

	performLeaderDuties(ctx, nodeID, session)

	election.Resign(ctx) // -> resigns leadership
	fmt.Printf("[%s] 📤 Resigned from leadership\n", nodeID)
}

func performLeaderDuties(ctx context.Context, nodeID string, session *concurrency.Session) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	workTimer := time.After(30 * time.Second)
	workCount := 0

	for {
		select {
		case <-ticker.C:
			workCount++
			fmt.Printf("[%s] 💼 Performing leader work #%d (processing tasks, coordinating cluster...)\n", nodeID, workCount)

		case <-workTimer:
			fmt.Printf("[%s] ✅ Leadership term complete (30s), stepping down to allow others\n", nodeID)
			return

		case <-session.Done():
			fmt.Printf("[%s] ⚠️  Session expired! Lost leadership (network issue or etcd restart).\n", nodeID)
			return

		case <-ctx.Done():
			return
		}
	}
}
