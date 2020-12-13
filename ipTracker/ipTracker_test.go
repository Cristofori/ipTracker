package ipTracker

import (
	"fmt"
	"math/rand"
	"testing"
)

func Test_Empty(t *testing.T) {
	Clear()
	result := Top100()

	if len(result) != 0 {
		t.Error("Expected an empty list")
	}
}

func Test_One(t *testing.T) {
	Clear()

	const ip = "192.168.0.1"
	RequestHandled(ip)
	results := Top100()

	if len(results) != 1 {
		t.Errorf("Expected one entry in the list, got %v", len(results))
	}

	expected := fmt.Sprintf("%v (1)", ip)
	result := results[0]

	if result != expected {
		t.Errorf("Unexpected entry in results. Expected (%v), received (%v)", expected, result)
	}
}

func Test_OneIpMultiple(t *testing.T) {
	Clear()

	const ip = "8.8.8.8"
	for i := 0; i < 1000; i++ {
		RequestHandled(ip)
	}

	results := Top100()

	if len(results) != 1 {
		t.Errorf("Expected one entry in the list, got %v", len(results))
	}

	expected := fmt.Sprintf("%v (1000)", ip)
	result := results[0]

	if result != expected {
		t.Errorf("Unexpected entry in results. Expected (%v), received (%v)", expected, result)
	}
}

func Test_Tie(t *testing.T) {
	Clear()

	const ip1 = "10.0.0.1"
	const ip2 = "192.168.0.1"

	RequestHandled(ip1)
	RequestHandled(ip2)

	results := Top100()

	if len(results) != 2 {
		t.Errorf("Expected 2 entries, got %v", len(results))
	}

	// If both entries have the same count the first one to reach that count should be higher in the list
	expected := fmt.Sprintf("%v (1)", ip1)
	if results[0] != expected {
		t.Errorf("Expected %v to be the first entry in the list, found %v", expected, results[0])
	}

	expected = fmt.Sprintf("%v (1)", ip2)
	if results[1] != expected {
		t.Errorf("Expected %v to be the second entry in the list, found %v", expected, results[1])
	}
}

func Test_BasicSort(t *testing.T) {
	Clear()

	const ip1 = "1.0.0.1"
	const ip2 = "10.0.0.1"

	RequestHandled(ip2)
	RequestHandled(ip1)
	RequestHandled(ip2)

	results := Top100()

	if len(results) != 2 {
		t.Errorf("Expected 2 entries, got %v", len(results))
	}

	expected := fmt.Sprintf("%v (2)", ip2)
	if results[0] != expected {
		t.Errorf("Expected %v to be the first entry in the list, found %v", expected, results[0])
	}

	expected = fmt.Sprintf("%v (1)", ip1)
	if results[1] != expected {
		t.Errorf("Expected %v to be the second entry in the list, found %v", expected, results[1])
	}
}

func Test_SortBubbling(t *testing.T) {
	Clear()

	const ip1 = "1.2.3.4"
	const ip2 = "192.168.0.1"
	const ip3 = "9.9.9.9"

	RequestHandled(ip1)
	RequestHandled(ip1)
	RequestHandled(ip1)

	RequestHandled(ip2)
	RequestHandled(ip2)
	RequestHandled(ip2)

	RequestHandled(ip3)
	RequestHandled(ip3)
	RequestHandled(ip3)

	results := Top100()

	entry3 := results[2]
	if entry3 != fmt.Sprintf("%v (3)", ip3) {
		t.Errorf("Expected IP %v to be the 3rd entry, found %v", ip3, entry3)
	}

	RequestHandled(ip3)
	results = Top100()

	// It should have hopped over multiple entries in the list
	entry1 := results[0]
	if entry1 != fmt.Sprintf("%v (4)", ip3) {
		t.Errorf("Sorting failed, %v should be the first IP in the list. Found %v", ip3, entry1)
	}
}

func Test_TrivialSort(t *testing.T) {
	Clear()

	RequestHandled("1")
	RequestHandled("1")
	RequestHandled("1")
	RequestHandled("1")
	RequestHandled("1")

	RequestHandled("2")
	RequestHandled("2")
	RequestHandled("2")
	RequestHandled("2")

	RequestHandled("3")
	RequestHandled("3")
	RequestHandled("3")

	RequestHandled("4")
	RequestHandled("4")

	RequestHandled("5")

	results := Top100()

	if results[0] != "1 (5)" {
		t.Errorf("Expected 1 to be first")
	}

	if results[1] != "2 (4)" {
		t.Errorf("Expected 2 to be second")
	}

	if results[2] != "3 (3)" {
		t.Errorf("Expected 3 to be third")
	}

	if results[3] != "4 (2)" {
		t.Errorf("Expected 4 to be fourth")
	}

	if results[4] != "5 (1)" {
		t.Errorf("Expected 5 to be fifth")
	}
}

func Test_ReverseSort(t *testing.T) {
	Clear()

	RequestHandled("5")

	RequestHandled("4")
	RequestHandled("4")

	RequestHandled("3")
	RequestHandled("3")
	RequestHandled("3")

	RequestHandled("2")
	RequestHandled("2")
	RequestHandled("2")
	RequestHandled("2")

	RequestHandled("1")
	RequestHandled("1")
	RequestHandled("1")
	RequestHandled("1")
	RequestHandled("1")

	results := Top100()

	if results[0] != "1 (5)" {
		t.Errorf("Expected 1 to be first")
	}

	if results[1] != "2 (4)" {
		t.Errorf("Expected 2 to be second")
	}

	if results[2] != "3 (3)" {
		t.Errorf("Expected 3 to be third")
	}

	if results[3] != "4 (2)" {
		t.Errorf("Expected 4 to be fourth")
	}

	if results[4] != "5 (1)" {
		t.Errorf("Expected 5 to be fifth")
	}
}

func Test_SortMax(t *testing.T) {
	Clear()

	for i := 0; i < SortedSize; i++ {
		for j := 0; j <= i; j++ {
			RequestHandled(fmt.Sprintf("%v", i+1))
		}
	}

	results := Top100()
	if len(results) != SortedSize {
		t.Errorf("Results should have %v entries, got %v", SortedSize, len(results))
	}

	first := results[0]
	last := results[99]

	expectedFirst := fmt.Sprintf("%v (%v)", SortedSize, SortedSize)
	expectedLast := "1 (1)"

	if first != expectedFirst {
		t.Errorf("Expect first entry to be %v, got %v", expectedFirst, first)
	}

	if last != expectedLast {
		t.Errorf("Expect first entry to be %v, got %v", expectedLast, last)
	}
}

func Test_SortMaxPlusOne(t *testing.T) {
	Clear()

	for i := 0; i < 101; i++ {
		for j := 0; j <= i; j++ {
			RequestHandled(fmt.Sprintf("%v", i+1))
		}
	}

	results := Top100()
	if len(results) != SortedSize {
		t.Errorf("Results should have %v entries, got %v", SortedSize, len(results))
	}

	first := results[0]
	last := results[99]

	expectedFirst := fmt.Sprintf("%v (%v)", SortedSize+1, SortedSize+1)
	expectedLast := "2 (2)"

	if first != expectedFirst {
		t.Errorf("Expect first entry to be %v, got %v", expectedFirst, first)
	}

	if last != expectedLast {
		t.Errorf("Expect first entry to be %v, got %v", expectedLast, last)
	}
}

func Test_SortMaxTimesTwo(t *testing.T) {
	Clear()

	for i := 0; i < SortedSize*2; i++ {
		for j := 0; j <= i; j++ {
			RequestHandled(fmt.Sprintf("%v", i+1))
		}
	}

	results := Top100()
	if len(results) != SortedSize {
		t.Fatalf("Results should have %v entries, got %v", SortedSize, len(results))
	}

	first := results[0]
	last := results[99]

	expectedFirst := fmt.Sprintf("%v (%v)", SortedSize*2, SortedSize*2)
	expectedLast := "101 (101)"

	if first != expectedFirst {
		t.Errorf("Expect first entry to be %v, got %v", expectedFirst, first)
	}

	if last != expectedLast {
		t.Errorf("Expect first entry to be %v, got %v", expectedLast, last)
	}
}

// This is a bad test as far as proper unit testing goes since the random
// element could lead to intermittent failures. It does, however, demonstrate
// that the algorithm is capable of handling millions of entries within seconds.
// It's also good at finding problems I failed to account for in the other
// tests.
func Test_Stress(t *testing.T) {
	Clear()

	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	// 10M requests
	for i := 0; i < 10000000; i++ {
		// Across 1M "ip addresses"
		RequestHandled(fmt.Sprintf("%v", rand.Intn(1000000)))
	}

	entry := sortedEntries.Front()

	for entry != nil {
		next := entry.Next()

		if next != nil {
			thisCount := getCount(entry)
			nextCount := getCount(next)

			if getCount(next) < getCount(entry) {
				t.Errorf("Entries list (len: %v) is not sorted correctly: %v (%v) -> %v (%v)",
					sortedEntries.Len(), toIp(entry), thisCount, toIp(next), nextCount)
			}
		}

		entry = next
	}

	lowCount := getCount(sortedEntries.Front())

	for ip, count := range ipCounts {
		if sortedLookup[ip] == nil && count > lowCount {
			t.Errorf("Found a record that should have been in the top 100 %v (%v):", ip, count)
		}
	}
}

func Benchmark_HighLoad(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Clear()

		for j := 0; j < 100000; j++ {
			fakeIp := fmt.Sprintf("%v", rand.Intn(10000))
			RequestHandled(fakeIp)
			Top100()
		}
	}
}
