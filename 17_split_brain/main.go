package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	myID   = os.Getenv("NODE_ID")
	peers  = strings.Split(os.Getenv("PEERS"), ",")
	status = "Follower"
	client = &http.Client{
		Timeout: 1 * time.Second, // 1 second timeout
	}
)

func main() {
	// Background check: Ping peers to see who is alive
	go func() {
		for {
			aliveCount := 0
			for _, peer := range peers {
				if peer == "" {
					continue
				}
				resp, err := client.Get(fmt.Sprintf("http://%s:8080/health", peer))
				if err == nil && resp.StatusCode == 200 {
					aliveCount++
					resp.Body.Close()
				}
			}

			// NAIVE LOGIC: If I can't see others, I must be the leader!
			// (A real system would require aliveCount > 2 for a 5-node cluster)
			if aliveCount == 0 {
				status = "LEADER (Split-Brain)"
			} else {
				status = "Follower"
			}

			fmt.Printf("Node %s: I see %d peers. Status: %s\n", myID, aliveCount, status)
			time.Sleep(2 * time.Second)
		}
	}()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	fmt.Printf("Node %s starting...\n", myID)
	http.ListenAndServe(":8080", nil)
}
