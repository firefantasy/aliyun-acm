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
type AfterUpdate func(UpdateUnits)

type Observer struct {
	AfterUpdate AfterUpdate
	confs       sync.Map
}

// 用来添加想要关心的配置
func (o *Observer) AddUnit(uf Info) {
	o.confs.LoadOrStore(uf, nil)
}

// ACM配置更新后的回调函数
func (o *Observer) Modify(unit Unit, config Config) {
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
						updateUnit.Group = unit.Group
						updateUnit.DataID = unit.DataID
						copyUnits = append(copyUnits, updateUnit)
					}
				}
			}
		}
		return true
	})

	if readFlag && foundFlag && o.AfterUpdate != nil {
		o.AfterUpdate(copyUnits)
	}
}
