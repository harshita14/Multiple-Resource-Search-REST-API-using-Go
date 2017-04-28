package main

import(
	"fmt"
    "net/http"
    "io/ioutil"
    "time"
    "strings"
)
func main() {
	//fmt.Println(duckDuckGoSearch("sweets"))
	fmt.Println(googleSearch("the dark knight"))
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

func googleSearch(query string) string{
	query = strings.Replace(query, " ", "", -1)
	timeout := time.Duration(time.Second*5)  //timeout set to one second
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
    posFirst := strings.Index(bodyString, "\"snippet\":")
    //8 is added to remove the "Text": part
    result := bodyString[posFirst+10:]             
    posLast := strings.Index(result, "\",")
    return result[:posLast]
}