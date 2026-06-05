package evidence

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type Status string

const (
	StatusPass Status = "pass"
	StatusFail Status = "fail"
	StatusSkip Status = "skip"

	StatusPassed  Status = StatusPass
	StatusFailed  Status = StatusFail
	StatusSkipped Status = StatusSkip
)

const SchemaVersion = "testkitx.evidence.v1"

type Run struct {
	Suite     string    `json:"suite"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
	Cases     []Case    `json:"cases"`
}

type Case struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Status   Status        `json:"status"`
	Duration time.Duration `json:"duration"`
	Message  string        `json:"message,omitempty"`
}

func (r Run) Validate() error {
	if r.Suite == "" {
		return fmt.Errorf("suite is required")
	}
	if r.StartedAt.IsZero() {
		return fmt.Errorf("started_at is required")
	}
	if r.EndedAt.IsZero() {
		return fmt.Errorf("ended_at is required")
	}
	if r.EndedAt.Before(r.StartedAt) {
		return fmt.Errorf("ended_at must not be before started_at")
	}
	if len(r.Cases) == 0 {
		return fmt.Errorf("cases must not be empty")
	}
	for i, c := range r.Cases {
		if c.ID == "" {
			return fmt.Errorf("cases[%d].id is required", i)
		}
		if c.Name == "" {
			return fmt.Errorf("cases[%d].name is required", i)
		}
		if !validStatus(c.Status) {
			return fmt.Errorf("cases[%d].status is invalid", i)
		}
		if c.Duration < 0 {
			return fmt.Errorf("cases[%d].duration must not be negative", i)
		}
	}
	return nil
}

type Report struct {
	SchemaVersion string    `json:"schema_version"`
	ID            string    `json:"id"`
	Status        Status    `json:"status"`
	GeneratedAt   time.Time `json:"generated_at"`
	Checks        []Check   `json:"checks"`
}

type Check struct {
	Name    string `json:"name"`
	Status  Status `json:"status"`
	Message string `json:"message,omitempty"`
}

func NewReport(id string, checks ...Check) Report {
	return Report{
		SchemaVersion: SchemaVersion,
		ID:            id,
		Status:        aggregateStatus(checks),
		GeneratedAt:   time.Now().UTC(),
		Checks:        append([]Check(nil), checks...),
	}
}

func Passed(name string) Check {
	return Check{Name: name, Status: StatusPassed}
}

func Failed(name, message string) Check {
	return Check{Name: name, Status: StatusFailed, Message: message}
}

func Skipped(name, message string) Check {
	return Check{Name: name, Status: StatusSkipped, Message: message}
}

func (r Report) Validate() error {
	if r.SchemaVersion == "" {
		return fmt.Errorf("schema_version is required")
	}
	if r.SchemaVersion != SchemaVersion {
		return fmt.Errorf("schema_version is invalid")
	}
	if r.ID == "" {
		return fmt.Errorf("id is required")
	}
	if r.Status == "" {
		return fmt.Errorf("status is required")
	}
	if !validStatus(r.Status) {
		return fmt.Errorf("status is invalid")
	}
	if r.GeneratedAt.IsZero() {
		return fmt.Errorf("generated_at is required")
	}
	if len(r.Checks) == 0 {
		return fmt.Errorf("checks must not be empty")
	}
	for i, check := range r.Checks {
		if check.Name == "" {
			return fmt.Errorf("checks[%d].name is required", i)
		}
		if !validStatus(check.Status) {
			return fmt.Errorf("checks[%d].status is invalid", i)
		}
	}
	if want := aggregateStatus(r.Checks); r.Status != want {
		return fmt.Errorf("status must match aggregate check status")
	}
	return nil
}

func Marshal(report Report) ([]byte, error) {
	if err := report.Validate(); err != nil {
		return nil, err
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

func Digest(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func aggregateStatus(checks []Check) Status {
	if len(checks) == 0 {
		return StatusSkipped
	}
	status := StatusPassed
	for _, check := range checks {
		if check.Status == StatusFailed {
			return StatusFailed
		}
		if check.Status == StatusSkipped {
			status = StatusSkipped
		}
	}
	return status
}

func validStatus(status Status) bool {
	switch status {
	case StatusPassed, StatusFailed, StatusSkipped:
		return true
	default:
		return false
	}
}
