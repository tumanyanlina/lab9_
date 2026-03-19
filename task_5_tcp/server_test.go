package main

import (
	"reflect"
	"testing"
)

func TestProcessNumbers_NormalInput(t *testing.T) {
	input := []int{1, 2, 3}
	wantSum := 1*1 + 2*2 + 3*3 // = 14

	resp := processNumbers(input)
	if resp.Error != "" {
		t.Errorf("unexpected error: %v", resp.Error)
	}
	if resp.Sum != wantSum {
		t.Errorf("got sum %v, want %v", resp.Sum, wantSum)
	}
	if !reflect.DeepEqual(resp.Original, input) {
		t.Errorf("got original %v, want %v", resp.Original, input)
	}
}

func TestProcessNumbers_EmptyArray(t *testing.T) {
	input := []int{}
	resp := processNumbers(input)
	if resp.Error != "no numbers provided" {
		t.Errorf("got error %v, want 'no numbers provided'", resp.Error)
	}
	if resp.Sum != 0 {
		t.Errorf("got sum %v, want 0", resp.Sum)
	}
	if resp.Original != nil && len(resp.Original) > 0 {
		t.Errorf("got original %v, want nil or empty", resp.Original)
	}
}

func TestProcessNumbers_NumberTooLarge(t *testing.T) {
	input := []int{10, 2000, 2}
	resp := processNumbers(input)
	if resp.Error != "number too large" {
		t.Errorf("got error %v, want 'number too large'", resp.Error)
	}
	if resp.Sum != 0 {
		t.Errorf("got sum %v, want 0", resp.Sum)
	}
	// For errors, the instruction does not specify about original, but let's check it is not filled
	if resp.Original != nil && len(resp.Original) > 0 {
		t.Errorf("got original %v, want nil or empty", resp.Original)
	}
}
