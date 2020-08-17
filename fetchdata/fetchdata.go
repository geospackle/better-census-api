package fetchdata

import (
	"census-api/helpers"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"net/url"
    "time"
)

var allStates [50]string = [50]string{"01", "02", "04", "05", "06", "08", "09", "10", "12", "13", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31", "32", "33", "34", "35", "36", "37", "38", "39", "40", "41", "42", "44", "45", "46", "47", "48", "49", "50", "51", "53", "54", "55", "56"}

type Datasets struct {
	ConformsTo string
	Dataset    []DatasetInfo
}

type DatasetInfo struct {
	C_vintage    int
	Title        string
	C_groupsLink string
	Distribution []DistributionInfo
}

type DistributionInfo struct {
	AccessURL string
}

type foundDatasets struct {
	ID        int `json:"Dataset_ID"`
	C_vintage int `json:"Vintage"`
	Title     string
}

type tables struct {
	Groups []groupInfo
}

type groupInfo struct {
	Name        string
	Description string
	Variables   string
}

type groups struct {
	Dataset     string
	AccessURL   string
    Groups      []map[string]string
}

type AllGroups struct {
   Groups []group
   }

type group struct {
    Description string
    Name        string
    Variables   string 
}   

type varProperties struct {
	varType, label string
}

func makeGeoID(varNames []string, varValues []string) string {
	var geoID string
	geoNames := [5]string{"state", "county", "tract", "blockgroup", "block"}
	for _, v := range geoNames {
		geoIdx, ok := helpers.StringIndex(varNames, v)
		if ok {
			geoID += varValues[geoIdx]
		}
	}
	return geoID
}

func FindDataset(datasets []DatasetInfo, vintage string, searchStr string) ([]foundDatasets, error) {
	var pickedValues []foundDatasets
	errors := []error{errors.New("vintage not available"), errors.New("search term not found")}
	if searchStr == "ACS" {
        searchStr = "American Community Survey, ACS"
    }
    for i := range datasets {
		//fmt.Printf("%v", table[i].Title)
		dataVintage := strconv.Itoa(datasets[i].C_vintage)
		if vintage == "*" || vintage == dataVintage {
			errors[0] = nil
			if searchStr == "*" || helpers.Match(searchStr, datasets[i].Title) {
				errors[1] = nil
				pickedValues = append(
					pickedValues, foundDatasets{i, datasets[i].C_vintage, datasets[i].Title})
			}
		}
	}
	err := helpers.GetError(errors)
	return pickedValues, err
}

func FindTable(datasets []DatasetInfo, datasetID int, searchStr string) (groups, error) {
	var descr = make([]map[string]string, 0)
	var err error
	if len(datasets) < datasetID {
		return groups{}, errors.New("dataset ID does not exist")
	}
	tableURL := datasets[datasetID].C_groupsLink
	accessURL := datasets[datasetID].Distribution[0].AccessURL
	dataset := datasets[datasetID].Title
	allTables := new(tables)
	helpers.GetJSON(tableURL, allTables)
    fmt.Println(tableURL, allTables)
	err = errors.New("search term not found")
	for _, v := range allTables.Groups {
		if searchStr == "*" || helpers.Match(searchStr, v.Description) {
			err = nil
            m := make(map[string]string)
            m[v.Name] = v.Description
			descr = append(descr, m)
		}
	}
    out := groups{dataset, accessURL, descr}
	return out, err
}

func mapVar(slc []string, prefix string, conn string) []string {
	out := make([]string, len(slc))
	for i, v := range slc {
		out[i] = prefix + conn + v
	}
	return out
}

func getGroups(vintage int, dataset string) *AllGroups {
    url := fmt.Sprintf("https://api.census.gov/data/%d/acs/%s/groups.json", vintage, dataset)
    allGroups := new(AllGroups)
    helpers.GetJSON(url, allGroups)
    return allGroups
}


func getData(key string, vintage int, dataset string, group string, variable string, geography string, state string, county string) ([]byte, int) {
	var cenVar string
	if variable == "*" {
		cenVar = fmt.Sprintf("group(%s)", group)
	} else {
		vars := strings.Split(variable, ",")
		newVars := mapVar(vars, group, "_")
		cenVar = strings.Join(newVars, ",")
	}
	var myClient = &http.Client{Timeout: 10 * time.Second}
	var URL string
    geography = url.QueryEscape(geography)
	if county == "*" {
		//&in=tract: gives specific tracts
		URL = fmt.Sprintf("https://api.census.gov/data/%d/acs/%s?get=%s&for=%s:*&in=state:%s&key=%s", vintage, dataset, cenVar, geography, state, key)
	} else {
		URL = fmt.Sprintf("https://api.census.gov/data/%d/acs/%s?get=%s&for=%s:*&in=state:%s&in=county:%s&key=%s", vintage, dataset, cenVar, geography, state, county, key)
	}
    r, err := myClient.Get(URL)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	return body, r.StatusCode
}

//gives tract data at state level, block group data at county
func GetTable(key string, vintage int, dataset string, group string, variable string,
	geography string, state string, county string) ([]byte, int) {
	// get census data
	var tableHeader []string
	var data []byte
	var tableBody [][]string
	var table [][]string
	var statusCode int
	states := strings.Split(state, ",")
	if states[0] == "*" {
		county = "*"
		for _, state := range allStates {
			data, statusCode = getData(key, vintage, dataset, group, variable, geography, state, county)
			target := make([][]string, len(data)) //needs to be created each time for copy
			err := json.Unmarshal(data, &target)
			if err != nil {
				return data, statusCode
			} //return error message from census api
			body := make([][]string, len(target[1:]))
			copy(body, target[1:])
			tableBody = append(tableBody, body...)
			tableHeader = target[0]
		}
	} else if len(states) > 1 {
		county = "*"
		for _, state := range states {
            state = strings.TrimSpace(state)
			data, statusCode = getData(key, vintage, dataset, group, variable, geography, state, county)
			target := make([][]string, len(data)) //needs to be created each time for copy
			err := json.Unmarshal(data, &target)
			if err != nil {
				return data, statusCode
			} //return error message from census api
			body := make([][]string, len(target[1:]))
			copy(body, target[1:])
			tableBody = append(tableBody, body...)
			tableHeader = target[0]
		}
	} else {
		counties := strings.Split(county, ",")
		for _, county := range counties {
            county = strings.TrimSpace(county)
			data, statusCode = getData(key, vintage, dataset, group, variable, geography, state, county)
            target := make([][]string, len(data)) //needs to be created each time for copy
			err := json.Unmarshal(data, &target)
			if err != nil {
				return data, statusCode
			} // return error message from census API
			body := make([][]string, len(target[1:]))
			copy(body, target[1:])
			tableBody = append(tableBody, body...)
			tableHeader = target[0]
		}
	}
	var varNames []string
	var varTypes []string
	allVars := getAllVars(vintage, dataset, group)
	for _, code := range tableHeader {
		varProps := getVarProperties(allVars, code)
		varName := varProps.label
		varType := varProps.varType
		varNames = append(varNames, varName)
		varTypes = append(varTypes, varType)
	}
	table = tableBody
	table = append([][]string{varTypes}, table...)
	table = append([][]string{varNames}, table...)
	table = append([][]string{tableHeader}, table...)
	res := tableToJSON(table, vintage, dataset, group)
	return res, 200
}

//types used in tableToJSON
type pair struct {
	Index int
	Value string
}

type keyvalue map[string]interface{}

type varDef struct {
	VarName string `json:"name"`
	VarType string `json:"type"`
} 

type groupDef struct {
    Code string `json:"code"`
    Vintage int `json:"vintage"` 
    Description string `json:"description"`
}

type censusData struct {
    Group     keyvalue `json:"group"`
	Variables keyvalue `json:"variables"`
	GeoID     keyvalue `json:"geoid"`
}

func tableToJSON(table [][]string, vintage int, dataset string, group string) []byte {
	lRows := len(table) - 3            //3 header rows
	dataslc := make([]keyvalue, lRows) // or make []keyvalue and use append
	exp := `[A-Z]+\d.+`
	vars := helpers.RegSlice(exp, table[0])
	for j := 0; j < lRows; j++ {
		dataslc[j] = make(keyvalue) //needs to  initialize
		for i := range vars {
			dataslc[j][vars[i].Value] = table[j+3][vars[i].Index] //can't create an intermediate variable and overwrite, bco referencing
		}
	}
	withid := make(keyvalue, lRows)
	for i := range dataslc {
		geoid := makeGeoID(table[0], table[i+3])
		withid[geoid] = dataslc[i] 
	}
	variables := make(keyvalue, len(vars))
	for i, v := range vars {
		vCode := v.Value
		vType := table[2][vars[i].Index]
		vName := table[1][vars[i].Index]
		variables[vCode] = varDef{VarName: vName, VarType: vType}
	}
    var groupDescr string
    allGroups := getGroups(vintage, dataset)
    for _,v := range(allGroups.Groups) {
        if v.Name == group {
            groupDescr = v.Description
        }
    }
    groupInfo := make(keyvalue, 1)
    groupInfo["description"] = groupDescr
    groupInfo["vintage"] = vintage
    groupInfo["code"] = group
    out := censusData{Group: groupInfo, GeoID: withid, Variables: variables}
	res, _ := json.MarshalIndent(out, "", "    ")
	return res
}

func getAllVars(vintage int, dataset string, group string) interface{} {
	url := fmt.Sprintf("https://api.census.gov/data/%d/acs/%s/groups/%s.json", vintage, dataset, group)
	var target interface{}
	allVars := helpers.GetJSON(url, target)
	return allVars
}

func getVarProperties(allVars interface{}, varCode string) varProperties {
	var varType string
	var label string
	//check if varCode is another type of attribute, e.g. geoid, tract
	checkCode := allVars.(map[string]interface{})["variables"].(map[string]interface{})[varCode]
	if checkCode == nil {
		varType = "string"
		label = varCode
	} else {
		varType = allVars.(map[string]interface{})["variables"].(map[string]interface{})[varCode].(map[string]interface{})["predicateType"].(string)
		label = allVars.(map[string]interface{})["variables"].(map[string]interface{})[varCode].(map[string]interface{})["label"].(string)
	}
	properties := varProperties{varType, label}
	return properties
}
