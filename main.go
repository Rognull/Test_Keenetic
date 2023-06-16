package main

import (
	// "context"
	"flag"
	"log"
	"net/http"
	"strconv"

	// "sync"
	// "time"
	"encoding/json"
	"fmt"
)

type item struct {
	name string
	next *item
}

type queueItems struct {
	head *item
	tail *item
}

type queue struct {
	items map[string]*queueItems
}

var (
	port    = flag.Int("p", 8080, "server port")
	manager = queue{
		items: make(map[string]*queueItems),
	}
)

func (q queue) putInQueue(queueName string, itemName string) {
	j, ok := q.items[queueName]
	if !ok {
		buf := &item{name: itemName}
		q.items[queueName] = &queueItems{
			head: buf, tail: buf,
		}
	} else {
		j.tail.next = &item{
			name: itemName,
		}
		j.tail = j.tail.next
	}
}

func (q queue) getFromQueue(queueName string, timeout string) (string, error) {
	j, ok := q.items[queueName]
	res := ""
	if !ok {
		// q.items[queueName] = &queueItems{
		// 	head: &item{name: itemName}, tail: &item{name: itemName}, // make timeout
		fmt.Println(timeout)
	} else {
		if j.head == nil {
			fmt.Println(timeout)
		} else {
			res = j.head.name
			j.head = j.head.next
		}
	}
	return res, nil
}

func putRequest(w http.ResponseWriter, r *http.Request) {
	queueName := r.URL.Path
	name := r.URL.Query().Get("v")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		manager.putInQueue(queueName, name)
		WrapOK(w, "")
	}
}

func getRequest(w http.ResponseWriter, r *http.Request) {
	queueName := r.URL.Path
	time := r.URL.Query().Get("timeout")
	// if time == "" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// } else {
	res, err := manager.getFromQueue(queueName, time)
	if err == nil {
		fmt.Println(res)
		WrapOK(w, res)
	}
}

// }

func handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			putRequest(w, r)
		case http.MethodGet:
			getRequest(w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func WrapOK(w http.ResponseWriter, answer string) {
	res, _ := json.Marshal(answer)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(res))
}

func main() {
	flag.Parse()
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), handler()))
	// manager.putInQueue("animal", "a")
	// manager.putInQueue("animal", "b")
	// manager.putInQueue("animal", "c")
	// manager.getFromQueue("animal", "1")
	// manager.getFromQueue("animal", "1")
	// manager.getFromQueue("animal", "1")
}
