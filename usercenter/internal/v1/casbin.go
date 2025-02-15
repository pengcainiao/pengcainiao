package v1

import (
	casbin "github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/google/martian/log"
	"go.uber.org/zap"
	"pp/usercenter/internal/auth/adapter"
	"sync"
)

var (
	enforcer     *casbin.Enforcer
	enforcerLock = &sync.Mutex{}
)

func SetUp() {
	text :=
		`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`
	m, _ := model.NewModelFromString(text)
	var err error

	a := adapter.NewMysqlAdapter()
	enforcer, err = casbin.NewEnforcer(m, a)
	if err != nil {
		log.Errorf("AnchorSystemExtObj failed", zap.Error(err))
	}
}

func getAllPermsByRole(role string) [][]string {
	perms := enforcer.GetFilteredNamedPolicy("p", 0, role, "", "", "")
	return perms
}

func EnforceList(level string, path string) (bool, error) {
	if pass, _ := enforce(level, path, path); pass {
		return true, nil
	}
	return false, nil
}

func enforce(params ...interface{}) (bool, error) {
	enforcerLock.Lock()
	defer enforcerLock.Unlock()
	return enforcer.Enforce(params...)
}
