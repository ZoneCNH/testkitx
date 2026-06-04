// Package obstest records fake observability events without provider SDKs.
package obstest

import "sync"

type Recorder struct {
	mu       sync.Mutex
	Counters map[string]int
	Logs     []string
}

func NewRecorder() *Recorder        { return &Recorder{Counters: map[string]int{}} }
func (r *Recorder) Inc(name string) { r.mu.Lock(); defer r.mu.Unlock(); r.Counters[name]++ }
func (r *Recorder) Log(message string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Logs = append(r.Logs, message)
}
func (r *Recorder) Count(name string) int { r.mu.Lock(); defer r.mu.Unlock(); return r.Counters[name] }
