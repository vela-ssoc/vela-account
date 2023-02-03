package account

import (
	"github.com/StackExchange/wmi"
	cond "github.com/vela-ssoc/vela-cond"
	"os/user"
)

type UserAccount struct {
	AccountType uint32
	Caption     string
	Description string
	Disabled    bool
	Domain      string
	FullName    string
	Lockout     bool
	Name        string
	SID         string
	SIDType     uint8
	Status      string
}

func convert(ua []UserAccount, cnd *cond.Cond) []Account {
	var av []Account
	n := len(ua)
	if n == 0 {
		return av
	}

	add := func(a UserAccount) {
		item := Account{
			Name:        a.Name,
			Description: a.Description,
			UID:         a.SID,
			Status:      a.Status,
		}

		uv, er := user.Lookup(a.Name)
		if er == nil {
			item.GID = uv.Gid
			item.Home = uv.HomeDir
		}

		if cnd.Match(&item) {
			av = append(av, item)
		}
	}

	for i := 0; i < n; i++ {
		add(ua[i])
	}

	return av
}

func By(cnd *cond.Cond) ([]Account, error) {
	var dst []UserAccount
	err := wmi.Query("SELECT * FROM Win32_UserAccount where LocalAccount=TRUE", &dst)
	if err != nil {
		return nil, err
	}

	return convert(dst, cnd), nil
}
