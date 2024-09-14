package main

import "os"

// If file exists and is not dir -> true.
// If file not exists or it's a dir -> false.
func FileExists(path string) bool {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) || stat.IsDir() {
		return false
	}
	return true
}
