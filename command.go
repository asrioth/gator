package main

type Command struct {
	Name              string
	Arguments         []string
	NumberOfArguments int
}

type Commands struct {
	Commands     map[string]func(*State, Command) error
	BaseCommands map[string]Command
}

func (C Commands) Register(name string, numArguments int, f func(*State, Command) error) {
	C.Commands[name] = f
	cmd := Command{Name: name, Arguments: make([]string, numArguments), NumberOfArguments: numArguments}
	C.BaseCommands[name] = cmd
}

func (C Commands) Run(s *State, cmd Command) error {
	return C.Commands[cmd.Name](s, cmd)
}
