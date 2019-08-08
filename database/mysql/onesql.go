//
// 按照JAVA的操作习惯，也声明一个OneSql的结构体
//
// robot.guo
package mysql

type OneSql struct {
	Sql         string        // sql
	EfRows      int           // sql预期影响的行数，-1表示大于0即可 -2表示无限制
	Params      []interface{} // 参数
	DatabaseKey string        // databaseKey
}

const (
	EffectRowMoreThanZero = -1
	EffectRowAny          = -2
)

func (sql *OneSql) FutureEfRows(efRows int) bool {
	if efRows < 0 {
		return false
	}
	switch {
	case sql.EfRows == EffectRowAny:
		return true
	case sql.EfRows == EffectRowMoreThanZero:
		if efRows > 0 {
			return true
		}
	default:
		if sql.EfRows == efRows {
			return true
		}
	}
	return false
}
