# Rate Limiter

In system design, a **rate limiter** is a tool used to control the amount of incoming or outgoing traffic to a network or service. It acts as a defensive shield that limits the number of requests a user or client can make within a specific timeframe (e.g., "5 requests per second").

If a user exceeds the defined limit, the rate limiter blocks the additional requests, often returning an **HTTP 429 (Too Many Requests)** status code.

## Why Use a Rate Limiter?

- **Prevent DoS Attacks**: It stops malicious actors from overwhelming your servers with a flood of requests (Denial of Service).

- **Cost Control**: If you use third-party APIs that charge per request, a rate limiter ensures you don't exceed your budget.

- **Server Stability**: It prevents "noisy neighbors" (a single user taking up all resources) from slowing down the experience for everyone else.

- **Managing Traffic Spikes**: It helps maintain a consistent service level during unexpected surges in popularity.

## Common Algorithms

There are several ways to implement the logic behind a rate limiter. Here are the most popular ones:

1. **Token Bucket**

   A "bucket" holds a fixed number of tokens. Each request consumes one token. Tokens are added back to the bucket at a fixed rate. If the bucket is empty, the request is rejected.

   `Pro`: Allows for occasional bursts of traffic.

2. **Leaking Bucket**

   Similar to the token bucket, but requests are processed at a constant, fixed rate (like water leaking from a hole in a bucket).

   `Pro`: Ensures a very stable and smooth flow of requests to the backend.

3. **Fixed Window Counter**

   The timeline is divided into fixed blocks (e.g., 1-minute windows). Each window has a counter.

   `Con`: Can allow double the allowed traffic if many requests happen right at the edge of two windows (the "boundary problem").

4. **Sliding Window Log/Counter**

   A more advanced method that tracks the exact timestamp of each request or uses a rolling window to smooth out the "boundary problem" found in fixed windows.

## Where is it Implemented?

1. **Client-Side**: Unreliable, as users can bypass it by modifying the code.

2. **Server-Side**: Inside the application code (e.g., using a library in Go or Node.js).

3. **API Gateway/Middleware**: The most common approach. Tools like Nginx, Kong, or Amazon API Gateway handle the limiting before the request even reaches your application logic.

## Resources

[Design Rate Limiter](https://bytebytego.com/courses/system-design-interview/design-a-rate-limiter)
