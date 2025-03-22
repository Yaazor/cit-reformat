package main

import (
	"cit-transform/cit"
	"cit-transform/transform"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/magiconair/properties"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var itemsMap map[string][]transform.ConvertElement

func main() {
	color.RGB(156, 225, 235).Println("🔧 Démarrage du programme.")
	itemsMap := make(map[string][]transform.ConvertElement)

	data, err := GetPropertiesFiles("./src")
	if err != nil {
		log.Fatal(err)
	}

	color.RGB(174, 235, 156).Println("🖊️ Création du JSON...")
	for _, d := range *data {
		p, err := properties.LoadString(d)
		if err != nil {
			log.Fatal(err)
		}

		c, err := PropertiesToStruct(*p)
		if err != nil {
			log.Fatal(err)
		}

		conv := StructToOutput(c)
		itemsMap[c.Item] = append(itemsMap[c.Item], *conv)
	}

	js, _ := json.MarshalIndent(itemsMap, "", "\t")
	err = os.WriteFile("output.json", js, 0755)

	if err != nil {
		log.Fatal(err)
	}

	color.RGB(217, 156, 235).Println("🖥️ Conversion terminée.")
}

func StructToOutput(c *cit.SourceCIT) (o *transform.ConvertElement) {
	t := "custom_data"
	if c.CustomData == "none" {
		t = "regex"
	}

	convertStr := strings.ReplaceAll(c.Pattern, "pattern:", "")

	if t == "regex" {
		r := regexp.MustCompile("[^a-zA-Z0-9 -]")
		convertStr = string(r.ReplaceAll([]byte(convertStr), []byte("")))
	}

	o = &transform.ConvertElement{
		Criteria: transform.OutputCriteria{
			Type:  t,
			Match: strings.ReplaceAll(c.Pattern, "pattern:", ""),
		},
		Transform: transform.OutputTransform{
			Type:    "item_name",
			Convert: convertStr,
		},
	}

	return
}

func GetPropertiesFiles(root string) (f *[]string, e error) {
	var files []string
	removedPrefixes := []string{
		"ipattern", "iregex", ": ", ":",
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if strings.Contains(path, "cit") && strings.HasSuffix(info.Name(), ".properties") {
				txt, err := os.ReadFile(path)

				if err != nil {
					fmt.Print(err)
				} else {
					stx := string(txt)
					color.RGB(235, 221, 156).Printf("📄 Lecture du fichier '%v'\n", info.Name())
					if strings.Contains(stx, "components.custom_name=") {
						stx = strings.ReplaceAll(stx, "matchItems", "items")
						stx = strings.ReplaceAll(stx, "items=minecraft:", "items=")

						for r := range removedPrefixes {
							stx = strings.ReplaceAll(stx, removedPrefixes[r], "")
						}

						files = append(files, stx)
					}

				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	f = &files

	return
}

func PropertiesToStruct(bp properties.Properties) (*cit.SourceCIT, error) {
	m := bp.Map()
	convertedMap := make(map[string]string)

	for k, v := range m {
		prefix := k
		suffix := v

		if strings.HasPrefix(k, "components.custom_data.") {
			prefix = strings.ReplaceAll(k, "components.custom_data.", "")
			convertedMap["custom_data"] = prefix + ":" + suffix
		} else {
			convertedMap[prefix] = suffix
		}
	}

	p := properties.LoadMap(convertedMap)
	var dec cit.SourceCIT
	err := p.Decode(&dec)

	if err == nil {
		return &dec, nil
	}

	return nil, fmt.Errorf("erreur dans la conversion")
}
