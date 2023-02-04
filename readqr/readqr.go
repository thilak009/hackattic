package readqr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"io/ioutil"
	"net/http"

	"github.com/liyue201/goqr"
)

type ProblemJSON struct {
	ImageUrl string `json:"image_url"`
}

type ReqPayload struct {
	Code string `json:"code"`
}

func ReadQR(token string) {
	problemUrl := fmt.Sprintf("https://hackattic.com/challenges/reading_qr/problem?access_token=%s", token)
	res, err := http.Get(problemUrl)

	if err != nil {
		fmt.Println("unable to fetch problem " + err.Error())
		return
	}
	if res.StatusCode != 200 {
		fmt.Println("invalid response code " + fmt.Sprint(res.StatusCode))
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("cannot read response body" + err.Error())
		return
	}

	var resJson ProblemJSON
	err = json.Unmarshal(body, &resJson)
	if err != nil {
		fmt.Println("JSON unmarshilling failed " + err.Error())
		return
	}

	res, err = http.Get(resJson.ImageUrl)
	if err != nil {
		fmt.Println("unable to fetch the QR image ", err.Error())
		return
	}

	if res.StatusCode != 200 {
		fmt.Println("invalid response code ", res.StatusCode)
		return
	}

	img, _, err := image.Decode(res.Body)
	if err != nil {
		fmt.Printf("image.Decode error: %v\n", err)
		return
	}

	qrCodes, err := goqr.Recognize(img)
	if err != nil {
		fmt.Printf("Recognize failed: %v\n", err)
		return
	}

	if len(qrCodes) > 1 {
		fmt.Println("found more than one QR code in the image")
		return
	}

	payload := ReqPayload{
		Code: string(qrCodes[0].Payload),
	}
	payloadJSON, err := json.Marshal(payload)

	res, err = http.Post(fmt.Sprintf("https://hackattic.com/challenges/reading_qr/solve?access_token=%s", token), "application/json", bytes.NewBuffer(payloadJSON))
	if err != nil {
		fmt.Println("error submitting solution ", err.Error())
		return
	}

	resBody, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(resBody, &resJson)
	if res.StatusCode != 200 {
		fmt.Println("incorrect solution ", resJson)
		return
	}
	fmt.Println("solution success response ", resJson)
}
