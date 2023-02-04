package unpack

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Payload struct {
	Int             int32   `json:"int"`
	Uint            uint32  `json:"uint"`
	Short           int16   `json:"short"`
	Float           float64 `json:"float"`
	Double          float64 `json:"double"`
	BigEndianDouble float64 `json:"big_endian_double"`
}

type ProblemJSON struct {
	Bytes string `json:"bytes"`
}

func Unpack(token string) {
	getUrl := fmt.Sprintf("https://hackattic.com/challenges/help_me_unpack/problem?access_token=%s", token)
	res, err := http.Get(getUrl)
	if err != nil {
		fmt.Println("unable to get problem data " + err.Error())
		return
	}

	if res.StatusCode != 200 {
		fmt.Println("not 200 response")
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("cannot read problem json " + err.Error())
		return
	}

	var resJSON ProblemJSON
	err = json.Unmarshal(body, &resJSON)
	decodedData, err := base64.StdEncoding.DecodeString(resJSON.Bytes)
	if err != nil {
		fmt.Println("unable to decode and read base64 data " + err.Error())
		return
	}

	buffer := bytes.NewBuffer(decodedData)

	iBytesBuf := bytes.NewBuffer(buffer.Next(4))
	var i int32
	err = binary.Read(iBytesBuf, binary.LittleEndian, &i)
	if err != nil {
		fmt.Println("not able to read interger " + err.Error())
		return
	}

	iBytesBuf = bytes.NewBuffer(buffer.Next(4))
	var ui uint32
	err = binary.Read(iBytesBuf, binary.LittleEndian, &ui)
	if err != nil {
		fmt.Println("not able to read unsigned interger " + err.Error())
		return
	}

	iBytesBuf = bytes.NewBuffer(buffer.Next(4))
	var s int16
	err = binary.Read(iBytesBuf, binary.LittleEndian, &s)
	if err != nil {
		fmt.Println("not able to read interger " + err.Error())
		return
	}

	iBytesBuf = bytes.NewBuffer(buffer.Next(4))
	var f float32
	err = binary.Read(iBytesBuf, binary.LittleEndian, &f)
	if err != nil {
		fmt.Println("not able to read interger " + err.Error())
		return
	}

	iBytesBuf = bytes.NewBuffer(buffer.Next(8))
	var d float64
	err = binary.Read(iBytesBuf, binary.LittleEndian, &d)
	if err != nil {
		fmt.Println("not able to read interger " + err.Error())
		return
	}

	iBytesBuf = bytes.NewBuffer(buffer.Next(8))
	var db float64
	err = binary.Read(iBytesBuf, binary.BigEndian, &db)
	if err != nil {
		fmt.Println("not able to read interger " + err.Error())
		return
	}

	payload := Payload{
		Int:             i,
		Uint:            ui,
		Short:           s,
		Float:           float64(f),
		Double:          d,
		BigEndianDouble: db,
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("unable to marshal post request payload " + err.Error())
		return
	}

	res, err = http.Post(fmt.Sprintf("https://hackattic.com/challenges/help_me_unpack/solve?access_token=%s", token), "application/json", bytes.NewBuffer(payloadJSON))
	if err != nil {
		fmt.Println("error submitting the POST request " + err.Error())
		return
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("unable to decode and read base64 data " + err.Error())
		return
	}

	var result any
	err = json.Unmarshal(data, &result)
	if res.StatusCode != 200 {
		fmt.Println("failure", result)
		return
	}
	fmt.Println(result)
}
