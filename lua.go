package account

import (
	cond "github.com/vela-ssoc/vela-cond"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
	"os/user"
)

var xEnv vela.Environment

/*
	local v = vela.account.all(cnd)
	local u = vela.account.current()
	local g = vela.group.all(cnd)

	local snap = vela.account.snapshot(true)
	snap.sync()
	snap.on_delete()
	snap.on_create()
	snap.on_update()
	snap.poll(5)
*/

func allL(L *lua.LState) int {
	cnd := cond.CheckMany(L)
	var ret lua.Slice

	data, err := By(cnd)
	if err != nil {
		L.Push(ret)
		return 1
	}

	for _, av := range data {
		if cnd.Match(&av) {
			ret = append(ret, &av)
		}
	}
	L.Push(ret)
	return 1
}

func currentL(L *lua.LState) int {
	u, err := user.Current()
	if err != nil {
		L.RaiseError("got current user fail %v", err)
		return 0
	}

	L.Push(&Account{
		GID:    u.Gid,
		UID:    u.Uid,
		Name:   u.Name,
		Home:   u.HomeDir,
		Status: "OK",
	})

	return 1
}

func snapshotL(L *lua.LState) int {
	snap := newSnapshot()
	snap.co = xEnv.Clone(L)
	snap.enable = L.IsTrue(1)
	proc := L.NewVelaData(snap.Name(), typeof)
	proc.Set(snap)
	L.Push(proc)
	return 1
}

func WithEnv(env vela.Environment) {
	xEnv = env
	kv := lua.NewUserKV()
	kv.Set("all", lua.NewFunction(allL))
	kv.Set("current", lua.NewFunction(currentL))
	kv.Set("snapshot", lua.NewFunction(snapshotL))
	xEnv.Set("account", kv)

	xEnv.Mime(Account{}, encode, decode)
}
