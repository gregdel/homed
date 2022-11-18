package schedule

import (
	"encoding/json"
	"testing"
)

func TestTimeString(t *testing.T) {
	time := NewTime(9, 12, 3)
	got := time.String()
	expected := "09:12:03"
	if got != expected {
		t.Errorf("got %s, expected %s", got, expected)
	}
}

func TestTimeBefore(t *testing.T) {
	tt := []struct {
		name     string
		first    Time
		second   Time
		expected bool
	}{
		{
			name:     "time before",
			first:    NewTime(10, 30, 00),
			second:   NewTime(11, 30, 00),
			expected: true,
		},
		{
			name:     "time after",
			first:    NewTime(23, 30, 00),
			second:   NewTime(13, 30, 00),
			expected: false,
		},
		{
			name:     "same hour before",
			first:    NewTime(10, 00, 00),
			second:   NewTime(10, 30, 00),
			expected: true,
		},
		{
			name:     "same hour after",
			first:    NewTime(10, 30, 00),
			second:   NewTime(10, 00, 00),
			expected: false,
		},
		{
			name:     "same hour and minute",
			first:    NewTime(10, 30, 00),
			second:   NewTime(10, 30, 30),
			expected: true,
		},
		{
			name:     "exact same time",
			first:    NewTime(10, 30, 30),
			second:   NewTime(10, 30, 30),
			expected: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.first.Before(tc.second)
			if got != tc.expected {
				word := "before"
				if !tc.expected {
					word = "after"
				}

				t.Errorf(
					"%s should be %s %s",
					tc.first.String(),
					word,
					tc.second.String(),
				)
			}

			after := tc.first.After(tc.second)
			if got == after {
				t.Errorf("after should be the opposite of before")
			}
		})
	}
}

func TestTimeMarshalJSON(t *testing.T) {
	tt := []struct {
		Time     *Time
		expected string
	}{
		{
			Time:     NewTimePointer(10, 30, 00),
			expected: "\"10:30:00\"",
		},
		{
			Time:     NewTimePointer(23, 30, 17),
			expected: "\"23:30:17\"",
		},
		{
			Time:     NewTimePointer(8, 3, 4),
			expected: "\"08:03:04\"",
		},
	}

	for _, tc := range tt {
		t.Run(tc.expected, func(t *testing.T) {
			out, err := json.Marshal(tc.Time)
			if err != nil {
				t.Fatalf("failed to marshal json: %s", err.Error())
			}

			if string(out) != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, string(out))
			}
		})
	}
}

func TestTimeUnmarshalJSON(t *testing.T) {
	tt := []struct {
		expected          *Time
		input             []byte
		expectedErr       error
		jsonEncodingError bool
	}{
		{
			input:    []byte("\"10:30:00\""),
			expected: NewTimePointer(10, 30, 00),
		},
		{
			input:    []byte("\"10:30\""),
			expected: NewTimePointer(10, 30, 00),
		},
		{
			input:    []byte("\"8:2\""),
			expected: NewTimePointer(8, 2, 00),
		},
		{
			input:       []byte("\"xxxxx\""),
			expectedErr: ErrInvalidTime,
		},
		{
			input:       []byte("\"xx:xx:xx\""),
			expectedErr: ErrInvalidTime,
		},
		{
			input:             []byte("\""),
			jsonEncodingError: true,
		},
	}

	for _, tc := range tt {
		t.Run(string(tc.input), func(t *testing.T) {
			got := &Time{}
			err := got.UnmarshalJSON(tc.input)
			if err != nil {
				_, ok := err.(*json.SyntaxError)
				if tc.jsonEncodingError && !ok {
					t.Fatalf("expected JSON encoding error got %s", err)
					return
				}

				if !tc.jsonEncodingError && tc.expectedErr != err {
					t.Fatalf("expected %s, got %s", tc.expectedErr, err)
				}
			}

			if tc.expected == nil {
				return
			}

			if *got != *tc.expected {
				t.Errorf("expected %+v, got %+v", tc.expected, got)
			}
		})
	}
}
