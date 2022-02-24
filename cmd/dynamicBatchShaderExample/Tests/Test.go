package tests

type Test struct {
	M_CurrentTest string
}

type TestMenu struct {
	Menus        []string
	CurrentTests string
}

func SetUpTest() *Test {
	test := &Test{}
	return test
}

func (test *Test) OpenTestMenu(menu *TestMenu) {
	test.M_CurrentTest = menu.CurrentTests
}

func (test *Test) deleteTest() {

}

func (test *Test) OnUpdate(deltaTime float32) {

}
func (test *Test) OnRender() {

}
func (test *Test) OnUIRender() {

}
