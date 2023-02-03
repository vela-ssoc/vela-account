package account

import (
	"github.com/vela-ssoc/vela-kit/opcode"
	"github.com/vela-ssoc/vela-kit/vela"
)

func (snap *snapshot) Create(bkt vela.Bucket) {
	for name, item := range snap.current {
		bkt.Store(name, item, 0)
		snap.report.doCreate(item)
		snap.onCreate.Do(&item, snap.co, func(err error) {
			xEnv.Errorf("account snapshot create pipe call fail %v", err)
		})
	}

}

func (snap *snapshot) Update(bkt vela.Bucket) {
	for name, item := range snap.update {
		bkt.Store(name, item, 0)
		snap.report.doUpdate(item)
		snap.onUpdate.Do(&item, snap.co, func(err error) {
			xEnv.Errorf("account snapshot update pipe call fail %v", err)
		})
	}

}

func (snap *snapshot) Delete(bkt vela.Bucket) {
	for name, item := range snap.delete {
		bkt.Delete(name)
		snap.report.doDelete(name)
		snap.onDelete.Do(&item, snap.co, func(err error) {
			xEnv.Errorf("account snapshot delete pipe call fail %v", err)
		})
	}
}

func (snap *snapshot) Report() {
	if snap.enable && snap.report.len() > 0 {
		xEnv.TnlSend(opcode.OpAccountDiff, snap.report)
	}
}
