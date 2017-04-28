package main

import(
	"fmt"
    "net/http"
    "io/ioutil"
    "time"
    "strings"
    "github.com/julienschmidt/httprouter"
)

func main() {
	r := httprouter.New()
    // Add a handler on /:querystr
    r.GET("/:querystr", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
        query := p.ByName("querystr")
        //using a response channel to be used with goroutines
        response := make(chan string)
        search(query, response)
        resultfirst := <-response
        resultsecond := <-response
        output := "{\"query\": \"%s\",\"results\": {%s,%s}}"
        output = fmt.Sprintf(output, query, resultfirst, resultsecond)
        fmt.Fprint(w, output)
    })
    // the server runs here
    http.ListenAndServe("localhost:3000", r)
}

//function which launches different goroutines for each search
func search(query string, ch chan<-string){
	go duckDuckGoSearch(query, ch)
	go googleSearch(query, ch)
}

func duckDuckGoSearch(query string, ch chan<-string){
	//URL strings to be printed in json response
	duckDuckStringArr := []string{"http://api.duckduckgo.com/?q=", query, "&format=json"}
	duckurl := strings.Join(duckDuckStringArr, "")
	//timeout set to one second
	timeout := time.Duration(time.Second)  
	client := http.Client{
    	Timeout: timeout,
	}
    // Make a get request
	s := []string{"http://api.duckduckgo.com/?q=", query, "&format=json"}
	queryString := strings.Join(s, "")
    rs, err := client.Get(queryString)
    // Process response
    if err != nil {
        panic(err) 
    }
    //close the response at the end of the function
    defer rs.Body.Close()  
    bodyBytes, err := ioutil.ReadAll(rs.Body)
    if err != nil {
        panic(err)
    }
    //bodyString is the json response as a string
    bodyString := string(bodyBytes)

    //slicing bodystring for getting first result
    posFirst := strings.Index(bodyString, "\"Text\":")
    //8 is added to remove the "Text": part
    result := bodyString[posFirst+8:]             
    posLast := strings.Index(result, "\"}")
    result = result[:posLast]
    jsonString := "\"duckduckgo\": {\"url\": \"%s\",\"text\": \"%s\"}"
    ch <- fmt.Sprintf(jsonString, duckurl, result)
}

func googleSearch(query string, ch chan<-string){
	//URL strings to be printed in json response
	googleStringArr := []string{"https://www.googleapis.com/customsearch/v1?q=", query}
	googleurl := strings.Join(googleStringArr, "")
	//timeout set to one second
	timeout := time.Duration(time.Second*5)  
	client := http.Client{
    	Timeout: timeout,
	}
    // Make a get request
	s := []string{"https://www.googleapis.com/customsearch/v1?key=AIzaSyBwkNbaIKhYBqghBwxsNH1ES4Ze_4WphPk&cx=017576662512468239146:omuauf_lfve&q=", query}
	queryString := strings.Join(s, "")
    rs, err := client.Get(queryString)
    // Process response
    if err != nil {
        panic(err) 
    }
    //close the response at the end of the function
    defer rs.Body.Close()  
    bodyBytes, err := ioutil.ReadAll(rs.Body)
    if err != nil {
        panic(err)
    }
    //bodyString is the json response as a string
    bodyString := string(bodyBytes)
    //slicing bodystring for getting first result
    posFirst := strings.Index(bodyString, "\"title\":")
    //8 is added to remove the "Text": part
    result := bodyString[posFirst+10:]             
    posLast := strings.Index(result, "\",")
    result = result[:posLast]
    jsonString := "\"google\": {\"url\": \"%s\",\"text\": \"%s\"}"
    ch <- fmt.Sprintf(jsonString, googleurl, result)
}