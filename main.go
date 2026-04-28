package main

import (
	"sync"
	"time"
)

type queue struct {
	msgs []string
	wait []chan string
}

type broker struct {
	mu sync.Mutex
	qs map[string]*queue
}

func (b *broker) q(name string) *queue {
	q, ok := b.qs[name]
	if !ok {
		q = &queue{}
		b.qs[name] = q
	}
	return q
}

func (b *broker) put(name, msg string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	q := b.q(name)
	if len(q.wait) > 0 {
		ch := q.wait[0]
		q.wait = q.wait[1:]
		ch <- msg
		return
	}
	q.msgs = append(q.msgs, msg)
}

func (b *broker) take(name string, timeout time.Duration) (string, bool) {
	b.mu.Lock()
	q := b.q(name)
	if len(q.msgs) > 0 {
		m := q.msgs[0]
		q.msgs = q.msgs[1:]
		b.mu.Unlock()
		return m, true
	}
	if timeout == 0 {
		b.mu.Unlock()
		return "", false
	}
	ch := make(chan string, 1)
	q.wait = append(q.wait, ch)
	b.mu.Unlock()

	select {
	case m := <-ch:
		return m, true
	case <-time.After(timeout):
		b.mu.Lock()
		for i, w := range q.wait {
			if w == ch {
				q.wait = append(q.wait[:i], q.wait[i+1:]...)
				b.mu.Unlock()
				return "", false
			}
		}
		b.mu.Unlock()
		return <-ch, true
	}
}

func main() {}
