package testdata

type IFace interface{}
type TestInterface interface {
	Func2(arr []int) (IFace, t1.T1Interface, map[string]string, error)
	Func1(i int, pt1 t1.T1) (rt1 *t1.T1, xx t1.T1Interface, err error)
}
