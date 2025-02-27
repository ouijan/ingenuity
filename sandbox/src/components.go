package main

// Aether:Component
type DemoComponent1 struct {
	Name string
	Age  int
}

type DemoCategory int

// Aether:Enum
const (
	FirstCategory DemoCategory = iota
	SecondCategory
	ThirdCategory
)

// Aether:Component
type DemoComponent2 struct {
	Category DemoCategory `json:"category"`
}
