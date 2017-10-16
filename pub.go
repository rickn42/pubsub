package pubsub

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"sync"
)

var NotComparable = fmt.Errorf("Arguments should be comparable type (for unsubscribe)!")

type event struct {
	issue interface{}
	value interface{}
}

type Subscriber interface {
	Receive(issue interface{}, value interface{})
}

type Publisher interface {
	Publish(issue interface{}, value interface{}) (err error)
	Subscribe(issue interface{}, sub Subscriber) error
	Unsubscribe(issue interface{}, sub Subscriber)
}

type publisher struct {
	q  chan event
	mu sync.RWMutex
	ss map[interface{}][]Subscriber
}

func NewPub() *publisher {
	p := &publisher{
		q:  make(chan event),
		ss: make(map[interface{}][]Subscriber),
	}

	go func() {
		var e event
		var ss []Subscriber

		for e = range p.q {
			p.mu.RLock()
			ss, _ = p.ss[e.issue]
			p.mu.RUnlock()
			for _, s := range ss {
				listenSafe(s, e)
			}
			runtime.Gosched()
		}
	}()

	return p
}

func listenSafe(sub Subscriber, e event) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("warn: Receive panic recovered\n", r)
		}
	}()

	sub.Receive(e.issue, e.value)
}

func (p *publisher) Publish(issue interface{}, value interface{}) (err error) {
	if !reflect.TypeOf(issue).Comparable() {
		return NotComparable
	}
	defer func() {
		if e, ok := recover().(error); ok {
			err = e
		}
	}()
	p.q <- event{issue, value}
	return nil
}

// Subscribe register Subscriber with issue.
// sub (Subscriber) should be comparable type for unsubscribe later.
func (p *publisher) Subscribe(issue interface{}, sub Subscriber) error {
	if !reflect.TypeOf(sub).Comparable() || !reflect.TypeOf(issue).Comparable() {
		return NotComparable
	}

	p.mu.Lock()
	p.ss[issue] = append(p.ss[issue], sub)
	p.mu.Unlock()

	return nil
}

func (p *publisher) Unsubscribe(issue interface{}, sub Subscriber) {
	if !reflect.TypeOf(issue).Comparable() {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if ss, ok := p.ss[issue]; ok {
		for i, s := range ss {
			if s == sub {
				p.ss[issue] = append(p.ss[issue][:i], p.ss[issue][i+1:]...)
				return
			}
		}
	}
}

func (p *publisher) Remove() {
	close(p.q)
}
