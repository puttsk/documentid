package documentid

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var varRegex *regexp.Regexp
var varIndex int

// ReplaceVars replaces a string containing named variables (`%(var)`) with values from `args`.
func ReplaceVars(s string, args map[string]interface{}) string {
	if varRegex == nil {
		varRegex = regexp.MustCompile(`%\((?P<var>\w+)\)`) // Regex for variable
		varIndex = varRegex.SubexpIndex("var")
	}

	variables := varRegex.FindAllStringSubmatch(s, -1)

	for _, v := range variables {
		val := ""

		if arg, ok := args[v[varIndex]]; ok {
			val = fmt.Sprintf("%v", arg)
		}

		s = strings.Replace(s, v[0], val, 1)
	}

	return s
}

// GenerateDocumentID generates an ID for a document with user-specified prefix, suffix, and zero-padded number.
// The geneated ID will be in `<prefix><number><suffix>` format. The function replaces variables, specified by
// `%(var)`, with values from `args`.
//
// Reserved variable names
// `y`: Year, e.g., 2022
// `m`: Numeric month of the year with zero padding
// `d`: Numeric day of the month with zero padding
// `date`: Date in `YYYYMMDD` format
func GenerateDocumentID(prefix string, suffix string, t time.Time, zeroes int64, number int64, args map[string]interface{}) string {
	if args == nil {
		args = map[string]interface{}{}
	}

	// Override reserved variables
	args["y"] = t.Format("2006")
	args["m"] = t.Format("01")
	args["d"] = t.Format("02")
	args["date"] = t.Format("20060102")

	if prefix != "" {
		prefix = ReplaceVars(prefix, args)
	}
	if suffix != "" {
		suffix = ReplaceVars(suffix, args)
	}

	numStr := fmt.Sprintf("%%0%dd", zeroes)
	return fmt.Sprintf("%s"+numStr+"%s", prefix, number, suffix)
}
