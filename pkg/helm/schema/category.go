package schema

import (
	"fmt"
	"github.com/imdario/mergo"
	"helm.sh/helm/v3/pkg/chart"
	"strconv"
	"strings"
)

func LoadCategories(currCh, parentCh *chart.Chart) map[int]string {
	dest := make(map[int]string, 0)
	if parentCh != nil {
		dest = loadCategories(parentCh)
	}

	curr := loadCategories(currCh)

	err := mergo.Merge(&dest, curr, mergo.WithOverwriteWithEmptyValue)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("merged dest:", dest)
	return dest
}

func loadCategories(ch *chart.Chart) map[int]string {
	var result = make(map[int]string, 0)
	for _, file := range ch.Files {

		if strings.HasPrefix(file.Name, schemaDirName) {
			if strings.HasSuffix(file.Name, jsonSuffix) {
				l := len(file.Name)
				n := file.Name[len(schemaDirName)+1 : l-len(jsonSuffix)]
				frames := strings.Split(n, "-")
				if len(frames) == 2 {
					id, err := strconv.Atoi(frames[0])
					if err != nil {
						continue
					}
					category := frames[1]
					result[id] = category
				}
			}
		}
	}

	return result
}
