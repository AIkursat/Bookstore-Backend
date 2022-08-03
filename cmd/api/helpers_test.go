package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_readJSON(t *testing.T){
  // create sample json
   sampleJSON := map[string]interface{}{
	"foo": "bar",
   }

   body, _ := json.Marshal(sampleJSON)

   // declare a variable that we can read into
   var decodededJSON struct{
	FOO string `json:"foo"`
   }

   // create a request
   req, err := http.NewRequest("POST", "/", bytes.NewReader(body))
   if err != nil {
	t.Log(err)
   }

   // create a test response recorder
   rr := httptest.NewRecorder() // rr = response recoder
   defer req.Body.Close()
   
   // call readJSON
   err = testApp.readJSON(rr, req, &decodededJSON)
   if err != nil {
	t.Error("failed to decode json", err)
   }



}

func Test_writeJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	payload := jsonResponse{
		Error: false,
		Message: "foo",
	}

	headers := make(http.Header)
	headers.Add("FOO", "BAR")
	err := testApp.writeJSON(rr, http.StatusOK, payload, headers)
	if err != nil {
		t.Errorf("Failed to write JSON: %v", err)
	}

	testApp.environment = "production"
	err = testApp.writeJSON(rr, http.StatusOK, payload, headers)
	if err != nil {
		t.Errorf("Failed to write JSON in production env: %v", err)
	}
	
	testApp.environment = "development"
}