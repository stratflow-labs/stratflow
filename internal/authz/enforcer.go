package authz

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

func NewEnforcer() (*casbin.Enforcer, error) {
	m, err := model.NewModelFromFile("pkg/authz/model.conf")
	if err != nil {
		return nil, err
	}

	e, err := casbin.NewEnforcer(m, "pkg/authz/policy.csv")
	if err != nil {
		return nil, err
	}

	return e, nil
}
