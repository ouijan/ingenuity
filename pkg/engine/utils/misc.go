package utils

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

const ServerTickRate = 160 * time.Millisecond

/**
 * Assert is a simple assertion function that panics if the condition is false.
 */
func Assert(condition bool, message ...string) {
	if !condition {
		panic(fmt.Sprintf("Assert Failed! %s", message))
	}
}

/**
 * Cast is a type assertion function that attempts to cast a value to the specified type T.
 */
func Cast[T any](value any) (T, bool) {
	castedVal, ok := value.(T)
	return castedVal, ok
}

/*
*
  - ReadStdIn reads from standard input line by line and calls the provided yield function for each line.
*/
func ReadStdIn(yield func(string) bool) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if !yield(line) {
			return
		}
	}
	// if err := scanner.Err(); err != nil {
	// 	log.Error("Error reading from stdin: %v", err)
	// }
}

/**
 * ChanSelect reads a message from a channel but does NOT block if empty
 */
func ChanSelect[T any](ch chan T) (T, bool) {
	select {
	case msg, ok := <-ch:
		if !ok {
			return msg, false
		}
		return msg, true
	default:
		return *new(T), false
	}
}
