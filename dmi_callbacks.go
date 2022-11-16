// Code generated by "callbackgen -type DMI"; DO NOT EDIT.

package indicator

import ()

func (inc *DMI) OnUpdate(cb func(diplus float64, diminus float64, adx float64)) {
	inc.updateCallbacks = append(inc.updateCallbacks, cb)
}

func (inc *DMI) EmitUpdate(diplus float64, diminus float64, adx float64) {
	for _, cb := range inc.updateCallbacks {
		cb(diplus, diminus, adx)
	}
}
