package apps

type Directive struct {
	Arg      string
	Callback func()
}
