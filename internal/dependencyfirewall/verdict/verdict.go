package verdict

import "fmt"

type Verdict string

const (
	Blocked Verdict = "blocked"
	Warning Verdict = "warning"
)

type Entry struct {
	Package   string  `json:"package"`
	Version   string  `json:"version,omitempty"`
	Verdict   Verdict `json:"verdict"`
	Reason    string  `json:"reason,omitempty"`
	Status    int     `json:"status,omitempty"`
	Timestamp string  `json:"timestamp,omitempty"`
}

func (e Entry) Key() string {
	return fmt.Sprintf("%s@%s:%s", e.Package, e.Version, e.Verdict)
}
