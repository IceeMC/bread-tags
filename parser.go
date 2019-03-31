// Copyright (c) 2019 Soumil07. All rights reserved. BSD-3 License
//

package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Tag struct {
	Name    string   `json:"name"`
	Aliases []string `json:"aliases"`
	Run     func(value string, context map[string]interface{}) string
}

type Token struct {
	Parent *Token `json:"parent"`
	Type   int    `json:"type"`
	Child  *Token `json:"child"`
	Text   string `json:"text"`
}

var Tags = map[string]*Tag{}

func LoadTags(tags ...*Tag) {
	for _, t := range tags {
		Tags[t.Name] = t
		if len(t.Aliases) > 0 {
			for _, a := range t.Aliases {
				Tags[a] = t
			}
		}
	}
}

// runTag executes the given tag
func runTag(token Token, context map[string]interface{}) string {
	name := ""
	value := ""
	if token.Child != nil {
		name = strings.Split(token.Text, ":")[0]
		value = runTag(*token.Child, context)
	} else {
		s := strings.Split(token.Text, ":")
		name = s[0]
		if len(s) == 2 {
			value = s[1]
		}
	}
	t, exists := Tags[name]
	if value == "" {
		value = ""
	}
	if !exists {
		return fmt.Sprintf("{%s:%s}", name, value)
	}
	return t.Run(value, context)
}

func getTokenType(c string) int {
	t := 1
	if c != "{" {
		t = 0
	}
	return t
}

func MakeToken(parent *Token, t int) *Token {
	return &Token{Parent: parent, Child: nil, Type: t, Text: ""}
}

// lex scans the input string and generates a stream of tokens
func lex(tag string) []Token {
	tok := MakeToken(nil, getTokenType(string(tag[:1])))
	if tok.Type == 0 {
		tok.Text += string(tag[:1])
	}
	var tokens []Token
	for i := 1; i < len(tag); i++ {
		c := tag[i]
		switch {
		case c == '{':
			{
				if tok.Type == 0 {
					tokens = append(tokens, *tok)
					tok = MakeToken(nil, 1)
				} else {
					tok.Child = MakeToken(tok, 1)
					tok = tok.Child
				}
				break
			}
		case c == '}':
			{
				if tok.Type == 1 {
					if tok.Parent == nil {
						tokens = append(tokens, *tok)
						if i+1 == len(tag) { // Fix for list out of range
							tok = MakeToken(nil, getTokenType(string(tag[i])))
						} else {
							tok = MakeToken(nil, getTokenType(string(tag[i+1])))
						}
					} else {
						tok = tok.Parent
					}
				}
				break
			}
		default:
			tok.Text += string(c)
		}
	}
	return tokens
}

func Parse(input string, context map[string]interface{}) string {
	out := ""
	tokens := lex(input)
	for _, token := range tokens {
		if token.Type == 0 {
			out += token.Text
		} else {
			out += runTag(token, context)
		}
	}
	return out
}

func init() {
	argsTag := &Tag{
		Name: "args",
		Run: func(value string, context map[string]interface{}) string {
			if context == nil || context["args"] == nil {
				return "No arguments passed."
			}
			joinChar := " "
			if context["joiner"] != nil {
				joinChar = context["joiner"].(string)
			}
			return strings.Join(context["args"].([]string), joinChar)
		},
		Aliases: []string{"allargs"},
	}
	capitalizeTag := &Tag{
		Name: "capitalize",
		Run: func(value string, context map[string]interface{}) string {
			return strings.ToUpper(string(value[0])) + value[1:]
		},
		Aliases: []string{"titlecase"},
	}
	chooseTag := &Tag{
		Name: "choose",
		Run: func(value string, context map[string]interface{}) string {
			rand.Seed(time.Now().Unix())
			choices := strings.Split(value, ";")
			return choices[rand.Intn(len(choices))]
		},
		Aliases: []string{"choice"},
	}
	rangeTag := &Tag{
		Name: "range",
		Run: func(value string, context map[string]interface{}) string {
			s := strings.Split(value, ";")
			min, _ := strconv.ParseInt(s[0], 10, 64)
			max := 0
			if len(s) == 2 {
				m, _ := strconv.ParseInt(s[1], 10, 64)
				max = int(m)
			}
			rand.Seed(time.Now().Unix())
			return fmt.Sprintf("%f", math.Floor(float64(rand.Int())*float64(max))+float64(min))
		},
	}
	upperCase := &Tag{
		Name: "uppercase",
		Run: func(value string, context map[string]interface{}) string {
			return strings.ToUpper(value)
		},
		Aliases: []string{"upper"},
	}

	LoadTags(argsTag, capitalizeTag, chooseTag, rangeTag, upperCase)
}
