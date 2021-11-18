package example

import "github.com/coc1961/gomock/example/t1"

type IFace interface{}
type TestInterface interface {
	Func4(string) error
	Func3(str1, str2 string) error
	Func2(arr []int) (IFace, t1.T1Interface, map[string]string, error)
	Func1(i int, pt1 t1.T1) (rt1 *t1.T1, xx t1.T1Interface, err error)
}
