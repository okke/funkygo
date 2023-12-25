package main

import "testing"

func TestMap2Struct(t *testing.T) {

	type Person struct {
		Name string
		Age  int
		City string
	}

	data := map[string]interface{}{
		"Name": "John",
		"Age":  30,
		"City": "New York",
		"Foo":  "Bar",
	}

	result, err := map2struct(data, &Person{})
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if result.Name != "John" {
		t.Errorf("Expected John, got %s", result.Name)
	}

	if result.Age != 30 {
		t.Errorf("Expected 30, got %d", result.Age)
	}

	if result.City != "New York" {
		t.Errorf("Expected New York, got %s", result.City)
	}
}
