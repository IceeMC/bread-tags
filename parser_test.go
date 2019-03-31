package main

import (
	"math/rand"
	"testing"
)

func TestParser(t *testing.T) {
	parser := Parser{}
	t.Run("args", func(t *testing.T) {
		expected := "a b c"
		context := make(map[string]interface{})
		context["args"] = []string{"a", "b", "c"}
		result := parser.Parse("{args}", context)
		if result != expected {
			t.Errorf("Test failed, expected %s, got %s", expected, result)
		}
	})
	t.Run("args with joiner", func(t *testing.T) {
		expected := "a;b;c"
		context := make(map[string]interface{})
		context["args"] = []string{"a", "b", "c"}
		context["joiner"] = ";"
		result := parser.Parse("{args}", context)
		if result != expected {
			t.Errorf("Test failed, expected %s, got %s", expected, result)
		}
	})
	t.Run("args aliases", func(t *testing.T) {
		expected := "a b c"
		context := make(map[string]interface{})
		context["args"] = []string{"a", "b", "c"}
		result := parser.Parse("{allargs}", context)
		if result != expected {
			t.Errorf("Test failed, expected %s, got %s", expected, result)
		}
	})
	t.Run("capitalize", func(t *testing.T) {
		expected := "Abc"
		result := parser.Parse("{capitalize:abc}", nil)
		if result != expected {
			t.Errorf("Test failed, expected %s, got %s", expected, result)
		}
	})
	t.Run("uppercase", func(t *testing.T) {
		expected := "HI"
		result := parser.Parse("{uppercase:HI}", nil)
		if result != expected {
			t.Errorf("Test failed, expected %s, got %s", expected, result)
		}
	})
	emptyString := &Tag{
		Name: "emptystring",
		Run: func(value string, context map[string]interface{}) string {
			return ""
		},
		Aliases: []string{"emptystr"},
	}
	parser.LoadTags(emptyString)
	t.Run("empty string", func(t *testing.T) {
		expected := ""
		result := parser.Parse("{emptystring}", nil)
		if result != expected {
			t.Errorf("Test failed, expected nothing, got %s", result)
		}
	})
}

func BenchmarkParser(b *testing.B) {
	parser := Parser{}
	b.Run("args", func(b *testing.B) {
		context := make(map[string]interface{})
		context["args"] = []string{"a", "b", "c"}
		for i := 0; i < b.N; i++ {
			parser.Parse("{args}", context)
		}
	})
	b.Run("args with joiner", func(b *testing.B) {
		context := make(map[string]interface{})
		context["args"] = []string{"a", "b", "c"}
		context["joiner"] = ";"
		for i := 0; i < b.N; i++ {
			parser.Parse("{args}", context)
		}
	})
	b.Run("capitalize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			parser.Parse("{capitalize:abc}", nil)
		}
	})
	b.Run("uppercase", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			parser.Parse("{uppercase:abc}", nil)
		}
	})
	b.Run("load tag", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tag := &Tag{
				Name: string(rand.Int()),
				Run: func(value string, context map[string]interface{}) string {
					return ""
				},
			}
			parser.LoadTag(tag)
		}
	})
}
