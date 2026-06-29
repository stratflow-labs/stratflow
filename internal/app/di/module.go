package di

import "github.com/samber/do"

type ModuleRegistrar func(inj *do.Injector)

var modules []ModuleRegistrar

func RegisterModule(fn ModuleRegistrar) {
	if fn == nil {
		return
	}
	modules = append(modules, fn)
}
func Register(inj *do.Injector) {
	for _, fn := range modules {
		fn(inj)
	}
}
