package main

import "strings"

func getIdFromUrl(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-2]
}
