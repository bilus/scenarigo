package context

import (
	"github.com/bilus/scenarigo/template"
)

// ExecuteTemplate executes template strings in context.
func (c *Context) ExecuteTemplate(i interface{}) (interface{}, error) {
	return template.Execute(i, c)
}
