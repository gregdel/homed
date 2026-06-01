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
		time *time.Time
		name string
	}{
		{
			time: nil,
			name: "nil time",
		},
		{
			time: &t1,
			name: "non nil time",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			toEncode := struct {
				Name string
				Time *time.Time
			}{
				Name: tc.name,
				Time: tc.time,
			}

			data, err := json.Marshal(toEncode)
			require.Nil(t, err)

			var got Test
			err = json.Unmarshal(data, &got)
			require.Nil(t, err)

			assert.Equal(t, got.Name, tc.name)

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
		Time Time   `json:"time"`
	}

	t1 := time.Now()

	tt := []struct {
		expected *time.Time
		name     string
		time     *time.Time
	}{
		{
			expected: nil,
			name:     "nil time",
			time:     nil,
		},
		{
			expected: &t1,
			name:     "non nil time",
			time:     &t1,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			toEncode := Test{
				Name: tc.name,
			}
			toEncode.Time.Store(tc.time)

			data, err := json.Marshal(&toEncode)
			require.Nil(t, err)

			toDecode := struct {
				Name string     `json:"name"`
				Time *time.Time `json:"time"`
			}{}

			err = json.Unmarshal(data, &toDecode)
			require.Nil(t, err)

			assert.Equal(t, tc.name, toDecode.Name)
			if tc.expected == nil {
				assert.Nil(t, toDecode.Time)
			} else {
				require.NotNil(t, toDecode.Time)
				assert.True(t, tc.expected.Equal(*toDecode.Time))
			}
		})
	}
}

func TestTimeConcurrentInitialLoadStore(t *testing.T) {
	var tm Time
	value := time.Now()
	done := make(chan struct{})

	for range 100 {
		go func() {
			defer func() { done <- struct{}{} }()
			for range 100 {
				tm.Store(&value)
				_ = tm.Load()
				tm.Store(nil)
			}
		}()
	}

	for range 100 {
		<-done
	}
}
