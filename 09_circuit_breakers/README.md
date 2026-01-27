# Circuit Breaker

In system design, a Circuit Breaker is a design pattern used to prevent a failure in one service from cascading to others. It acts as a safety net that "trips" (stops traffic) when a downstream service is struggling, giving it time to recover rather than overwhelming it with more requests.

It is inspired by the electrical circuit breakers in your home, which shut off the flow of electricity if they detect a dangerous surge.

## Why Do We Use It?

In a distributed system (like microservices), services call each other over a network. If a downstream service becomes slow or unresponsive:

1. `Resource Exhaustion`: The calling service will have many threads waiting for responses that never come, eventually running out of memory or CPU.

2. `Cascading Failures`: If Service A fails because Service B is slow, Service C (which calls A) might also fail. This can take down an entire system.

The circuit breaker prevents this by failing fast. Instead of waiting for a 30-second timeout, it returns an error immediately.

## How It Works (The 3 States)

A circuit breaker is essentially a state machine that wraps around a network call.

1. ### Closed (Normal Operation)

- `Behavior`: Requests flow through normally to the downstream service.

- `Monitoring`: The breaker tracks the number of failures.

- `Trigger`: If the failure rate stays below a certain threshold (e.g., < 5 failures in 10 seconds), it stays Closed. If the threshold is exceeded, it trips and moves to Open.

2. ### Open (Failure State)

- `Behavior`: The breaker stops calling the downstream service entirely. It immediately returns a "fallback" response or an error.

- `Purpose`: This gives the failing service a "cool-down" period to restart or recover without being bombarded by traffic.

- `Trigger`: After a predefined "reset timeout" (e.g., 30 seconds), it moves to the **Half-Open** state.

3. ### Half-Open (Testing Recovery)

- `Behavior`: The breaker allows a small number of "test" requests to pass through.

- `Purpose`: To check if the downstream service has actually recovered.

- `Trigger`: If the test requests succeed, the breaker assumes the service is healthy and moves back to **Closed**.
  - If the test requests fail, it assumes the service is still broken and moves back to Open.

## Circuit Breaker vs. Retry

These are often confused, but they serve opposite purposes:

- `Retry Pattern`: "The service might have had a tiny hiccup; let's try again immediately."

- `Circuit Breaker`: "The service is clearly struggling; stop trying so we don't kill it."

<br/>
<br/>
<br/>

# Circuit Breaker Implementation with Go

This repository demonstrates a Circuit Breaker pattern implementation in Go using the Echo framework.

## Overview

The Circuit Breaker pattern prevents an application from repeatedly attempting operations that are likely to fail, allowing it to detect failures and handle them gracefully.

## Running the Application

### Option 1: Using Air (Hot Reload)

```bash
air
```

### Option 2: Using Go Run

```bash
go run .
```

## Testing the Circuit Breaker

Once the application is running, open a new terminal and test the API endpoint:

```bash
curl http://localhost:8080/users
```

### Expected Behavior

- **First 3 requests**: The application makes actual calls to the external service (simulated to fail)
- **4th request onwards**: Circuit breaker opens and returns "Service Unavailable" error without attempting external calls

This demonstrates how the circuit breaker protects your application from cascading failures by stopping requests to a failing service after a threshold is reached.

!["IMG"](./images/image.png)

## How It Works

1. Circuit breaker monitors external service call failures
2. After consecutive failures exceed the threshold (3 in this case), the circuit "opens"
3. Subsequent requests fail fast with a "Service Unavailable" error
4. The circuit breaker prevents unnecessary load on the failing service

---

For more details on the Circuit Breaker pattern, refer to the implementation in the source code.

## Reference Docs

https://learn.microsoft.com/en-us/previous-versions/msp-n-p/dn589784(v=pandp.10)
https://github.com/sony/gobreaker
