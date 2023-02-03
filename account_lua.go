package account

import (
	cond "github.com/vela-ssoc/vela-cond"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-process"
)

func (a *Account) String() string                         { return lua.B2S(a.Byte()) }
func (a *Account) Type() lua.LValueType                   { return lua.LTObject }
func (a *Account) AssertFloat64() (float64, bool)         { return 0, false }
func (a *Account) AssertString() (string, bool)           { return "", false }
func (a *Account) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (a *Account) Peek() lua.LValue                       { return a }

func (a *Account) ps(L *lua.LState) int {
	cnd := cond.New("username = " + a.Name)
	cnd.CheckMany(L, cond.Seek(0))
	L.Push(process.By(cnd))
	return 1
}

func (a *Account) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "name":
		return lua.S2L(a.Name)
	case "home":
		return lua.S2L(a.Home)
	case "uid":
		return lua.S2L(a.UID)
	case "gid":
		return lua.S2L(a.GID)
	case "status":
		return lua.S2L(a.Status)
	case "ps":
		return lua.NewFunction(a.ps)
	}

	return lua.LNil
}
