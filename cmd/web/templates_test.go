package main

import (
	"testing"
	"time"

	"github.com/mhrdini/snippetbox/internal/assert"
)

func TestPrettyDate(t *testing.T) {

	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 6, 2, 3, 27, 0, 0, time.UTC),
			want: "02 Jun 2024 at 03:27",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2022, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2022 at 09:15",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hd := prettyDate(test.tm)
			assert.Equal(t, hd, test.want)
		})
	}

}
