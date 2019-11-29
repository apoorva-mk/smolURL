package main

import (
        "encoding/json"
        //"fmt"
        "net/http"
        "log"
        //"strconv"
)

var urlMapShortLong map[string]string
var urlMapLongShort map[string]string
var base string
var base62 string



func main() {
        base = "http://localhost:8080/smol/"
        base62 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

        urlMapLongShort = make(map[string]string)
        urlMapShortLong = make(map[string]string)

        http.HandleFunc("/smol/", redirectServer)
        http.HandleFunc("/shorten", shortenServer)
        http.ListenAndServe(":8080", nil)
}

func convertBase42(n int) string{
        converted := ""
        for n>0 {
                converted = (string)(base62[n%42])+converted
                n = n/42
        }
        return converted
}


func redirectServer(w http.ResponseWriter, r *http.Request) {
        log.Println("redirecting")
        log.Println(r.URL)
        shortUrl := r.URL.Path[6:]
        longUrl := urlMapShortLong[shortUrl]
        //http.Redirect(w, r, longUrl, 301)
        redirect(longUrl, w, r)
        return

}

func redirect(target string, w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, target, 301)
        return 
}

func getShortUrl(longUrl string) string{
        url, ok := urlMapLongShort[longUrl]

        if ok == false {
                n := len(urlMapLongShort)+1
                url = convertBase42(n)
                for len(url) < 6 {
                        url = "a"+url
                }

                urlMapLongShort[longUrl] = url
                urlMapShortLong[url] = longUrl                 
        }
        return base+url
        
}

func shortenServer(w http.ResponseWriter, r *http.Request) {

        type RequestBody struct {
                LongUrl string
                ShortUrl string
        }

        type ResponseBody struct {
                ShortUrl string
        }

        var req_body RequestBody
    
        err := json.NewDecoder(r.Body).Decode(&req_body)

        if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
        }

        if(req_body.ShortUrl == "") {
                var res_body = ResponseBody{}
                res_body.ShortUrl = getShortUrl(req_body.LongUrl)

                respJSON, err := json.Marshal(res_body)

                if err != nil {
                        http.Error(w, err.Error(), http.StatusBadRequest)
                        return
                }

                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusOK)
                w.Write(respJSON)                        

        }

        // else {

        // }
}