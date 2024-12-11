package functions

import (
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var accentsMap = func() map[rune]string {
	accents := []string{
		"[aàáâã]",
		"[cç]",
		"[eéê]",
		"[ií]",
		"[oóôõ]",
		"[uú]",
	}
	res := make(map[rune]string)
	for _, str := range accents {
		for _, char := range str[1 : len(str)-1] {
			res[char] = str
		}
	}
	return res
}()

func ToCaseInsensitiveRegex(arr []string) bson.M {
	var regex string
	for i := range arr {
		if i != 0 {
			regex += "|"
		}
		escaped := regexp.QuoteMeta(strings.ToLower(arr[i]))
		for _, char := range escaped {
			accs, in := accentsMap[char]
			if in {
				regex += accs
			} else {
				regex += string(char)
			}
		}
	}
	return bson.M{"$regex": primitive.Regex{Pattern: regex, Options: "i"}}
}
