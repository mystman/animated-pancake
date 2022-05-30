package main

import "testing"

func TestMain(t *testing.T) {
	testCases := []struct {
		desc  string
		input interface{}
		want  interface{}
	}{
		{
			desc: "Test 1",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			//main()
		})
	}
}
