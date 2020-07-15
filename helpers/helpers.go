package helpers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func Map(in []string, f func(string) string) []string {
	out := make([]string, len(in))
	for i, v := range in {
		out[i] = f(v)
	}
	return out
}

func Match(strToMatch string, matchedTo string) bool {
	strToMatch = strings.ToUpper(strToMatch)
	matchedTo = strings.ToUpper(matchedTo)
	splitStr := strings.Split(strToMatch, ",")
	for _, str := range splitStr {
		var exp = regexp.MustCompile(str)
		match := exp.MatchString(matchedTo)
		if match == false {
			return false
		}
	}
	return true
}

func GetJSON(url string, target interface{}) interface{} {
	var myClient = &http.Client{Timeout: 10 * time.Second}
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	//    return json.NewDecoder(r.Body).Decode(target)    //unsafe
	json.Unmarshal(body, &target)
	return target
}

func GetHTML(url string) string {
	r, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	return string(body)
}

type pair struct {
	Index int
	Value string
}

// returns a slice of index, value pairs
func RegSlice(exp string, vars []string) []pair {
	var out []pair
	for i, str := range vars {
		match, _ := regexp.MatchString(exp, str)
		if match == true {
			var p = pair{i, str}
			out = append(out, p)
		}
	}
	return out
}

func GetError(errors []error) error {
	var err error
	for i := range errors {
		if errors[i] != nil {
			err = errors[i]
			break
		}
	}
	return err
}

func IndexExists(slc []string, idx int) bool {
	return len(slc) > idx
}

func StringIndex(slc []string, str string) (int, bool) {
	for i, v := range slc {
		if str == v {
			return i, true
		}
	}
	return -1, false
}
