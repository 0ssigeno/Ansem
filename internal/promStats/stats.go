package promStats

import "sync/atomic"

/* Atomic counters */

/**
 * Implementation of counter, could be changed the synchronization
 * method
 */
type FlagCounter struct {
	val int32
}

func (f *FlagCounter) Increment() {
	atomic.AddInt32(&f.val, 1)
}

func (f *FlagCounter) Get() int32 {
	return atomic.LoadInt32(&f.val)
}

/**
 * Statistics for flag submitter
 */
type Statistics struct {
	/* Number of failed Flags */
	flagFailed FlagCounter
	/* Number of Flags Correctly submitted */
	flagSubmitted FlagCounter
	/* Number of duplicated Flags */
	flagDuplicated FlagCounter
}

func (s *Statistics) IncrementFailed() {
	s.flagFailed.Increment()
}

func (s *Statistics) GetFailed() int32 {
	return s.flagFailed.Get()
}

func (s *Statistics) IncrementSubmitted() {
	s.flagSubmitted.Increment()
}

func (s *Statistics) GetSubmitted() int32 {
	return s.flagSubmitted.Get()
}

func (s *Statistics) IncrementDuplicated() {
	s.flagDuplicated.Increment()
}

func (s *Statistics) GetDuplicated() int32 {
	return s.flagDuplicated.Get()
}

var Stats Statistics
