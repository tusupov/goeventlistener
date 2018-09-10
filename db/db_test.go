package db

import (
	"context"
	"testing"

	"gopkg.in/h2non/gock.v1"
)

func TestStorage_Add(t *testing.T) {

	testStorage := New()

	testCase := []struct {
		b bool
		r ListenerRequest
	}{
		{
			true,
			ListenerRequest{
				"event1",
				"listener1",
				"address1",
			},
		},
		{
			true,
			ListenerRequest{
				"event1",
				"listener2",
				"address2",
			},
		},
		{
			true,
			ListenerRequest{
				"event1",
				"listener3",
				"address3",
			},
		},
		{
			true,
			ListenerRequest{
				"event2",
				"listener4",
				"address4",
			},
		},
		{
			true,
			ListenerRequest{
				"event2",
				"listener5",
				"address5",
			},
		},
		{
			true,
			ListenerRequest{
				"event2",
				"listener6",
				"address6",
			},
		},
		{
			false,
			ListenerRequest{
				"event1",
				"listener1",
				"address1",
			},
		},
		{
			false,
			ListenerRequest{
				"event3",
				"listener1",
				"address1",
			},
		},
	}

	for _, test := range testCase {
		err := testStorage.Add(test.r)
		if (test.b && err != nil) || (!test.b && err == nil) {
			t.Fatalf("Expect %t, but error '%v'\nTest: %v", test.b, err, test)
		}
	}

}

func TestStorage_DeleteListener(t *testing.T) {

	testStorage := New()
	testStorage.Add(ListenerRequest{"event1", "listener1", "address1"})
	testStorage.Add(ListenerRequest{"event1", "listener2", "address2"})
	testStorage.Add(ListenerRequest{"event1", "listener3", "address3"})
	testStorage.Add(ListenerRequest{"event2", "listener4", "address4"})
	testStorage.Add(ListenerRequest{"event2", "listener5", "address5"})
	testStorage.Add(ListenerRequest{"event3", "listener6", "address6"})

	testCase := []struct {
		d bool
		l string
	}{
		{
			true,
			"listener1",
		},
		{
			true,
			"listener2",
		},
		{
			true,
			"listener3",
		},
		{
			false,
			"listener1",
		},
		{
			true,
			"listener4",
		},
		{
			true,
			"listener5",
		},
		{
			false,
			"listener5",
		},
		{
			true,
			"listener6",
		},
		{
			false,
			"listener7",
		},
	}

	for _, test := range testCase {
		err := testStorage.DeleteListener(test.l)
		if (test.d && err != nil) || (!test.d && err == nil) {
			t.Fatalf("Expect %t, but error '%v'\nTest: %v", test.d, err, test)
		}
	}

}

func TestStorage_Publish(t *testing.T) {

	defer gock.Off()

	testStorage := New()
	list := []ListenerRequest{
		{"event1", "listener1", "https://address1/"},
		{"event1", "listener2", "https://address2/"},
		{"event1", "listener3", "https://address3/"},
		{"event2", "listener4", "https://address4/"},
		{"event2", "listener5", "https://address5/"},
		{"event3", "listener6", "https://address6/"},
	}
	for _, r := range list {
		testStorage.Add(r)
		if r.Listener == "listener2" {
			gock.New(r.Address).Reply(404)
		} else {
			gock.New(r.Address).Reply(200)
		}
	}

	var err error

	err = testStorage.Publish(context.Background(), "event1")
	if err == nil {
		t.Fatalf("Error must be not nil, but nil")
	}

	err = testStorage.Publish(context.Background(), "event2")
	if err != nil {
		t.Fatalf("Error must be nil, but '%v'", err)
	}

	err = testStorage.Publish(context.Background(), "event0")
	if err == nil {
		t.Fatalf("Error must be not nil, but nil")
	}

}
