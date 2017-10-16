package pubsub_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/kohalamakai/pubsub"
)

func Example() {
	pub := pubsub.NewPub()

	var s1 pubsub.Subscriber
	s1 = pubsub.NewSubFunc(func(issue, v interface{}) {
		fmt.Println("s1", issue, v)
	})

	var s2 pubsub.Subscriber
	s2 = pubsub.NewSubFunc(func(issue, v interface{}) {
		fmt.Println("s2", issue, v)
	})

	var issue = "hello"
	var err error

	err = pub.Subscribe(issue, s1)
	if err != nil {
		fmt.Println("subscribe failed:", err)
		return
	}
	err = pub.Subscribe(issue, s2)
	if err != nil {
		fmt.Println("subscribe failed:", err)
		return
	}
	err = pub.Publish(issue, 100)
	if err != nil {
		fmt.Println("publish failed:", err)
		return
	}

	pub.Unsubscribe(issue, s1)
	fmt.Println("unsubscribe hello s1")
	time.Sleep(time.Millisecond)

	pub.Publish(issue, 200)
	time.Sleep(time.Millisecond)

	pub.Subscribe(issue, s1)
	fmt.Println("resubscribe hello s1")
	time.Sleep(time.Millisecond)

	pub.Publish(issue, 300)
	time.Sleep(time.Millisecond)

	// Output:
	// s1 hello 100
	// s2 hello 100
	// unsubscribe hello s1
	// s2 hello 200
	// resubscribe hello s1
	// s2 hello 300
	// s1 hello 300
}

func TestPublisher_Subscribe(t *testing.T) {

	result := make(chan struct{}, math.MaxInt64)

	s1 := pubsub.NewSubFunc(func(issue, value interface{}) {
		result <- struct{}{}
	})

	s2 := pubsub.NewSubFunc(func(issue, value interface{}) {
		result <- struct{}{}
	})

	pub := pubsub.NewPub()
	pub.Subscribe("hello", s1)
	pub.Subscribe("hello", s2)

	n := 1000000
	start := time.Now()
	go func() {
		for i := 0; i < n; i++ {
			pub.Publish("hello", struct{}{})
		}
	}()

	var sum int
	tick := time.NewTimer(3 * time.Second)
	for {
		select {
		case <-result:
			sum++
		case <-tick.C:
			t.Error("some evt missing", sum)
			return
		}
		if 2*n == sum {
			fmt.Println("publish", n, "count", time.Now().Sub(start), "time")
			return
		}
	}
}
