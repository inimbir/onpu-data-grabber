package commands

import (
	"fmt"
	"github.com/inimbir/onpu-data-grabber/app/clients"
	"strings"
)

type AddCriteriaCommand struct {
	Group   string
	Pattern string
}

func (c *AddCriteriaCommand) Matches(pattern string) bool {
	command := strings.Split(pattern, " ")
	return command[0] == "add"
}

func (c *AddCriteriaCommand) Parse(pattern string) (err error) {
	command := strings.Split(pattern, " ")
	for i, params := range command {
		if params == "-g" && i+1 < len(command) {
			c.Group = command[i+1]
		}
		if params == "-n" && i+1 < len(command) {
			c.Criteria = command[i+1]
		}
	}
	if c.Group == "" || c.Criteria == "" {
		return fmt.Errorf("Cannot add criteria [%s] to group [%s], because one of them is empty\n", c.Criteria, c.Group)
	}
	return
}

func (c *AddCriteriaCommand) Execute() error {
	return clients.InsertCriteria(c)
}

func (c *AddCriteriaCommand) String() string {
	return fmt.Sprintf("add -g group-name -n pattern-name")
}
