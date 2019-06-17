package exploration

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func FindImports(root string, includeVendor bool) (imports []Import, err error) {
	var paths []string
	importMap := map[string]Import{}

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !includeVendor && strings.HasSuffix(filepath.Dir(path), "vendor") {
			return filepath.SkipDir
		}

		if strings.HasSuffix(path, ".go") {
			paths = append(paths, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		contents, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		indexes := findKeyword("import", contents)

		parsedImports := parseImports(indexes, contents)

		for _, source := range parsedImports {
			i := importMap[source]
			i.Source = source
			i.UsedIn = append(i.UsedIn, path)
			importMap[source] = i
		}
	}

	for _, i := range importMap {
		imports = append(imports, i)
	}

	return imports, nil
}

func findKeyword(word string, bytes []byte) (indexes []int) {
	var inStringLit bool
	var stringDelimiter rune
	word += " "
	wordlen := len(word)

	for index := 0; index+wordlen < len(bytes); index++ {
		wordend := index + wordlen
		b := bytes[index]

		if b == '"' && (stringDelimiter == '"' || stringDelimiter == 0) {
			stringDelimiter = 0
			inStringLit = !inStringLit

			if inStringLit {
				stringDelimiter = '"'
			}
		}

		if b == '`' && (stringDelimiter == '`' || stringDelimiter == 0) {
			stringDelimiter = 0
			inStringLit = !inStringLit

			if inStringLit {
				stringDelimiter = '`'
			}
		}

		if inStringLit {
			continue
		}

		if string(bytes[index:wordend]) == word {
			indexes = append(indexes, index)
			index = wordend
		}
	}

	return indexes
}

func parseImports(indexes []int, contents []byte) (imports []string) {
indexLoop:
	for _, index := range indexes {
		var block bool
		var inImport bool
		var iport []byte

		for ; index < len(contents); index++ {
			if contents[index] == '\n' {
				continue indexLoop
			}

			if contents[index] == '(' {
				block = true
				break
			} else if contents[index] == '"' {
				break
			}
		}

		if block {
			for ; index < len(contents) && contents[index] != ')'; index++ {
				if contents[index] == '"' && contents[index-1] != '\\' {
					inImport = !inImport

					if !inImport {
						imports = append(imports, string(iport))
						iport = []byte{}
					}
				}

				if inImport && contents[index] != '"' {
					iport = append(iport, contents[index])
				}
			}

			continue indexLoop
		}

		index++
		for ; index < len(contents) && contents[index] != '"'; index++ {
			iport = append(iport, contents[index])
		}

		imports = append(imports, string(iport))
	}

	return imports
}
