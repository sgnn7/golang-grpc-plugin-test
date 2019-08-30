package echoer

const InterfaceName = "Echoer"

type IEcho interface {
	Reply(string) string
}

type Echoer struct {
	F func() IEcho
}

type EchoFunc func() IEcho
