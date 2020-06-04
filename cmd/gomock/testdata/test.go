package testdata

import "github.com/coc1961/gomock/cmd/gomock/testdata/t1"

type TestInterface interface {
	Func1(i int, pt1 t1.T1) (rt1 *t1.T1, xx t1.T1Interface, err error)
}
