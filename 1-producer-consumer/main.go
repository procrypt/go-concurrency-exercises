//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer szenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//


/*
We can make both the producer and consumer concurrent
Producer will send data to the channel and
Consumer will process it on the other end of the channel.
 */
package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(wg *sync.WaitGroup, ch chan<- *Tweet, stream Stream){
	defer wg.Done()
	wg.Add(1)
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			close(ch)
			return
		}
		ch <- tweet
	}
}

func consumer(wg *sync.WaitGroup, tweets <-chan *Tweet) {
	wg.Add(1)
	defer wg.Done()
	for t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()
	tweet := make(chan *Tweet)
	var wg sync.WaitGroup

	// Producer
	go producer(&wg, tweet, stream)

	// Consumer

	select {
	// If these is anything on the channel
	// send it to the consumer
	case <-tweet:
		go consumer(&wg, tweet)

	// time.After() returns a recevive only channel after the specified duration
	// if there is anything on the channel returned by time.After() close
	// the tweet channel and return
	case <-time.After(time.Second*1):
		close(tweet)
		return
	}
	wg.Wait()
	fmt.Printf("Process took %s\n", time.Since(start))
}
