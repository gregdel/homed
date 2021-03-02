package schedule

import (
	"reflect"
	"testing"
	"time"
)

func TestOverridesAdd(t *testing.T) {
	overrides := NewOverrides()

	n := now()

	tt := []struct {
		name     string
		override *Override
		expected error
	}{
		{
			name: "no override",
		},
		{
			name:     "invalid time",
			override: &Override{},
			expected: ErrInvalidTime,
		},
		{
			name: "first override",
			override: &Override{
				Start: n,
				Stop:  n.Add(1 * time.Hour),
			},
		},
		{
			name: "same time override",
			override: &Override{
				Start: n,
				Stop:  n.Add(1 * time.Hour),
			},
			expected: ErrOverlappingOverrides,
		},
		{
			name: "new timeslot starting at the end of the previous one",
			override: &Override{
				Start: n.Add(1 * time.Hour),
				Stop:  n.Add(2 * time.Hour),
			},
		},
		{
			name: "overlapping start",
			override: &Override{
				Start: n.Add(30 * time.Minute),
				Stop:  n.Add(3 * time.Hour),
			},
			expected: ErrOverlappingOverrides,
		},
		{
			name: "overlapping stop",
			override: &Override{
				Start: n.Add(-30 * time.Minute),
				Stop:  n.Add(30 * time.Minute),
			},
			expected: ErrOverlappingOverrides,
		},
		{
			name: "before everything",
			override: &Override{
				Start: n.Add(-2 * time.Hour),
				Stop:  n.Add(-1 * time.Hour),
			},
		},
		{
			name: "stop before start",
			override: &Override{
				Start: n.Add(-1 * time.Hour),
				Stop:  n.Add(-2 * time.Hour),
			},
			expected: ErrStopBeforeStart,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := overrides.Add(tc.override)

			if got != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, got)
			}
		})
	}
}

func TestOverridesAt(t *testing.T) {
	overrides := NewOverrides()

	if got := overrides.Now(); got != nil {
		t.Errorf("expected nothing, got %+v", got)
	}

	n := now()
	o1 := &Override{
		Start: n.Add(1 * time.Hour),
		Stop:  n.Add(2 * time.Hour),
	}
	o2 := &Override{
		Start: n.Add(3 * time.Hour),
		Stop:  n.Add(4 * time.Hour),
	}
	o3 := &Override{
		Start: n.Add(4 * time.Hour),
		Stop:  n.Add(5 * time.Hour),
	}

	for _, o := range []*Override{o2, o1, o3} {
		if err := overrides.Add(o); err != nil {
			t.Fatalf("expected no error, got %s", err)
		}
	}

	// Test sort
	expected := Overrides{o1, o2, o3}
	if !reflect.DeepEqual(expected, overrides) {
		t.Errorf("expected %+v, got %+v", expected, overrides)
	}

	tt := []struct {
		at       time.Time
		expected *Override
		name     string
	}{
		{
			at:       n.Add(1 * time.Hour),
			expected: o1,
			name:     "at start of o1",
		},
		{
			at:       n.Add(2 * time.Hour),
			expected: o1,
			name:     "at stop of o1",
		},
		{
			at:       n.Add(90 * time.Minute),
			expected: o1,
			name:     "in the middle of o1",
		},
		{
			at:   n.Add(-1 * time.Hour),
			name: "before o1",
		},
		{
			at:   n.Add(6 * time.Hour),
			name: "after o3",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := overrides.At(tc.at)

			if got != tc.expected {
				t.Errorf("expected %+v, got %+v", tc.expected, got)
			}
		})
	}
}

func TestOverridesDelete(t *testing.T) {
	overrides := NewOverrides()

	n := now()
	o1 := &Override{
		Start: n.Add(1 * time.Hour),
		Stop:  n.Add(2 * time.Hour),
	}
	o2 := &Override{
		Start: n.Add(3 * time.Hour),
		Stop:  n.Add(4 * time.Hour),
		ID:    "o2",
	}
	o3 := &Override{
		Start: n.Add(4 * time.Hour),
		Stop:  n.Add(5 * time.Hour),
		ID:    "o3",
	}

	for _, o := range []*Override{o2, o3, o1} {
		if err := overrides.Add(o); err != nil {
			t.Fatalf("expected no error, got %s", err)
		}
	}

	expected := Overrides{o1, o2, o3}
	if !reflect.DeepEqual(expected, overrides) {
		t.Errorf("expected %+v, got %+v", expected, overrides)
	}

	tt := []struct {
		name     string
		id       string
		err      error
		expected Overrides
	}{
		{
			name:     "removing o2",
			id:       o2.ID,
			expected: Overrides{o1, o3},
		},
		{
			name:     "removing o3",
			id:       o3.ID,
			expected: Overrides{o1},
		},
		{
			name:     "removing invalid o4",
			id:       "o4",
			expected: Overrides{o1},
			err:      ErrOverrideNotFound,
		},
		{
			name:     "removing o1",
			id:       o1.ID,
			expected: Overrides{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := overrides.Delete(tc.id)

			if !reflect.DeepEqual(tc.expected, overrides) {
				t.Errorf("expected %+v, got %+v", tc.expected, got)
			}

			if got != tc.err {
				t.Errorf("expected %s, got %s", tc.err, got)
			}
		})
	}
}

func TestOverridesCleanup(t *testing.T) {
	overrides := NewOverrides()

	n := now()
	o1 := &Override{
		Start: n.Add(-3 * time.Hour),
		Stop:  n.Add(-2 * time.Hour),
	}
	o2 := &Override{
		Start: n.Add(-2 * time.Hour),
		Stop:  n.Add(-1 * time.Hour),
	}
	o3 := &Override{
		Start: n.Add(1 * time.Hour),
		Stop:  n.Add(2 * time.Hour),
	}
	o4 := &Override{
		Start: n.Add(2 * time.Hour),
		Stop:  n.Add(3 * time.Hour),
	}

	// With nothing
	overrides.cleanup()
	expected := Overrides{}
	if !reflect.DeepEqual(expected, overrides) {
		t.Fatalf("expected %+v, got %+v", expected, overrides)
	}

	// With only old stuff
	overrides = Overrides{o1, o2}
	overrides.cleanup()
	expected = Overrides{}
	if !reflect.DeepEqual(expected, overrides) {
		t.Fatalf("expected %+v, got %+v", expected, overrides)
	}

	// With old and new stuff
	overrides.Add(o1)
	overrides.Add(o4)
	overrides.Add(o3)
	overrides.Add(o2)
	overrides.cleanup()
	expected = Overrides{o3, o4}
	if !reflect.DeepEqual(expected, overrides) {
		t.Fatalf("expected %+v, got %+v", expected, overrides)
	}
}

func TestOverridesNextTime(t *testing.T) {

	n := now()
	o1 := &Override{
		Start: n.Add(1 * time.Hour),
		Stop:  n.Add(2 * time.Hour),
	}
	o2 := &Override{
		Start: n.Add(2 * time.Hour),
		Stop:  n.Add(3 * time.Hour),
	}
	o3 := &Override{
		Start: n.Add(-1 * time.Hour),
		Stop:  n.Add(1 * time.Hour),
	}

	tt := []struct {
		name     string
		expected *time.Time
		toAdd    []*Override
	}{
		{
			name: "empty",
		},
		{
			name:  "only one in progress",
			toAdd: []*Override{o3},
		},
		{
			name:     "first one",
			expected: &o1.Start,
			toAdd:    []*Override{o1},
		},
		{
			name:     "second one",
			expected: &o1.Start,
			toAdd:    []*Override{o3, o1},
		},
		{
			name:     "classic case",
			expected: &o1.Start,
			toAdd:    []*Override{o3, o1, o2},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			overrides := NewOverrides()
			for _, o := range tc.toAdd {
				if err := overrides.Add(o); err != nil {
					t.Fatalf("expected no error, got %s", err)
				}
			}

			got := overrides.NextTime()
			if got != nil && tc.expected != nil {
				if !got.Equal(*tc.expected) {
					t.Errorf("expected %q, got %q", tc.expected, got)
				}
			} else {
				if got != tc.expected {
					t.Errorf("expected %q, got %q", tc.expected, got)
				}
			}
		})
	}
}
