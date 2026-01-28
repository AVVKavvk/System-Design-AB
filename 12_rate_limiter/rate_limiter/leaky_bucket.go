package rate_limiter

import (
	"fmt"
	"time"
)

type LeakyBucket struct {
	Capacity    int
	Bucket      chan struct{}
	ProcessRate time.Duration
}

func GetLeakyBucket(capacity int, processRate time.Duration) *LeakyBucket {

	lb := &LeakyBucket{
		Capacity:    capacity,
		Bucket:      make(chan struct{}, capacity),
		ProcessRate: processRate,
	}

	go removeFromBucketWithFixedRate(lb)

	return lb

}

func (lb *LeakyBucket) Allow() bool {
	select {
	case lb.Bucket <- struct{}{}: // Add one element to the bucket, if it's not full
		return true
	default:
		return false
	}
}

func removeFromBucketWithFixedRate(lb *LeakyBucket) {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Printf("Recovered from panic for removeFromBucketWithFixedRate: %v\n", r)
			go removeFromBucketWithFixedRate(lb)
		}
	}()

	ticker := time.NewTicker(lb.ProcessRate)

	for range ticker.C {
		select {
		case <-lb.Bucket: // Remove one element from the bucket
		default:
		}
	}
}
