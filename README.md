# bread-tags
bread-tags is a minimal, customizable, fast tas parser for discord bots.

# Installation
`go get https://github.com/IceeMC/bread-tags`

# Usage
```go
// This code assumes you have defined it in a main/similar function
context := make(map[string]interface{})
context["args"] = []string{"Ice"}
breadtags.Parse("Hi, {args}", context) // Hi, Ice
```

# Creating custom tags
```go
// Create the tag
tag := &breadtags.Tag{
	Name: "mytag",
	Run: func (value string, context map[string]interface{}) string {
		// What to do when the tag is ran
		return value
	}
}
breadtags.LoadTags(tag)
// Using the tag
breadtags.Parse("{mytag:Some value}", nil) // Some value
```
