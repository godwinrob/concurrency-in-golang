package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"log"
	rand "math/rand"
	"time"
)

type Response struct {
	Result  string        `db:"result" json:"result"`
	Runtime time.Duration `db:"runtime" json:"runtime"`
}

func main() {

	// Simulate a call to a function that calls multiple databases
	result1 := multipleDatabaseCalls()
	log.Println(result1)

	// Simulate a call to a function that will get killed after 10 seconds runtime
	result2 := functionWithHardTimeLimit()
	log.Println(result2)
}

// This function will make a number of database calls one after another
func multipleDatabaseCalls() Response {
	log.Println("Starting multipleDatabaseCalls")
	start := time.Now()

	results := ""

	// "DB call 1"
	time.Sleep(8 * time.Second)
	log.Println("Result set 1 returned")
	results += "'db call 1 result set' "
	// "DB call 2"
	time.Sleep(4 * time.Second)
	log.Println("Result set 2 returned")
	results += "'db call 2 result set' "
	// "DB call 3"
	time.Sleep(9 * time.Second)
	log.Println("Result set 3 returned")
	results += "'db call 3 result set' "

	return Response{
		Result:  results,
		Runtime: time.Since(start),
	}
}

// This function will fail randomly when database response takes too long
func functionWithHardTimeLimit() Response {
	log.Println("Starting functionWithHardTimeLimit")

	// Time in seconds that Lambda will be killed after
	maxDuration := 10

	// Randomize if run time will exceed max allowed run time
	// runTime can randomly take between 1 and 19 seconds
	// ok will return false if Runtime is longer than Max Duration
	runTime, ok := randomTime(maxDuration)

	if !ok {
		// Simulated failed call
		time.Sleep(10 * time.Second)
		log.Println("Function took too long and timed out. No response")
		return Response{}
	} else {
		// Simulated successful response
		time.Sleep(time.Duration(runTime) * time.Second)
		return Response{
			Runtime: time.Duration(runTime) * time.Second,
			Result:  "Database Result returned successfully",
		}
	}

}

func randomTime(maxDuration int) (int, bool) {
	var src cryptoSource
	rnd := rand.New(src)
	runTime := rnd.Intn(20)

	if runTime > maxDuration {
		return runTime, false
	}

	return runTime, true
}

type cryptoSource struct{}

func (s cryptoSource) Seed(seed int64) {}

func (s cryptoSource) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func (s cryptoSource) Uint64() (v uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		log.Fatal(err)
	}
	return v
}
