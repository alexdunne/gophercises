package main

import "testing"

type normaliseTestCase struct {
	input string
	want  string
}

func TestNormalise(t *testing.T) {
	testCases := []normaliseTestCase{
		{input: "1234567890", want: "1234567890"},
		{input: "123 456 7891", want: "1234567891"},
		{input: "(123) 456 7892", want: "1234567892"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			actual := normalise(testCase.input)

			if actual != testCase.want {
				t.Errorf("got %s; want %s", actual, testCase.want)
			}
		})
	}
}
