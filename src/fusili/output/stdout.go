package output

import (
	"sort"

	"fusili/logger"
)

type StdoutOutput struct {
	name string
}

func init() {
	Outputs["stdout"] = func(name string, settings map[string]interface{}) (Output, error) {
		o := &StdoutOutput{name: name}

		return o, nil
	}
}

func (o *StdoutOutput) Report(hosts map[string][]int) error {
	for host, hostPorts := range hosts {
		sort.Ints(hostPorts)

		for _, port := range hostPorts {
			logger.Warning("report", "%s: found port %d/tcp open", host, port)
		}
	}

	return nil
}
