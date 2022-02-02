package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/getkin/kin-openapi/openapi3filter"
)

type NitricTarget struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func main() {
	// TODO: Figure out where to get this file
	// Just use a relative openapi file for now, will move to an env variable soon
	wd, _ := os.Getwd()
	router := openapi3filter.NewRouter().WithSwaggerFromFile(fmt.Sprintf("%s/openapi.json", wd))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if route, _, err := router.FindRoute(r.Method, r.URL); err == nil {
			// now we want to pass the request verbatim to the specified nitric target
			var target NitricTarget
			if err := json.Unmarshal(route.Operation.ExtensionProps.Extensions["x-nitric-target"].(json.RawMessage), &target); err != nil {
				// Throw an internal server error here...
				w.WriteHeader(500)
				w.Write([]byte("There was an error reading the provided openapi spec"))
			}

			// target should be all good here...
			targetURL := fmt.Sprintf("http://%s:9001", target.Name)

			// Append the path
			targetURL = fmt.Sprintf("%s%s", targetURL, r.URL.Path)

			if r.URL.RawQuery != "" {
				targetURL = fmt.Sprintf("%s?%s", targetURL, r.URL.RawQuery)
			}

			// pass through the request
			newReq, _ := http.NewRequest(r.Method, targetURL, r.Body)
			newReq.Header = r.Header

			if res, err := http.DefaultClient.Do(newReq); err == nil {
				resBody, _ := ioutil.ReadAll(res.Body)

				for key, val := range res.Header {
					w.Header()[key] = val
				}

				// TODO: Add Proxy header?
				w.WriteHeader(res.StatusCode)
				w.Write(resBody)
			} else {
				w.WriteHeader(404)
				w.Write([]byte("Function not available"))
			}
		} else {
			// There was an error, lets return a not found response for now

			w.WriteHeader(404)
			w.Write([]byte("Function not found"))
		}
	})

	// invoke the target directly based on the extension (validate and pass through the request...)
	fmt.Println("Gateway Listening @ 0.0.0.0:8080")
	err := http.ListenAndServe("0.0.0.0:8080", nil)

	// Log and exit(1)
	log.Fatal(err)
}
