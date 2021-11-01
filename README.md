# URL Finder
*About the project*
This software helps people find out information about links(urls). It has a user interface built with html and its core logic written in golang. 

*Technical information*
Goquery is used to perse the html results of the web page. We check for the following:
- if the web page has a login form
- the number of internal and external links
- the number of working and non working links

*Running the program*
- clone the repository 
- cd into the code directory 
- run **go install**
- run **go tidy**
- run **go run app.go** - this starts a web server running on port 300
- open up your browser and head to **localhost:300** 