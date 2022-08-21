// roadmap2 represents a simplified set of data structures to represent a roadmap.
package roadmap2

import "time"

// Item represents a "roadmap item" or box on a roadmap slide.
type Item struct {
	Name        string
	Description string
	ReleaseTime time.Time
	StreamName  string
}

type Items []Item
