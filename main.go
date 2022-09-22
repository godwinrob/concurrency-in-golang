package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	rand "math/rand"
	"sync"
	"time"
)

type Response struct {
	Result       string        `db:"result" json:"result"`
	Runtime      time.Duration `db:"runtime" json:"runtime"`
	ErrorMessage error         `db:"error" json:"error,omitempty"`
}

type ResponseResult struct {
	Response Response `json:"response"`
	Error    error    `json:"error"`
}

func main() {

	// Simulate a call to a function that calls multiple databases
	result1, err := multipleDatabaseCalls()
	if err != nil {
		err := fmt.Errorf("multipleDatabaseCalls returned an error: %v", result1.ErrorMessage.Error())
		log.Println(err)
	} else {
		log.Println("multipleDatabaseCalls success! result: " + result1.Result)
		log.Printf("time to complete: %s", result1.Runtime)
	}

	// Simulate a call to a function that will get killed after 10 seconds runtime
	result2, err := functionWithHardTimeLimit()
	if err != nil {
		err := fmt.Errorf("functionWithHardTimeLimit returned an error: %v", result2.ErrorMessage.Error())
		log.Println(err)
	} else {
		log.Println("multipleDatabaseCalls success! result: " + result2.Result)
		log.Printf("time to complete: %s", result2.Runtime)
	}
}

// This function will make a number of database calls one after another
func multipleDatabaseCalls() (Response, error) {
	log.Println("Starting multipleDatabaseCalls")
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(3)

	results := ""

	go func() {
		// "DB call 1"
		log.Println("calling db 1")
		time.Sleep(8 * time.Second)
		log.Println("Result set 1 returned")
		results += "'db 1 result set' "
		wg.Done()
	}()

	go func() {
		// "DB call 2"
		log.Println("calling db 2")
		time.Sleep(4 * time.Second)
		log.Println("Result set 2 returned")
		results += "'db 2 result set' "
		wg.Done()
	}()

	go func() {
		// "DB call 3"
		log.Println("calling db 3")
		time.Sleep(9 * time.Second)
		log.Println("Result set 3 returned")
		results += "'db 3 result set' "
		wg.Done()
	}()

	wg.Wait()

	return Response{
		Result:  results,
		Runtime: time.Since(start),
	}, nil
}

// This simulates the signature of functionWithHardTimeLimit to return before lambda timeout
// kills the function
func functionWithHardTimeLimit() (Response, error) {
	result := make(chan ResponseResult, 1)

	go func() {
		result <- doFunctionWithHardTimeLimit()
	}()

	select {
	case <-time.After(9 * time.Second):
		return Response{
			ErrorMessage: errors.New("function reached timeout limit and returned"),
		}, errors.New("function timed out")
	case result := <-result:
		return Response{
			Result:  result.Response.Result,
			Runtime: result.Response.Runtime,
		}, nil
	}
}

// This function will fail randomly when database response takes too long
func doFunctionWithHardTimeLimit() ResponseResult {
	log.Println("Starting functionWithHardTimeLimit")

	// Time in seconds that Lambda will be killed after
	maxDuration := 10

	// Randomize run time to randomly exceed max allowed run time
	// runTime can randomly take between 1 and 19 seconds
	// ok will return false if Runtime is longer than Max Duration
	runTime, ok := randomTime(maxDuration)

	// Simulated failed call
	if !ok {
		time.Sleep(10 * time.Second)
		// Return is meant to simulate a non-existent response
		// We would expect an error to return here for a timed-out response
		// but instead our application will panic to simulate a lambda being killed
		panic("Function took too long and timed out. Response was not returned")
	}

	// Simulated successful response
	time.Sleep(time.Duration(runTime) * time.Second)
	return ResponseResult{
		Response{
			Runtime: time.Duration(runTime) * time.Second,
			Result:  "Database Result returned successfully",
		},
		nil,
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
