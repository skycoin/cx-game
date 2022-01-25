package tests

type Test struct {
}

type TestMenu struct {
	Menus []string
}

func SetUpTest() *Test {
	test := &Test{}
	return test
}

func (test *Test) deleteTest() {

}

func (test *Test) OnUpdate(deltaTime float32) {

}
func (test *Test) OnRender() {

}
func (test *Test) OnUIRender() {

}
