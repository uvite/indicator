// Code generated by "callbackgen -type Drift"; DO NOT EDIT.

package indicator

import ()

func (inc *Drift) OnUpdate(cb func(value float64)) {
	inc.UpdateCallbacks = append(inc.UpdateCallbacks, cb)
}

func (inc *Drift) EmitUpdate(value float64) {
	for _, cb := range inc.UpdateCallbacks {
		cb(value)
	}
}