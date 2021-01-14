package schedule

import "testing"

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
				t.Errorf("after should be the oposite of before")
			}
		})
	}
}
