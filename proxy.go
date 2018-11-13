package main

import (
	"net/http"
	"fmt"
	"log"
	"io/ioutil"
	"encoding/json"
	"flag"
	"bytes"
	"time"
	"math/rand"
	"strconv"
)

const (
	webhookKey string = "webhook"
	refKey string = "ref"
	dataKey string = "data"
	callbackKey string = "callback"
	callbackContext string = "c"
)

var baseUrl *string
var proxy map[string]chan(string)
var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	proxy = make(map[string]chan(string))
	baseUrl = flag.String("url", "xxx", "base url to callback")
	port := flag.String("port", "9090", "port to listen")
	flag.Parse()

	http.HandleFunc("/", handler)
	http.HandleFunc("/" + callbackContext, callbackHandler)
	log.Fatal(http.ListenAndServe(":" + *port, nil))
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got a callback!")
	str, dat, err := readJSON(r)
	if err != nil {
		fmt.Fprintf(w, "%q\n", err)
		return
	}
	data := dat["ops"].([]interface{})[0].(map[string]interface{})[dataKey].(map[string]interface{})
	// get channel by ref and write into
	ref := data[refKey].(string)
	ch := proxy[ref]
	ch <- str
	delete(proxy, ref)

	// return response
	fmt.Fprintf(w, "%v\n", "{\"api_result\":\"ok\"}")
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got a request!")
	_, dat, err := readJSON(r)
	if err != nil {
		fmt.Fprintf(w, "%q\n", err)
		return
	}

	// make channel and put to proxy map
	id := strconv.Itoa(seededRand.Int())
	ch := make(chan string)
	proxy[id] = ch

	// put id & callback to request
	data := dat[dataKey].(map[string]interface{})
	data[callbackKey] = *baseUrl + "/" + callbackContext
	data[refKey] = id
	url := dat[webhookKey].(string)

	// send request
	resp, err := sendJSON(url, data)
	if err != nil {
		fmt.Fprintf(w, "%q\n", err)
		return
	}
	fmt.Println("Sync response:", resp)

	// read resp from channel
	response := <-ch
	fmt.Println("Async response:", response)

	// send response
	fmt.Fprintf(w, "%v\n", response)
}

// ------------------------------------------------------------------------------

func readJSON(r *http.Request) (string, map[string]interface{}, error) {
	// read request as string
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", nil, err
	}
	s := string(body)
	fmt.Println("Request:", s)

	// parse json
	var dat map[string]interface{}
	err = json.Unmarshal(body, &dat)
	if err != nil {
		return s, nil, err
	}
	return s, dat, nil
}

func sendJSON(url string, dat map[string]interface{}) (string, error) {
	data, err := json.Marshal(dat)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "invalid response " + resp.Status, nil
	}

	// add channel to proxy
	rand.Seed(time.Now().UnixNano())
	rand.Int()


	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}


