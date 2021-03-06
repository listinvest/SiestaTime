//{{{{{{{ Hive Miscelanious Functions and external sources }}}}}}}

//// Extra functions to help Hive with different tasks

//By Rebujacker - Alvaro Folgado Rueda as an open source educative project
package main

import (
	"time"
	"bytes"
	"encoding/json"
	"strings"
	"math/rand"
)

type Log struct {
    Pid  string   `json:"pid"`              // Parent Id: Hive, R-<ID>/B-<ID>
    Time string   `json:"time"`
    Error  string   `json:"error"`
}


//Add a Log to the Jobs to send

func addLog(error string){

	var(
		log Log
		job Job
	)

	//Lock shared Slice
	jobsToHive.mux.Lock()
	defer jobsToHive.mux.Unlock()

	if !strings.Contains(bid,"B-"){
		return
	}

	//Skip logs when there is a job/log overhead
	if len(jobsToHive.Jobs) > 10 {
		return
	}

	time := time.Now().Format("02/01/2006 15:04:05 MST")
	log = Log{bid,time,error}

	bufRP := new(bytes.Buffer)
	json.NewEncoder(bufRP).Encode(log)
	resultRP := bufRP.String()
	param := "["+resultRP+"]"

	//Check that the size of the Result doesn't exceed 20 MB
	bytesLog := len(param)
	if (bytesLog >= 20000000){
		param = "This log was too long...Skipped to avoid crashes"
	}

	job = Job{"","","",bid,"log","","","",param}
	jobsToHive.Jobs = append(jobsToHive.Jobs, &job)
}

func randomString(length int) string{

	charset := "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
  	for i := range b {
    	b[i] = charset[seededRand.Intn(len(charset))]
  	}

  	return string(b)
}