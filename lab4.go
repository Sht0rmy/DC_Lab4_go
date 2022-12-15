package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type city struct {
	name string
	id   int
}

type ticket struct {
	origin      city
	destination city
	price       int
}

type travelAgencies struct {
	tickets []ticket
	cities  []city
}

func addTicket(ta *travelAgencies, lock *sync.RWMutex, t ticket) {
	lock.Lock()
	ta.tickets = append(ta.tickets, t)
	lock.Unlock()
}

func addCity(ta *travelAgencies, lock *sync.RWMutex, c city) {
	lock.Lock()
	ta.cities = append(ta.cities, c)
	lock.Unlock()
}

func removeCity(ta *travelAgencies, lock *sync.RWMutex, c city) {
	lock.Lock()
	for i, city := range ta.cities {
		if city == c {
			ta.cities = append(ta.cities[:i], ta.cities[i+1:]...)
		}
	}
	for i, ticket := range ta.tickets {
		if ticket.origin == c || ticket.destination == c {
			ta.tickets = append(ta.tickets[:i], ta.tickets[i+1:]...)
		}
	}
	for i, city := range ta.cities {
		city.id = i
	}
	lock.Unlock()
}

func removeTicket(ta *travelAgencies, lock *sync.RWMutex, t ticket) {
	lock.Lock()
	for i, ticket := range ta.tickets {
		if ticket == t {
			ta.tickets = append(ta.tickets[:i], ta.tickets[i+1:]...)
		}
	}
	lock.Unlock()
}

func findPathPriceDijkstra(ta *travelAgencies, lock *sync.RWMutex, origin city, destination city) int {
	lock.RLock()
	var visited []city
	var unvisited []city
	var distance []int
	var previous []city
	for _, c := range ta.cities {
		if c.id == origin.id {
			distance = append(distance, 0)
		} else {
			distance = append(distance, 999999)
		}
		previous = append(previous, city{})
		unvisited = append(unvisited, c)
	}
	for len(unvisited) > 0 {
		var current city
		var currentDistance int
		for _, c := range unvisited {
			if distance[c.id] < currentDistance {
				current = c
				currentDistance = distance[c.id]
			}
		}
		for _, t := range ta.tickets {
			if t.origin == current {
				var alt int = distance[current.id] + t.price
				if alt < distance[t.destination.id] {
					distance[t.destination.id] = alt
					previous[t.destination.id] = current
				}
			}
		}
		visited = append(visited, current)
		for i, c := range unvisited {
			if c == current {
				unvisited = append(unvisited[:i], unvisited[i+1:]...)
			}
		}
	}
	lock.RUnlock()
	return distance[destination.id]
}

func citiesGenerator(ta *travelAgencies, lock *sync.RWMutex) {
	for i := 0; i < 10; i++ {
		c := city{"City" + string(i), i}
		addCity(ta, lock, c)
		fmt.Println("Added city")
		time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
	}
}

func ticketsGenerator(ta *travelAgencies, lock *sync.RWMutex) {
	time.Sleep(time.Duration(4000 * time.Millisecond))
	for i := 0; i < 100; i++ {
		t := ticket{ta.cities[rand.Intn(len(ta.cities))], ta.cities[rand.Intn(len(ta.cities))], rand.Intn(1000)}
		addTicket(ta, lock, t)
		fmt.Println("Added ticket")
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}
}

func removeCities(ta *travelAgencies, lock *sync.RWMutex) {
	time.Sleep(time.Duration(10000 * time.Millisecond))
	for i := 0; i < 5; i++ {
		removeCity(ta, lock, ta.cities[rand.Intn(len(ta.cities))])
		fmt.Println("Removed city")
		time.Sleep(time.Duration(rand.Intn(10000)) * time.Millisecond)

	}
}

func removeTickets(ta *travelAgencies, lock *sync.RWMutex) {
	time.Sleep(time.Duration(10000 * time.Millisecond))
	for i := 0; i < 5; i++ {
		removeTicket(ta, lock, ta.tickets[rand.Intn(len(ta.tickets))])
		fmt.Println("Removed ticket")
		time.Sleep(time.Duration(rand.Intn(70000)) * time.Millisecond)

	}
}

func findRandomPath(ta *travelAgencies, lock *sync.RWMutex) {
	time.Sleep(time.Duration(10000 * time.Millisecond))
	origin := ta.cities[rand.Intn(len(ta.cities))-1]
	destination := ta.cities[rand.Intn(len(ta.cities))-1]
	fmt.Println("Searching path from", origin.id, "to", destination.id)
	fmt.Println("Path from", origin.name, "to", destination.name, "is", findPathPriceDijkstra(ta, lock, origin, destination))
}

func main() {
	var rwMutex sync.RWMutex
	var ta travelAgencies

	var cities []city
	var tickets []ticket

	ta.cities = cities
	ta.tickets = tickets

	go citiesGenerator(&ta, &rwMutex)
	go ticketsGenerator(&ta, &rwMutex)
	go removeCities(&ta, &rwMutex)
	go removeTickets(&ta, &rwMutex)
	go findRandomPath(&ta, &rwMutex)

	for {
		time.Sleep(5 * time.Second)
	}
}
