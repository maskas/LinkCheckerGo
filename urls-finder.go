package main

import "regexp"

func findUrls(html string) []string {
	var urls = []string{}

	re := regexp.MustCompile("<a .*href=\"([^\"]*)")
	matches := re.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		urls = append(urls, match[1])
	}

	re = regexp.MustCompile("<img .*src=\"([^\"]*)")
	matches = re.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		urls = append(urls, match[1])
	}

	re = regexp.MustCompile("<link .*href=\"([^\"]*)")
	matches = re.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		urls = append(urls, match[1])
	}

	re = regexp.MustCompile("<script .*src=\"([^\"]*)")
	matches = re.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		urls = append(urls, match[1])
	}

	return urls
}