package postgres

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	// значения массива без кавычек не должны содержать: (" , \ { } пробел NULL)
	//и должны быть хотя бы одним символом
	unquotedChar  = `[^",\\{}\s(NULL)]`
	unquotedValue = fmt.Sprintf("(%s)+", unquotedChar)

	// Значения массива в кавычках заключены в двойные кавычки, могут быть
	//любыми символами, кроме " или \, которые должны быть экранированы обратной косой чертой:

	quotedChar  = `[^"\\]|\\"|\\\\`
	quotedValue = fmt.Sprintf("\"(%s)*\"", quotedChar)

	// значение массива может быть как в кавычках, так и без кавычек:
	arrayValue = fmt.Sprintf("(?P<value>(%s|%s))", unquotedValue, quotedValue)

	// Значения массива разделяются запятой, ЕСЛИ значений больше одного
	arrayExp = regexp.MustCompile(fmt.Sprintf("((%s)(,)?)", arrayValue))

	valueIndex int
)

// Находим индекс именованного выражения 'value'
func init() {
	for i, exp := range arrayExp.SubexpNames() {
		if exp == "value" {
			valueIndex = i
			break
		}
	}
}

type Array []string

func (s *Array) String() string {
	return `{"` + strings.Join(*s, `","`) + `"}`
}

// Scan реализует sql.Scanner для типа String slice
// Сканеры принимают значение базы данных (в данном случае в виде фрагмента байта)
// и устанавливает значение типа.  Здесь мы приводим к строке и
// выполняем синтаксический анализ на основе регулярных выражений
func (s *Array) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return errors.New("scan source was not []bytes")
	}
	*s = parseArray(string(bytes))

	return nil
}

// разобрать массив проанализируйте выходную строку из типа массива.
// Используется регулярное выражение: (((?P<значение>(([^",\\{}\ s(НУЛЕВОЙ)])+|"([^"\\]|\\"|\\\\)*")))(,)?)
func parseArray(array string) (results []string) {
	matches := arrayExp.FindAllStringSubmatch(array, -1)
	for _, match := range matches {
		s := match[valueIndex]
		// строка _может_ быть заключена в кавычки, поэтому обрежьте их:
		s = strings.Trim(s, "\"")
		results = append(results, s)
	}
	return
}
