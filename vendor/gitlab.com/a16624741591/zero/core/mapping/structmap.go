package mapping

import (
	"bytes"
	"encoding/json"
)

// StructToMap 结构体转map
func StructToMap(s interface{}, number ...bool) (map[string]interface{}, error) {
	var result = make(map[string]interface{})
	if bytesArr, err := json.Marshal(s); err != nil {
		return nil, err
	} else {
		d := json.NewDecoder(bytes.NewReader(bytesArr))
		// 设置将float64转为一个number
		if len(number) > 0 && number[0] {
			d.UseNumber()
		}
		if err := d.Decode(&result); err != nil {
			return nil, err
		}
		return result, nil
	}
}

// MapToStruct map解析到结构体
func MapToStruct(m interface{}, s interface{}) error {
	jsonByte, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonByte, s)
}
