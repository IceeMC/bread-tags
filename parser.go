package main

import (
	"errors"
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

type Parser struct {
}

func (parser Parser) GetTags(t *Tag) map[string]*Tag {
	return Tags
}

func (parser Parser) LoadTag(t *Tag) (*Tag, error) {
	if t == nil || t.Name == "" {
		return nil, errors.New("expected a tag")
	}
	if _, exists := Tags[t.Name]; exists {
		return nil, fmt.Errorf("%s already exists", t.Name)
	}
	Tags[t.Name] = t
	if len(t.Aliases) > 0 {
		for _, a := range t.Aliases {
			Tags[a] = t
		}
	}
	return t, nil
}

func LoadTagNoParser(t *Tag) (*Tag, error) {
	if t == nil || t.Name == "" {
		return nil, errors.New("expected a tag")
	}
	if _, exists := Tags[t.Name]; exists {
		return nil, fmt.Errorf("%s already exists", t.Name)
	}
	Tags[t.Name] = t
	if len(t.Aliases) > 0 {
		for _, a := range t.Aliases {
			Tags[a] = t
		}
	}
	return t, nil
}

func (parser Parser) LoadTags(tags ...*Tag) []*Tag {
	var out []*Tag
	for _, tag := range tags {
		tag, _ := parser.LoadTag(tag)
		if tag != nil {
			out = append(out, tag)
		} else {
			out = append(out, nil)
		}
	}
	return out
}

// A function to run tags
func RunTag(token Token, context map[string]interface{}) string {
	name := ""
	value := ""
	if token.Child != nil {
		name = strings.Split(token.Text, ":")[0]
		value = RunTag(*token.Child, context)
	} else {
		s := strings.Split(token.Text, ":")
		name = s[0]
		if len(s) > 1 {
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

func GetTokenType(c string) int {
	t := 1
	if c != "{" {
		t = 0
	}
	return t
}

func MakeToken(parent *Token, t int) *Token {
	return &Token{Parent: parent, Child: nil, Type: t, Text: ""}
}

// A function to scan for tokens
func Lex(tag string) []Token {
	tok := MakeToken(nil, GetTokenType(string(tag[:1])))
	if tok.Type == 0 {
		tok.Text += string(tag[:1])
	}
	var tokens []Token
	for i := 1; i < len(tag); i++ {
		c := []rune(tag)[i]
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
						if i+1 == len(tag) {
							tok = MakeToken(nil, GetTokenType(string(tag[i])))
						} else {
							tok = MakeToken(nil, GetTokenType(string(tag[i+1])))
						}
					} else {
						tok = tok.Parent
					}
				}
				break
			}
		default:
			if tok.Type == 0 {
				tok.Type = 0
			}
			tok.Text += string(c)
		}
	}
	return tokens
}

func (parser Parser) Parse(input string, context map[string]interface{}) string {
	out := ""
	tokens := Lex(input)
	for _, token := range tokens {
		if token.Type == 0 {
			out += token.Text
		} else {
			out += RunTag(token, context)
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
		Name: "capitalize",
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
	LoadTagNoParser(argsTag)
	LoadTagNoParser(capitalizeTag)
	LoadTagNoParser(chooseTag)
	LoadTagNoParser(rangeTag)
	LoadTagNoParser(upperCase)
}
