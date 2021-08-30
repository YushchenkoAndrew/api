package helper

func ValidateStr(str1 string, str2 string) (equal bool) {
	var len1 = len(str1)
	var len2 = len(str2)

	var max = len2
	if len1 > len2 {
		max = len1
	}

	equal = true
	for i := 0; i < max; i++ {
		if i > len1 || i > len2 {
			equal = false
		} else if str1[i] != str2[i] {
			equal = false
		}
	}

	return
}
