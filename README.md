# bread-tags
bread-tags is a minimal, customizable, fast tas parser for discord bots.

# Installation
`go get https://github.com/IceeMC/bread-tags`

# Usage
```go
// This code assumes you have defined it in a main/similar function
parser := Parser{}
context := make(map[string]interface{})
context["args"] = []string{"Ice"}
parser.Parse("Hi, {args}", context) // Hi, Ice
```

# Creating custom tags
```go
// Create the tag
tag := &Tag{
	Name: "mytag",
	Run: func (value string, context map[string]interface{}) string {
		// What to do when the tag is ran
		return value
	}
}
_, err := parser.LoadTag(tag)
if err != nil {
	log.Printf("Failed to load tag: %s", tag.Name)
}
// Using the tag
parser.Parse("{mytag:Some value}", nil) // Some value
```