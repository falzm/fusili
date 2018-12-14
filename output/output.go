package output

type Output interface {
	Report(map[string][]int) error
}

var (
	Outputs = make(map[string]func(string, map[string]interface{}) (Output, error))
)
