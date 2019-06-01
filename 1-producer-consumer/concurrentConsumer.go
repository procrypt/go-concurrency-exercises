//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer szenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(stream Stream) (tweets <-chan *Tweet) {
	ch := make(chan *Tweet,len(stream.tweets))
	defer close(ch)
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			return ch
		}
		ch <- tweet
	}

}

func consumer(wg *sync.WaitGroup, tweets <-chan *Tweet) {
	for t := range tweets {
		wg.Add(1)
		go func(s *Tweet) {
			defer wg.Done()
			if s.IsTalkingAboutGo() {
				fmt.Println(s.Username, "\ttweets about golang")
			} else {
				fmt.Println(s.Username, "\tdoes not tweet about golang")
			}
		}(t)
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()

	// Producer
	tweets := producer(stream)

	// Consumer
	var wg sync.WaitGroup
	consumer(&wg, tweets)
	wg.Wait()
	fmt.Printf("Process took %s\n", time.Since(start))
}
