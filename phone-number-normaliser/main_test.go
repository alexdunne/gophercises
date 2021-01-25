package main

import "testing"

func TestNormaliseWithCorrectFormat(t *testing.T) {
	input := "1234567890"
	output := Normalise(input)

	want := "1234567890"

	if want != output {
		t.Fatalf(`Normalise("%s") = %s, want match for %s`, input, output, want)
	}
}

func TestNormaliseWithSpaces(t *testing.T) {
	input := "123 456 7891"
	output := Normalise(input)

	want := "1234567891"

	if want != output {
		t.Fatalf(`Normalise("%s") = %s, want match for %s`, input, output, want)
	}
}

func TestNormaliseWithBrackets(t *testing.T) {
	input := "(123) 456 7892"
	output := Normalise(input)

	want := "1234567892"

	if want != output {
		t.Fatalf(`Normalise("%s") = %s, want match for %s`, input, output, want)
	}
}
