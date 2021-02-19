package schedule

import (
	"errors"
	"sort"
	"time"
)

var (
	// ErrOverlappingOverrides is returns two override overlaps
	ErrOverlappingOverrides = errors.New("schedule: overlapping overrides")
	// ErrOverrideNotFound is returned if no override are found
	ErrOverrideNotFound = errors.New("schedule: override not found")
)

// Overrides represents an odered slice of overrides
type Overrides []*Override

// NewOverrides returns a new override slice
func NewOverrides() Overrides {
	return []*Override{}
}

func (o *Overrides) cleanup() {
	j := 0
	for i, ov := range *o {
		if now().After(ov.Stop) {
			j = i
			continue
		}

		break
	}

	if j == 0 {
		return
	}

	if (j + 1) == len(*o) {
		*o = []*Override{}
		return
	}

	*o = (*o)[j+1:]
}

// Add add a new override ensuring non overlapping overrides
func (o *Overrides) Add(override *Override) error {
	if override == nil {
		return nil
	}

	if err := override.validate(); err != nil {
		return err
	}

	for _, ov := range *o {
		if ov.Start.After(override.Stop) {
			break
		}

		if ov.Start.Equal(override.Start) || ov.Stop.Equal(override.Stop) {
			return ErrOverlappingOverrides
		}

		if ov.Start.Before(override.Start) && ov.Stop.After(override.Start) {
			return ErrOverlappingOverrides
		}

		if ov.Start.Before(override.Stop) && ov.Stop.After(override.Stop) {
			return ErrOverlappingOverrides
		}
	}

	override.generateID()

	*o = append(*o, override)
	sort.Sort(o)
	o.cleanup()
	return nil
}

// At returns the override at at given time or nil
func (o *Overrides) At(t time.Time) *Override {
	if o == nil || len(*o) == 0 {
		return nil
	}

	// Cleanup if the overrive is in the past
	o.cleanup()

	for _, ov := range *o {
		if ov.Start.Equal(t) || ov.Stop.Equal(t) {
			return ov
		}

		if ov.Start.Before(t) && ov.Stop.After(t) {
			return ov
		}
	}

	return nil
}

// Now returns the override now or nil
func (o *Overrides) Now() *Override {
	return o.At(now())
}

// Len implements the sort interface
func (o Overrides) Len() int {
	return len(o)
}

// Less implements the sort interface
func (o Overrides) Less(i, j int) bool {
	return o[i].Start.Before(o[j].Start)
}

// Swap implements the sort interface
func (o Overrides) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

// Delete deletes a timeslot
func (o *Overrides) Delete(id string) error {
	var found *int
	for i, override := range *o {
		if override.ID == id {
			found = &i
			break
		}
	}

	if found == nil {
		return ErrOverrideNotFound
	}

	if *found == (len(*o) - 1) {
		// This is the last element remove everything until the index
		*o = (*o)[:*found]
		return nil
	}

	// Remove the element at index _found_ and rebuild the slice
	*o = append((*o)[:*found], (*o)[*found+1:]...)
	return nil
}

// NextTime returns the next time
func (o *Overrides) NextTime() *time.Time {
	if len(*o) == 0 {
		return nil
	}

	o.cleanup()

	// Check the first one
	if (*o)[0].Start.After(now()) {
		return &(*o)[0].Start
	}

	// Check the second one
	if len(*o) > 1 {
		return &(*o)[1].Start
	}

	return nil
}
