package main

// A CommandGroup represents a set of commands about a common topic.
type CommandGroup struct {
	name     string
	commands []Command
}

// Name returns the group name
func (g *CommandGroup) Name() string {
	return g.name
}

func (g CommandGroup) Len() int           { return len(g.commands) }
func (g CommandGroup) Less(i, j int) bool { return g.commands[i].Name() < g.commands[j].Name() }
func (g CommandGroup) Swap(i, j int)      { g.commands[i], g.commands[j] = g.commands[j], g.commands[i] }

// Sorting of a slice of command groups.
type byGroupName []*CommandGroup

func (p byGroupName) Len() int           { return len(p) }
func (p byGroupName) Less(i, j int) bool { return p[i].name < p[j].name }
func (p byGroupName) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
