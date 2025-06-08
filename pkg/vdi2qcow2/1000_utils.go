package vdi2qcow2

import "time"

// setError is a helper method to set error state and message
func (v *VditoQcow2JobStruct) setError(state, errMsg string) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	v.CurrentState = state
	v.ErrorMsg = errMsg
	v.EndTime = time.Now()
}
