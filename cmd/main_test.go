package main

import (
	"flag"
	"io"
	"testing"
	"time"
)

func TestLeaderElectionDefaults(t *testing.T) {
	d := defaultLeaderElectionDurations()
	if d.lease != 15*time.Second || d.renew != 10*time.Second || d.retry != 2*time.Second {
		t.Fatalf("upstream defaults changed: %+v", d)
	}
	if err := d.validate(); err != nil {
		t.Fatal(err)
	}
}

func TestLeaderElectionConfiguredFlags(t *testing.T) {
	d := defaultLeaderElectionDurations()
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	f.SetOutput(io.Discard)
	d.bindFlags(f)
	if err := f.Parse([]string{"--leader-elect-lease-duration=60s", "--leader-elect-renew-deadline=40s", "--leader-elect-retry-period=10s"}); err != nil {
		t.Fatal(err)
	}
	if err := d.validate(); err != nil {
		t.Fatal(err)
	}
	if d.lease != time.Minute || d.renew != 40*time.Second || d.retry != 10*time.Second {
		t.Fatalf("flags did not propagate: %+v", d)
	}
}

func TestLeaderElectionValidation(t *testing.T) {
	cases := []struct {
		name                string
		lease, renew, retry time.Duration
		valid               bool
	}{
		{"default", 15 * time.Second, 10 * time.Second, 2 * time.Second, true},
		{"nested", 60 * time.Second, 40 * time.Second, 10 * time.Second, true},
		{"maximum", 10 * time.Minute, 5 * time.Minute, time.Minute, true},
		{"zero lease", 0, 10 * time.Second, 2 * time.Second, false},
		{"zero renew", 15 * time.Second, 0, 2 * time.Second, false},
		{"zero retry", 15 * time.Second, 10 * time.Second, 0, false},
		{"negative", 15 * time.Second, 10 * time.Second, -time.Second, false},
		{"equal lease renew", 10 * time.Second, 10 * time.Second, 2 * time.Second, false},
		{"lease shorter", 9 * time.Second, 10 * time.Second, 2 * time.Second, false},
		{"equal jitter", 15 * time.Second, 12 * time.Second, 10 * time.Second, false},
		{"below jitter", 15 * time.Second, 11 * time.Second, 10 * time.Second, false},
		{"lease too long", 11 * time.Minute, 10 * time.Second, 2 * time.Second, false},
		{"renew too long", 15 * time.Second, 11 * time.Minute, 2 * time.Second, false},
		{"retry too long", 15 * time.Second, 10 * time.Second, 11 * time.Minute, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := leaderElectionDurations{lease: tc.lease, renew: tc.renew, retry: tc.retry}
			err := d.validate()
			if (err == nil) != tc.valid {
				t.Fatalf("valid=%t error=%v", tc.valid, err)
			}
		})
	}
}

func TestInvalidDurationFlag(t *testing.T) {
	d := defaultLeaderElectionDurations()
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	f.SetOutput(io.Discard)
	d.bindFlags(f)
	if err := f.Parse([]string{"--leader-elect-lease-duration=not-a-duration"}); err == nil {
		t.Fatal("invalid duration was accepted")
	}
}
