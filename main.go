package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

const timeoutInSecond = 10

func main() {
	bgCtx := context.Background()
	ctx, cancel := context.WithTimeout(bgCtx, timeoutInSecond*time.Second)
	defer cancel()

	timeoutWait := getTimeoutDetector(ctx)

	proctor := make(chan struct{})
	defer close(proctor)

	count := 0
	go func() {
		for {
			select {
			case <-proctor:
				count++
			}
		}
	}()

	go wordTest(proctor)

	<-timeoutWait

	fmt.Println(">>Num of correct answers=", count)
}

func getWords() []string {
	words := []string{
		"Desk",
		"Apple",
		"Monday",
		"June",
		"Sun",
		"Car",
	}
	return words
}

func getTimeoutDetector(ctx context.Context) <-chan bool {
	wait := make(chan bool)
	go func() {
		select {
		case <-ctx.Done():
			wait <- true
			close(wait)
			break
		}
	}()
	return wait
}

func wordTest(proctor chan<- struct{}) {
	words := getWords()
	s := bufio.NewScanner(os.Stdin)
	for _, w := range words {
		fmt.Println(">", w)
		for s.Scan() {
			input := s.Text()
			// Ignore case
			if strings.Compare(strings.ToLower(w), strings.ToLower(input)) == 0 {
				os.Stderr.WriteString(fmt.Sprintln("OK"))
				proctor <- struct{}{}
			} else {
				os.Stderr.WriteString(fmt.Sprintln("NG"))
			}
			break
		}
	}
}
