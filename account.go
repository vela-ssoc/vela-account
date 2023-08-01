package account

import (
	"github.com/vela-ssoc/vela-kit/kind"
)

type Account struct {
	FullName    string `json:"full_name"`
	Name        string `json:"name"`
	UID         string `json:"uid"`
	GID         string `json:"gid"`
	Home        string `json:"home_dir"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Action      string `json:"action"`
}

func (a *Account) Byte() []byte {
	enc := kind.NewJsonEncoder()
	enc.Tab("")
	enc.KV("minion_id", xEnv.ID())
	enc.KV("minion_inet", xEnv.Inet())
	enc.KV("name", a.Name)
	enc.KV("uid", a.UID)
	enc.KV("gid", a.GID)
	enc.KV("home_dir", a.Home)
	enc.KV("description", a.Description)
	enc.KV("status", a.Status)
	enc.KV("action", a.Action)
	enc.End("}")
	return enc.Bytes()
}

func (a *Account) equal(old Account) bool {
	if a.Name != old.Name {
		return false
	}

	if a.UID != old.UID {
		return false
	}

	if a.Home != old.Home {
		return false
	}

	if a.Description != old.Description {
		return false
	}

	if a.Status != a.Status {
		return false
	}

	return true
}
