# better-census-api
### API for easier access to U.S. Census data
#### About the Project
The purpose of Better Census API is to provide improved access to U.S. Census data. The API fetches data from the [official census API](https://api.census.gov),
parses it into a properly formatted JSON response, and includes additional information about datasets and variables.

#### Caveats
Data is currently available only for Detailed Tables of the American Community Survey.

#### Features
The API includes elementary search functions.

The following query returns all datasets with the phrase "american community survey" in the title (note that some tables are titled with the abbreviation **'ACS'**).

http://better-census-api.com/finddataset?search=american%20community%20survey,detailed%20tables&vintage=*

You will use the **Dataset_ID** in the response to search for tables within the dataset [Dataset_ID is **not static**].

http://better-census-api.com/findtable?search=income,household&datasetid=[Dataset_ID]

Once you have identified a data group (i.e. table), you can either pull the whole group or one or more variables (i.e. columns). You will need a personal [Census API key](https://api.census.gov/data/key_signup.html).

http://better-census-api.com/gettable?vintage=2018&dataset=acs5&state=36&county=*&group=B19049&variable=001E,002E&geography=tract&key=32dd72aa5e814e89c669a4664fd31dcfc3df333d

Parameter | Value | * 
----------|-------|---
vintage | data year | N
dataset | acs1, acs3, acs5 | N
state | state FIPS code(s) | Y
county | county FIPS code(s) | Y
group | group code | N
variable | variable code(s) | Y
geography | county, tract, block group | N
key | your census api key

Data is currently available only for American Community Survey detailed tables. ACS1 and ACS3 data is available at county level only. Many variables are available down to block group level for ACS5.
Calls for all of U.S. data ('state=*') may take 30 seconds or more to process.
