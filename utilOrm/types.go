package utilOrm

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type OrmMysqlJsonSliceString []string

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *OrmMysqlJsonSliceString) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Failed to unmarshal string value: %v ", value)
	}
	return json.Unmarshal(bytes, j)
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j *OrmMysqlJsonSliceString) Value() (driver.Value, error) {
	return json.Marshal(j)
}
