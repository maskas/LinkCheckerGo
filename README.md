# LinkCheckerGo
Extremely fast recursive link checker to detect broken links and images on your website. Launches multiple routines at the same time to speed up the process.

## Features
- checks page for errors recursivelly, always staying within the boundaries of the same domain
- lists all urls, that return anything else than 200
- has a parameter to limit count of URLs to be checked
- supports HTTP and HTTPS
- ignores invalid ssl certificates to let you check on the staging server
- display progress counter
- returns propper exit status

## How to use
- download link-checker.go file
- install go language. Installation depends on your OS. Use Google to find out how ot do this step
- open a terminal window
- navigate to the directory, where you have downloaded the file
- type: `go run *.go http://www.example.com 99` where 99 is the max count of URLs to be checked
- add false at the end to disable progress 
 
