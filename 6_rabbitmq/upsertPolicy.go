package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// PolicyDefinition represents the "definition" body of a RabbitMQ policy
type PolicyDefinition map[string]interface{}

// UpsertPolicy handles both creation and updates of RabbitMQ policies
func UpsertPolicy(name string, pattern string, definition PolicyDefinition) error {
	const (
		user  = "guest"
		pass  = "guest"
		host  = "localhost:15672"
		vhost = "%2f" // Default vhost "/" URL encoded
	)

	// Construct the payload
	payload := map[string]interface{}{
		"pattern":    pattern,
		"definition": definition,
		"apply-to":   "queues",
		"priority":   0,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/api/policies/%s/%s", host, vhost, name)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.SetBasicAuth(user, pass)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	fmt.Printf(" [✓] Policy '%s' synchronized successfully\n", name)
	return nil
}
