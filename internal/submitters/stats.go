package submitters

import "sync/atomic"

/* Atomic counters */

/**
 * Implementation of counter, could be changed the synchronization
 * method
 */
type FlagCounter struct {
	val int32
};

func (f *FlagCounter) Increment() {
	atomic.AddInt32(&f.val, 1)
}

func (f *FlagCounter) Get() int32 {
	return atomic.LoadInt32(&f.val);
}

/**
 * Statistics for flag submitter
 */
type Statistics struct {
	/* Number of failed Flags */
	flag_failed     FlagCounter
	/* Number of Flags Correctly submitted */
	flag_submitted  FlagCounter
	/* Number of duplicated Flags */
	flag_duplicated FlagCounter
};

func (s *Statistics) IncrementFailed() {
	s.flag_failed.Increment()
}

func (s *Statistics) GetFailed() (int32) {
	return s.flag_failed.Get()
}

func (s *Statistics) IncrementSubmitted() {
	s.flag_submitted.Increment()
}

func (s *Statistics) GetSubmitted() (int32) {
	return s.flag_submitted.Get()
}

func (s *Statistics) IncrementDuplicated() {
	s.flag_duplicated.Increment()
}

func (s *Statistics) GetDuplicated() (int32) {
	return s.flag_duplicated.Get()
}

var Stats Statistics
