package nagios

import "time"

type MergeRequest struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string
}
