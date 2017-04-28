package main

import(
	"fmt"
    "net/http"
    "io/ioutil"
    "time"
    "strings"
)
func main() {
	fmt.Println(duckDuckGoSearch("dog"))
}

func duckDuckGoSearch(query string) string{
	timeout := time.Duration(time.Second)  //timeout set to one second
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
    return result[:posLast]
}