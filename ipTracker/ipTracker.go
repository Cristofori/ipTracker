package ipTracker

import (
	"container/list"
	"fmt"
)

const SortedSize = 100

// Maps an IP address to a count of how many times we've handled it
var ipCounts map[string]uint64

// Maps an IP address to its corresponding node in the sorted linked list
var sortedLookup map[string]*list.Element

// Linked list of the most frequently handled IP addresses, sorted in place
var sortedEntries list.List

func init() {
	Clear()
}

// A couple of convenience functions to avoid having to type out some ugly code
// more than once
func toIp(el *list.Element) string {
	return el.Value.(string)
}

func getCount(el *list.Element) uint64 {
	return ipCounts[toIp(el)]
}

func RequestHandled(ip string) {
	// We're going to assume these are valid IP addresses
	ipCounts[ip]++
	count := ipCounts[ip]

	if count > 1 {
		// We've seen this IP address before
		frontElement := sortedEntries.Front()
		var elementToSort *list.Element

		if sortedLookup[ip] != nil {
			// It's already in our sorted list
			elementToSort = sortedLookup[ip]
		} else if getCount(frontElement) < count {
			// It's not in the sorted list but it just beat out the lowest element
			elementToSort = sortedEntries.InsertAfter(ip, frontElement)
			sortedLookup[ip] = elementToSort

			if sortedEntries.Len() > SortedSize {
				delete(sortedLookup, toIp(frontElement))
				sortedEntries.Remove(frontElement)
			}
		}

		if elementToSort != nil {
			// This request affected one of the top IP addresses, so update its
			// position in the sorted list
			var newPosition *list.Element
			nextElement := elementToSort.Next()

			// Finds its place in line
			for nextElement != nil && getCount(nextElement) < count {
				newPosition = nextElement
				nextElement = nextElement.Next()
			}

			if newPosition != nil {
				sortedEntries.MoveAfter(elementToSort, newPosition)
			}
		}
	} else if sortedEntries.Len() < SortedSize {
		// Add new IPs encountered to the sorted list until it's full
		element := sortedEntries.PushFront(ip)
		sortedLookup[ip] = element
	}
}

func Top100() []string {
	result := make([]string, sortedEntries.Len())

	current := sortedEntries.Back()

	if current == nil {
		return result
	}

	index := 0
	for current != nil {
		ip := toIp(current)

		// The requirements only specify to list the most frequent IP addresses, but
		// I wanted to see the count as well so I put it in there for good measure.
		result[index] = fmt.Sprintf("%v (%v)", ip, ipCounts[ip])

		current = current.Prev()
		index++
	}

	return result
}

func Clear() {
	sortedEntries.Init()
	ipCounts = map[string]uint64{}
	sortedLookup = map[string]*list.Element{}
}
