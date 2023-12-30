package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/BishopFox/jsluice"
)

func getJsluiceEndpoints(data string, source string) map[string][]string {
	analyzer := jsluice.NewAnalyzer([]byte(data))
	endpoints := make(map[string][]string)
	foundUrls := analyzer.GetURLs()
	for _, url := range foundUrls {
		if endpoints[string(url.URL)] == nil {
			endpoints[url.URL] = append(endpoints[url.URL], source)
		} else {
			for _, oldSource := range endpoints[url.URL] {
				if source != oldSource {
					endpoints[url.URL] = append(endpoints[url.URL], source)
					break
				}
			}
		}
	}
	return endpoints
}

func getRegexEndpoints(data string, source string, regex string) map[string][]string {
	toTrim := []string{"\"", " ", "\n"}
	endpoints := make(map[string][]string)
	re := regexp.MustCompile(regex)
	matches := re.FindAllString(data, -1)
	for _, endpoint := range matches {
		for _, char := range toTrim {
			endpoint = strings.Trim(endpoint, char)
		}
		if endpoints[endpoint] == nil {
			endpoints[endpoint] = append(endpoints[endpoint], source)
		} else {
			for _, oldSource := range endpoints[endpoint] {
				if source != oldSource {
					endpoints[endpoint] = append(endpoints[endpoint], source)
					break
				}
			}
		}
	}
	return endpoints
}

func loadJsFilesNames(folder string) map[int](string) {
	index := 0
	list := make(map[int](string))
	filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal("couldn't get list of js files: ", err)
		}
		if index >= 1 {
			list[index] = info.Name()
		}
		index += 1
		return nil
	})

	return list
}

func getFile(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("can't load ", filePath, " : ", err)
	}
	return string(data)
}

func appendUniq(allEndpoints map[string][]string, newEndpoints map[string][]string, endpoint string, sources []string) map[string][]string {
	if allEndpoints[endpoint] != nil {
		for _, source := range newEndpoints[endpoint] {
			allEndpoints[endpoint] = append(allEndpoints[endpoint], source)
		}
	} else {
		allEndpoints[endpoint] = sources
	}
	return allEndpoints
}

func main() {
	jsFolder := "./JSFiles"
	jsFiles := loadJsFilesNames(jsFolder)
	basicRegex := `(?:^|\"|'|\\n|\\r|\n|\r|\s)(((?:[a-zA-Z]{1,10}:\/\/|\/\/)([^\"'\/\s]{1,255}\.[a-zA-Z]{2,24}|localhost)[^\"'\n\s]{0,255})|((?:\/|\.\.\/|\.\/)[^\"'><,;| ()(%%$^\/\\\[\]][^\"'><,;|()\s]{1,255})|([a-zA-Z0-9_\-\/]{1,}\/[a-zA-Z0-9_\-\/]{1,255}\.(?:[a-zA-Z]{1,4})(?:[\?|\/][^\"|']{0,}|))|([a-zA-Z0-9_\-]{1,255}\.(?:php|php3|php5|asp|aspx|ashx|cfm|cgi|pl|jsp|jspx|json|js|action|html|xhtml|htm|bak|do|txt|wsdl|wadl|xml|xls|xlsx|bin|conf|config|bz2|bzip2|gzip|tar\.gz|tgz|log|src|zip|js\.map)(?:\?[^\"|^']{0,255}|)))(?:\"|'|\\n|\\r|\n|\r|\s|$)|^Disallow:\s([^\$\n]+)|^Allow:\s([^\$\n]+)| Domain\=([^\\";']+)|\<(https?:\/\/[^>\n]+)|(\"|\')([A-Za-z0-9_-]+\/)+[A-Za-z0-9_-]+(\.[A-Za-z0-9]{2,}|\/?(\?|\#)[A-Za-z0-9_\-&=\[\]])(\"|\')`
	customRegex := ``
	fileTypesToFilter := ".svg"

	allEndpoints := make(map[string][]string)

	for _, fileName := range jsFiles {
		fmt.Println("Loading ", fileName)
		filePath := jsFolder + "/" + fileName
		fileData := getFile(filePath)

		// Jsluice
		jsluiceEndpoints := getJsluiceEndpoints(fileData, fileName)
		for endpoint, sources := range jsluiceEndpoints {
			allEndpoints = appendUniq(allEndpoints, jsluiceEndpoints, endpoint, sources)
		}

		// Custom regex
		customRegexEndpoints := getRegexEndpoints(fileData, fileName, customRegex)
		for endpoint, sources := range customRegexEndpoints {
			allEndpoints = appendUniq(allEndpoints, jsluiceEndpoints, endpoint, sources)
		}

		// Classic Regex
		regexEndpoints := getRegexEndpoints(fileData, fileName, basicRegex)
		for endpoint, sources := range regexEndpoints {
			allEndpoints = appendUniq(allEndpoints, jsluiceEndpoints, endpoint, sources)
		}
	}
	urls := make([]string, 0, len(allEndpoints))
	for k := range allEndpoints {
		urls = append(urls, k)
	}
	sort.Strings(urls)
	for _, url := range urls {
		if !strings.Contains(url, fileTypesToFilter) {
			if !strings.Contains(url, "./") {
				fmt.Println(url)
			}
		}
	}
	// fmt.Println(allEndpoints)
}
