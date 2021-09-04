package runner

import (
	"time"
)

// Execution represents a sub-process of a session, being a separate "section" for output.
type Execution struct {
	// Details is mapped into output.text https://docs.github.com/en/rest/reference/checks#update-a-check-run
	Details string
	// ErrorDetails may hold a Golang error string, populated from runtime processing the request
	ErrorDetails string
	// Completed is a flag indicating whether the section has Completed (regardless of successfully or not)
	Completed bool
	// Successful is a flag representing the success of the section
	Successful bool
	// StartedAt is just useful meta to have
	StartedAt time.Time
	// SecondsRan represents how many seconds the section ran for until completion.
	SecondsRan int64
}

// AppendToDetails will append the input to the Details string, no formatting is done.
func (s *Execution) AppendToDetails(details string) {
	s.Details += details
}

func (s *Execution) AppendToErrorDetails(errorDetails string) {
	s.ErrorDetails += errorDetails
}

func (s *Execution) GetFinishedAt() time.Time {
	return time.Unix(s.StartedAt.Unix()+s.SecondsRan, 0)
}

func (s *Execution) GetSecondsRanFor() int64 {
	return s.SecondsRan
}

// MarkCompleted is a shortcut method to set multiple fields, should only be called once at completion.
func (s *Execution) MarkCompleted(success bool) {
	s.Completed = true
	s.Successful = success
	s.SecondsRan = s.StartedAt.Unix() - time.Now().Unix()
}
