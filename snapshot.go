package account

import (
	"fmt"
	cond "github.com/vela-ssoc/vela-cond"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/pipe"
	"gopkg.in/tomb.v2"
	"sync/atomic"
	"time"
)

type snapshot struct {
	lua.SuperVelaData
	name     string
	err      error
	bkt      []string
	data     []Account
	onCreate *pipe.Chains
	onDelete *pipe.Chains
	onUpdate *pipe.Chains

	co      *lua.LState
	tomb    *tomb.Tomb
	current map[string]Account
	create  map[string]Account
	delete  map[string]Account
	update  map[string]Account
	enable  bool
	report  *report
}

func newSnapshot() *snapshot {
	sub := atomic.AddUint32(&subscript, 1)

	snap := &snapshot{
		name:     fmt.Sprintf("vela.account.snapshot.%d", sub),
		bkt:      []string{"vela", "account", "snapshot"},
		onCreate: pipe.New(),
		onDelete: pipe.New(),
		onUpdate: pipe.New(),
	}
	snap.V(lua.VTMode, time.Now())
	return snap
}

func (snap *snapshot) init() {
	cnd := &cond.Cond{}
	snap.data, snap.err = By(cnd)
	snap.current = make(map[string]Account, 5)
	snap.create = make(map[string]Account, 5)
	snap.delete = make(map[string]Account, 5)
	snap.update = make(map[string]Account, 5)
	snap.report = &report{}
}

func (snap *snapshot) Start() error {
	return nil
}

func (snap *snapshot) Close() error {
	if snap.tomb != nil {
		snap.tomb.Kill(nil)
	}

	return nil
}

func (snap *snapshot) Name() string {
	return snap.name
}

func (snap *snapshot) Typeof() string {
	return typeof
}

func (snap *snapshot) ok() bool {
	return snap.err == nil
}

func (snap *snapshot) Map() {
	n := len(snap.data)
	if n == 0 {
		return
	}

	for i := 0; i < n; i++ {
		a := snap.data[i]
		snap.current[a.Name] = a
	}
}

func (snap *snapshot) find(name string) (Account, bool) {
	var a Account
	n := len(snap.data)
	if n == 0 {
		return a, false
	}

	for i := 0; i < n; i++ {
		if snap.data[i].Name == name {
			return snap.data[i], true
		}
	}
	return a, false
}

func (snap *snapshot) diff(name string, v interface{}) {
	old, ok := v.(Account)
	if !ok {
		snap.delete[name] = Account{Name: name}
		return
	}

	a, ok := snap.current[name]
	if !ok {
		snap.delete[name] = old
		return
	}
	delete(snap.current, name)

	if a.equal(old) {
		return
	}

	snap.update[name] = a
}

func (snap *snapshot) poll(dt time.Duration) {
	tk := time.NewTicker(dt)
	defer tk.Stop()

	for {
		select {
		case <-snap.tomb.Dying():
			xEnv.Errorf("%s account snapshot over", snap.Name())
			return
		case <-tk.C:
			snap.run()
		}
	}
}

func (snap *snapshot) reset() {
	snap.current = make(map[string]Account, 6)
	snap.create = make(map[string]Account, 6)
	snap.delete = make(map[string]Account, 6)
	snap.update = make(map[string]Account, 6)
	snap.report = nil
}

func (snap *snapshot) run() {
	snap.init()
	if !snap.ok() {
		xEnv.Errorf("init account snapshot fail %v", snap.err)
		return
	}

	snap.Map()
	bkt := xEnv.Bucket(snap.bkt...)
	bkt.Range(snap.diff)
	snap.Create(bkt)
	snap.Update(bkt)
	snap.Delete(bkt)
	snap.Report()
	snap.reset()
}

func (snap *snapshot) sync() {
	snap.init()
	if !snap.ok() {
		xEnv.Errorf("init account snapshot fail %v", snap.err)
		return
	}

	snap.Map()
	bkt := xEnv.Bucket(snap.bkt...)
	bkt.Range(snap.diff)
	snap.Create(bkt)
	snap.Update(bkt)
	snap.Delete(bkt)
	xEnv.Push("/api/v1/broker/collect/agent/account/full", snap.data)
	//xEnv.TnlSend(opcode.OpAccountFull, snap.data)
	snap.reset()
}
