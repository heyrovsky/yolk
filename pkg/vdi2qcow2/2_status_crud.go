package vdi2qcow2

import "time"

func (v *VditoQcow2JobStruct) updateState(state string) {
	v.mutex.Lock()
	v.CurrentState = state
	v.mutex.Unlock()
}

func (v *VditoQcow2JobStruct) setProgress(p int) {
	v.mutex.Lock()
	v.Progress = p
	v.mutex.Unlock()
}

func (v *VditoQcow2JobStruct) fail(state string, err error) error {
	v.mutex.Lock()
	v.CurrentState = state
	v.ErrorMsg = err.Error()
	v.EndTime = time.Now()
	v.mutex.Unlock()
	return err
}

func (v *VditoQcow2JobStruct) GetProgress() int {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	return v.Progress
}

func (v *VditoQcow2JobStruct) GetState() string {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	return v.CurrentState
}

func (v *VditoQcow2JobStruct) IsDone() bool {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	return v.IsCompleted
}
