package updateData

import (
	"sync"
	"time"
)

var Manager *UpdateDataManager

type UpdateDataManager struct {
	closed	bool
	mutex sync.RWMutex
	data map[uint64]*UpdateData
}

func init(){
	Manager := new(UpdateDataManager)
	Manager.data= make(map[uint64]*UpdateData)
	go Manager.Update()
}

func (this *UpdateDataManager)AddData(d *UpdateData){
	if  nil == d{
		return
	}
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.data[d.uid]= d
}

func (this *UpdateDataManager)DelData(uid uint64)bool{
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if  v,ok := this.data[uid];ok{
		if nil !=v{
			v.Close()
		}
		delete(this.data, uid)
		return true
	}
	return false
}

func (this *UpdateDataManager)GetData(uid uint64)*UpdateData{
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	if  v,ok := this.data[uid];ok{
		return v
	}
	return nil
}

func AddData(data *UpdateData){
	if nil== Manager || nil == data{
		return
	}
	Manager.AddData(data)
}

func DelData(uid uint64)bool{
	if nil== Manager {
		return false
	}
	return Manager.DelData(uid)
}

func GetData(uid uint64) *UpdateData {
	if nil== Manager {
		return nil
	}
	return Manager.GetData(uid)
}

func (this *UpdateDataManager)Close(){
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.closed= true
	for k, v := range this.data{
		if nil !=v{
			v.Close()
		}
		delete(this.data, k)
	}
}

func (this *UpdateDataManager)Update(){
	t := time.Tick(500 * time.Millisecond)
	for _ = range t {
		this.mutex.RLock()
		if true == this.closed{
			this.mutex.RUnlock()
			break
		}
		for _, v := range this.data{
			if nil !=v{
				v.UpdateData()
			}
		}
		this.mutex.RUnlock()
	}
}
