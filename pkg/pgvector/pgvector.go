/*
 * @Author: cloudyi.li
 * @Date: 2023-04-26 10:26:35
 * @LastEditTime: 2023-04-27 10:44:54
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/pgvector/pgvector.go
 */
package pgvector

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

type Vector struct {
	vec []float32
}

func (v Vector) String() string {
	var buf strings.Builder
	buf.WriteString("[")

	for i := 0; i < len(v.vec); i++ {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.FormatFloat(float64(v.vec[i]), 'f', -1, 32))
	}

	buf.WriteString("]")
	return buf.String()
}

func (v *Vector) Parse(s string) error {
	v.vec = make([]float32, 0)
	sp := strings.Split(s[1:len(s)-1], ",")
	for i := 0; i < len(sp); i++ {
		n, err := strconv.ParseFloat(sp[i], 32)
		if err != nil {
			return err
		}
		v.vec = append(v.vec, float32(n))
	}
	return nil
}

func (v *Vector) Scan(src interface{}) (err error) {
	switch src := src.(type) {
	case []byte:
		return v.Parse(string(src))
	case string:
		return v.Parse(src)
	default:
		return fmt.Errorf("unsupported data type: %T", src)
	}
}

func (v Vector) Value() (driver.Value, error) {
	return v.String(), nil
}
