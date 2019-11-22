package aliacm

import (
	"sync"
)

// Observer observes the config change.

type UpdateUnit struct {
	Info
	Conf Conf
}

type Conf []byte

func (c Conf) IsEqual(conf Conf) bool {
	if len(c) != len(conf) {
		return false
	}

	for _, cc := range c {
		for _, cconf := range conf {
			if cconf != cc {
				return false
			}
		}
	}
	return true
}

type UpdateUnits []UpdateUnit

// 配置更新完毕后的回调函数
type AfterUpdateHook func(UpdateUnits)

type Observer struct {
	AfterUpdateHook AfterUpdateHook
	confs           sync.Map
	infos           []Info
}

// 用来添加想要关心的配置
func (o *Observer) AddInfo(ufs ...Info) {
	for _, uf := range ufs {
		o.confs.LoadOrStore(uf, nil)
		o.infos = append(o.infos, uf)
	}
}

func (o *Observer) Infos() []Info {
	return o.infos[:]
}

// ACM配置更新后的回调函数
func (o *Observer) OnUpdate(unit Unit, config Config) {
	foundFlag := false
	readFlag := true
	var copyUnits UpdateUnits

	o.confs.Range(func(key, valueIf interface{}) bool {
		flag, ok := key.(Info)
		var conf Conf
		if ok {
			if flag.Group == unit.Group && flag.DataID == unit.DataID {
				conf = config.Content
				o.confs.Store(flag, conf)
				foundFlag = true
				updateUnit := UpdateUnit{Conf: config.Content}
				updateUnit.Group = unit.Group
				updateUnit.DataID = unit.DataID

				copyUnits = append(copyUnits, updateUnit)
			} else {
				if valueIf == nil {
					readFlag = false
				} else {
					value, ok := valueIf.(Conf)
					if ok {
						updateUnit := UpdateUnit{Conf: value}
						updateUnit.Group = flag.Group
						updateUnit.DataID = flag.DataID
						copyUnits = append(copyUnits, updateUnit)
					}
				}
			}
		}
		return true
	})
	if readFlag && foundFlag && o.AfterUpdateHook != nil {
		o.AfterUpdateHook(copyUnits)
	}
}
