# LinkCheckerGo
Extremely fast recursive link checker to detect broken links on your website. Launches multiple rutines at the same time to speed up the process.

## How to use
- download link-checker.go file and place it somewhere
- install go language. Installation depends on your OS. Use Google to find out how ot do this step
- open a terminal window
- navigate to the directtory, where you have downloaded the file
- type: go run link-checker.go http://www.example.com 99 where 99 is the max count of URLs to be checked

## Limitations
- checks only internal links. Outside links are not checked
- broken images are not checked (onlyl links)
