package bedrock

type Command struct {
	session *Session
}


func (c *Command) Run(command string) {
}

func (c *Command) Whisper(target string, message string) {
}


func (c *Command) Say(message string) {
}

func (c *Command) Teleport(target string, position Position) {
}

func (c *Command) Give(target string, item string, amount int) {
}

func (c *Command) Effect(target string, effect string) {
}


func (c *Command) execute(command string) {
}