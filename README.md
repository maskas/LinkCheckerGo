# LinkCheckerGo
Extremely fast recursive link checker to detect broken links on your website. Launches multiple routines at the same time to speed up the process.

## Features
- checks page for errors recursivelly, always staying within the boundaries of the same domain
- lists all urls, that returns anything else than 200
- has a parameter to limit count of URLs to be checked
- supports HTTP and HTTPS
- ignores invalid ssl certificates to let you check on the staging server
- display progress counter

## How to use
- download link-checker.go file
- install go language. Installation depends on your OS. Use Google to find out how ot do this step
- open a terminal window
- navigate to the directory, where you have downloaded the file
- type: `go run link-checker.go http://www.example.com 99` where 99 is the max count of URLs to be checked

## Limitations
- doesn't support relative urls
- doesn't support protocol relative urls, that starts with a double slash
- checks only internal links. Outside links are not checked
- broken images are not checked (only links)
 
