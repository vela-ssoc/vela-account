package account

import (
	"bufio"
	"fmt"
	cond "github.com/vela-ssoc/vela-cond"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"os"
	"strings"
	"sync"
	"time"
)

var colon = ":"

var cache = struct {
	mu   sync.Mutex
	tab  map[uint32]string
	last int64
}{
	tab:  make(map[uint32]string, 32),
	last: 0,
}

func convert(line string, v *Account) bool {
	u := strings.Split(line, colon)
	if len(u) < 7 {
		xEnv.Errorf("not convert %s to linux account", string(line))
		return false
	}

	v.Name = u[0]
	v.UID = u[2]
	v.GID = u[3]
	v.FullName = u[4]
	v.Home = u[5]
	v.Description = u[6]
	v.Status = "OK"

	return true
}

func By(cnd *cond.Cond) ([]Account, error) {

	f, err := os.Open("/etc/passwd")
	if err != nil {
		return nil, fmt.Errorf("read /etc/passwd fail %v", err)
	}
	defer f.Close()

	var av []Account
	add := func(v Account) { av = append(av, v) }

	rd := bufio.NewScanner(f)
	for rd.Scan() {
		v := Account{}
		if !convert(rd.Text(), &v) {
			continue
		}

		if cnd.Match(&v) {
			add(v)
		}

		if e := rd.Err(); e != nil {
			return nil, err
		}
	}

	return av, nil
}

func update() map[uint32]string {
	now := time.Now().Unix()

	if now-cache.last < 2*3600 {
		return cache.tab
	}
	tab := make(map[uint32]string, 32)
	f, err := os.Open("/etc/passwd")
	if err != nil {
		return nil
	}
	defer f.Close()

	add := func(v Account) {
		cache.mu.Lock()
		defer cache.mu.Unlock()
		uid, er := auxlib.ToUint32E(v.UID)
		if er != nil {
			return
		}
		tab[uid] = v.Name
	}

	rd := bufio.NewScanner(f)
	for rd.Scan() {
		v := Account{}
		if !convert(rd.Text(), &v) {
			continue
		}
		add(v)

		if e := rd.Err(); e != nil {
			goto done
		}
	}
done:
	cache.tab = tab
	cache.last = now
	return cache.tab
}

func ByUid(v uint32) string {
	tab := update()
	if tab == nil {
		return ""
	}

	return tab[v]
}
