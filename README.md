# better-census-api
### API for easier access to U.S. Census data
#### About the Project
The purpose of Better Census API is to provide improved access to U.S. Census data. The Go application fetches data from the [official census API](https://api.census.gov),
parses it into a properly formatted JSON response, and includes additional information about datasets and variables.

#### Caveats
Data is currently available only for Detailed Tables of the American Community Survey.

#### Features
The API includes elementary search functions.

The following query returns all datasets with the phrase "american community survey" in the title (note that some tables are titled with the abbreviation **'ACS'**).

http://better-census-api.com/finddataset?vintage=*&search=american%20community%20survey,detailed%20tables

You will use the **Dataset_ID** in the response to search for tables within the dataset [Dataset_ID is **not static**].

http://better-census-api.com/findtable?search=income,household&datasetid=[Dataset_ID]

Once you have identified a data group (i.e. table), you can either pull the whole group or one or more variables (i.e. columns). You will need a personal [Census API key](https://api.census.gov/data/key_signup.html).

http://better-census-api.com/gettable?vintage=2018&dataset=acs5&state=36&county=*&group=B01001&variable=001E,002E&geography=tract&key=[your_key]


Parameter | Value 
----------|-------
vintage | data year
dataset | acs1/acs3/acs5
state | state FIPS code(s)
county | county FIPS code(s)
group | group code
variable | variable code(s)
geography | county, tract, block group
key | your census api key

Parameters that accept multiple (comma separated) values, also accept an asterisk ('*') for all values. 

Data is currently available only for American Community Survey detailed tables. ACS1 and ACS3 data is available at county level only. Many variables are available down to block group level for ACS5.
Calls for complete U.S. data ('state=*') may take 30 seconds or more to process.
