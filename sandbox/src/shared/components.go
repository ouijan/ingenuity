package shared

// Ingenuity:Component
type DemoComponent1 struct {
	Name string
	Age  int
}

type DemoCategory int

// Ingenuity:Enum
const (
	FirstCategory DemoCategory = iota
	SecondCategory
	ThirdCategory
)

// Ingenuity:Component
type DemoComponent2 struct {
	Category DemoCategory `json:"category"`
}
