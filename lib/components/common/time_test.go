package common

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalTime(t *testing.T) {
	type Test struct {
		Name string `json:"name"`
		Time Time   `json:"time"`
	}

	t1 := time.Now()

	tt := []struct {
		time     *time.Time
		expected Test
	}{
		{
			time: nil,
			expected: Test{
				Name: "nil time",
			},
		},
		{
			time: &t1,
			expected: Test{
				Name: "non nil time",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.expected.Name, func(t *testing.T) {
			toEncode := struct {
				Name string
				Time *time.Time
			}{
				Name: tc.expected.Name,
				Time: tc.time,
			}

			data, err := json.Marshal(toEncode)
			require.Nil(t, err)

			var got Test
			err = json.Unmarshal(data, &got)
			require.Nil(t, err)

			assert.Equal(t, got.Name, tc.expected.Name)

			v := got.Time.Load()
			if tc.time == nil {
				assert.Nil(t, v)
			} else {
				assert.NotNil(t, v)
				assert.True(t, tc.time.Equal(*v))
			}
		})
	}
}

func TestMarshalTime(t *testing.T) {
	type Test struct {
		Name string `json:"name"`
		Time Time   `json:"time,omitempty"`
	}

	t1 := time.Now()
	var T1 Time
	var T2 Time
	T2.Store(&t1)

	tt := []struct {
		expected *time.Time
		time     Test
	}{
		{
			expected: nil,
			time: Test{
				Name: "nil time",
				Time: T1,
			},
		},
		{
			expected: &t1,
			time: Test{
				Name: "non nil time",
				Time: T2,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.time.Name, func(t *testing.T) {
			data, err := json.Marshal(tc.time)
			require.Nil(t, err)

			toDecode := struct {
				Name string     `json:"name"`
				Time *time.Time `json:"time"`
			}{}

			err = json.Unmarshal(data, &toDecode)
			require.Nil(t, err)

			assert.Equal(t, tc.time.Name, toDecode.Name)
			assert.Equal(t, tc.expected, tc.time.Time.Load())
		})
	}
}
