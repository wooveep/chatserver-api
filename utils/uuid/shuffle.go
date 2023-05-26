/*
 * @Author: cloudyi.li
 * @Date: 2023-05-19 10:37:15
 * @LastEditTime: 2023-05-26 16:46:08
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/utils/uuid/shuffle.go
 */
package uuid

func shuffleString(str string) string {
	lastChar := str[len(str)-1]
	pbox := getPbox(lastChar)

	runes := []rune(str)
	shuffled := make([]rune, len(runes))
	for i, newIndex := range pbox {
		shuffled[i] = runes[newIndex]
	}

	return string(shuffled)
}

func unshuffleString(str string) string {
	lastChar := str[len(str)-1]
	pbox := getPbox(lastChar)

	runes := []rune(str)
	unshuffled := make([]rune, len(runes))
	for i, newIndex := range pbox {
		unshuffled[newIndex] = runes[i]
	}

	return string(unshuffled)
}
func getPbox(lastChar byte) []int {
	switch lastChar {
	case '1':
		return []int{3, 10, 4, 15, 14, 6, 13, 1, 11, 8, 9, 7, 2, 5, 0, 12, 16}
	case '2':
		return []int{12, 9, 8, 5, 7, 6, 0, 4, 11, 2, 15, 10, 1, 3, 14, 13, 16}
	case '3':
		return []int{2, 15, 12, 11, 8, 1, 10, 14, 0, 4, 7, 13, 9, 3, 5, 6, 16}
	case '4':
		return []int{11, 9, 13, 5, 10, 4, 15, 12, 8, 3, 1, 14, 7, 0, 6, 2, 16}
	case '5':
		return []int{10, 0, 15, 4, 6, 7, 3, 9, 14, 2, 8, 13, 12, 5, 11, 1, 16}
	case '6':
		return []int{15, 4, 13, 7, 5, 2, 1, 12, 10, 14, 6, 3, 0, 9, 11, 8, 16}
	case '7':
		return []int{8, 15, 6, 13, 2, 9, 1, 14, 11, 3, 4, 12, 5, 0, 7, 10, 16}
	case '8':
		return []int{5, 10, 1, 9, 14, 11, 3, 0, 13, 7, 15, 8, 4, 6, 2, 12, 16}
	case '9':
		return []int{13, 6, 1, 2, 3, 14, 0, 7, 15, 11, 4, 12, 8, 5, 9, 10, 16}
	case 'A':
		return []int{0, 13, 11, 14, 1, 8, 7, 9, 2, 15, 5, 10, 12, 3, 6, 4, 16}
	case 'B':
		return []int{6, 11, 15, 5, 7, 13, 8, 2, 0, 12, 9, 1, 3, 4, 14, 10, 16}
	case 'C':
		return []int{10, 14, 7, 8, 3, 15, 4, 5, 11, 12, 13, 0, 9, 2, 1, 6, 16}
	case 'D':
		return []int{4, 6, 15, 3, 0, 13, 10, 1, 8, 2, 14, 7, 11, 12, 9, 5, 16}
	case 'E':
		return []int{13, 3, 11, 4, 7, 2, 6, 12, 5, 14, 10, 0, 9, 15, 8, 1, 16}
	case 'F':
		return []int{1, 8, 13, 3, 12, 10, 0, 11, 4, 6, 9, 15, 5, 14, 7, 2, 16}
	default:
		return []int{7, 1, 12, 2, 15, 5, 3, 4, 13, 6, 9, 0, 11, 10, 8, 14, 16}
	}
}
