package CivetTarsDataBase

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type WarpPacket struct {
	Base map[string]string //记录返回结果
	Keys []string          //记录列名顺序
}

type DataBase struct {
	UserName string
	PassWord string
	Host     string
	Port     string
}
type CivetData struct {
	conf *DataBase
	Db   *gorm.DB
}
type DBModel struct {
	DataBaseName string `json:"DataBaseName"`
	TableName    string `json:"TableName"`
	Column       []struct {
		Name string `json:"Name"`
		Type string `json:"Type"`
	} `json:"Column"`
}
type EmptyTable struct {
	ID int
}

type Column struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

func makeEmptyStr() *sql.NullString {
	return &sql.NullString{}
}
func (CD *CivetData) QueryRowsConditionWithOutModel(tableName string, Ins *[]*WarpPacket, SearchKey string, SearchValue string) bool {
	rows, _ := CD.Db.Raw("select * FROM " + tableName + " where " + SearchKey + "=" + SearchValue).Rows()
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
		return false
	}
	manages := Ins
	for rows.Next() {
		keys := make([]string, 0, len(columns))
		resultmap := make(map[string]string, len(columns))
		rowdatas := make([]interface{}, len(columns))
		for i := 0; i < len(rowdatas); i++ {
			rowdatas[i] = makeEmptyStr()
		}
		rows.Scan(rowdatas...)
		for i := 0; i < len(rowdatas); i++ {
			k := columns[i]
			val := (rowdatas[i].(*sql.NullString))
			v := ""
			if val.Valid {
				v = val.String
			}
			resultmap[k] = v
			keys = append(keys, k)
		}
		*manages = append(*manages, &WarpPacket{Base: resultmap, Keys: keys})
	}
	return true
}
func (CD *CivetData) QueryRowsAllWithOutModel(tableName string, Mana *[]*WarpPacket) bool {
	rows, _ := CD.Db.Raw("select * FROM " + tableName).Rows()
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
		return false
	}
	manages := Mana
	for rows.Next() {
		keys := make([]string, 0, len(columns))
		resultmap := make(map[string]string, len(columns))
		rowdatas := make([]interface{}, len(columns))
		for i := 0; i < len(rowdatas); i++ {
			rowdatas[i] = makeEmptyStr()
		}
		rows.Scan(rowdatas...)
		for i := 0; i < len(rowdatas); i++ {
			k := columns[i]
			val := (rowdatas[i].(*sql.NullString))
			v := ""
			if val.Valid {
				v = val.String
			}
			resultmap[k] = v
			keys = append(keys, k)
		}
		*manages = append(*manages, &WarpPacket{Base: resultmap, Keys: keys})
	}
	return true
}
func (CD *CivetData) DeleteRowWithOutModel(tableName string, SearchKey string, SearchValue string) bool {
	s := "DELETE FROM " + tableName + " WHERE " + SearchKey + "=" + SearchValue
	res := CD.Db.Exec(s)
	if res.Error != nil {
		print(res.Error.Error())
		return false
	} else {
		return true
	}
}
func (CD *CivetData) EditRowWithOutModel(tableName string, SearchKey string, SearchValue string, col []Column) (bool, error) {
	s := ""
	for index, value := range col {
		fmt.Println(index, value.Key, value.Value)
		s = s + fmt.Sprintf("%s=%s", value.Key, value.Value)
		if index+1 != len(col) {
			s = s + ","
		}
	}
	sql2 := fmt.Sprintf("update %s set %s where "+SearchKey+"="+SearchValue, tableName, s)
	res := CD.Db.Exec(sql2)
	if res.Error != nil {
		return false, res.Error
	} else {
		return true, nil
	}
}
func (CD *CivetData) CreateRowWithOutModel(tableName string, col []Column) bool {
	sqlK := ""
	sqlV := ""
	for index, value := range col {
		fmt.Println(index, value.Key, value.Value)
		sqlK = sqlK + value.Key
		sqlV = sqlV + value.Value
		if index+1 != len(col) {
			sqlK = sqlK + ","
			sqlV = sqlV + ","
		}
	}
	sqlH := fmt.Sprintf("INSERT INTO %s ", tableName)
	sql2 := fmt.Sprintf("(%s)", sqlK)
	sql3 := fmt.Sprintf(" VALUES (%s)", sqlV)
	fmt.Println(sqlH + sql2 + sql3)
	res := CD.Db.Exec(sqlH + sql2 + sql3)

	if res.Error != nil {
		return false
	} else {
		return true
	}

}

func (CD *CivetData) CreateTableByJson(js []byte, model *DBModel, prefix string) (string, error) {
	err := json.Unmarshal(js, &model)
	if err != nil {
		fmt.Println(err)
		return "", err
	} else {
		s := ""
		for index, value := range model.Column {
			fmt.Println(index, value.Name, value.Type)
			s = s + fmt.Sprintf("`%s` %s NOT NULL", value.Name, value.Type)
			if index+1 != len(model.Column) {
				s = s + ","
			}
		}
		fmt.Println(s)
		CD.Db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s`(%s,PRIMARY KEY ( `id` ))ENGINE=InnoDB DEFAULT CHARSET=utf8;", model.TableName+prefix, s))
		//CD.CreateTableByName(&EmptyTable{}, model.TableName)
		return model.TableName + prefix, nil
	}
}
func (CD *CivetData) ConnectDataBaseByJson(js []byte) {
	model := DBModel{}
	err := json.Unmarshal(js, &model)
	if err != nil {
		fmt.Println(err)
	} else {
		CD.ConnectOrCreateDataBase(model.DataBaseName)
	}
}
func (CD *CivetData) CreateDataBaseByJson(js []byte, model *DBModel) {
	err := json.Unmarshal(js, &model)
	if err != nil {
		fmt.Println(err)
	} else {
		CD.ConnectOrCreateDataBase(model.DataBaseName)
	}
}
func (CD *CivetData) CreateRowByName(TableName string, deviceModel interface{}) bool {
	res := true
	if res == true {
		result := CD.Db.Table(TableName).Create(deviceModel)
		if result.Error == nil {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}
func (CD *CivetData) RemoveDeviceByName(TableName string, TableModel interface{}, SearchKey string, SearchValue string) bool {
	res := CD.Db.Table(TableName).Where(SearchKey+" = ?", SearchValue).Delete(TableModel)
	if res.Error != nil {
		fmt.Println(res.Error)
		fmt.Println("删除出错了")
		return true
	} else {
		fmt.Println("删除没出错")
	}
	return false
}
func (CD *CivetData) EditDevice(TableName string, SearchKey string, SearchValue string, Key string, Value interface{}) bool {
	res := CD.Db.Table(TableName).Where(SearchKey+" = ?", SearchValue).Update(Key, Value)
	if res.Error != nil {
		fmt.Println(res.Error)
		fmt.Println("修改出错了")
		return true
	} else {
		fmt.Println("修改没出错")
	}
	return false
}

func (CD *CivetData) QueryRowsAll(modellist interface{}) bool {
	result := CD.Db.Find(modellist)
	if result.Error != nil {
		return false
	} else {
		return true
	}
}
func (CD *CivetData) QueryRowsWithCondition(modellist interface{}, key string, value interface{}) bool {
	f := fmt.Sprintf("%s = ?", key)
	result := CD.Db.Where(f, value).Find(modellist)
	if result.Error != nil {
		return false
	} else {
		return true
	}
}
func (CD *CivetData) QueryRowWithCondition(model interface{}, key string, value interface{}) bool {
	f := fmt.Sprintf("%s = ?", key)
	result := CD.Db.Where(f, value).First(model)
	if result.Error != nil {
		return false
	} else {
		return true
	}
}
func (CD *CivetData) QueryRowWithID(model interface{}, id interface{}) {
	CD.Db.First(model, id)
}
func (CD *CivetData) QueryRow(model interface{}) {
	CD.Db.First(model)
}
func (CD *CivetData) EditRowByCondition(model interface{}, key string, value interface{}, editKey string, editValue interface{}) {
	f := fmt.Sprintf("%s = ?", key)
	CD.Db.Model(model).Where(f, value).Update(editKey, editValue)
}
func (CD *CivetData) EditRow(value interface{}) {
	CD.Db.Save(value)
}
func (CD *CivetData) DelRowByCondition(model interface{}, key string, value interface{}) {
	f := fmt.Sprintf("%s = ?", key)
	CD.Db.Where(f, value).Delete(model)
}
func (CD *CivetData) DelRow(Value interface{}) {
	CD.Db.Delete(Value)
}
func (CD *CivetData) CreateRow(Value interface{}) bool {
	res := CD.Db.Create(Value)
	if res.Error != nil {
		return false
	} else {
		return true
	}
}
func (CD *CivetData) CheckTableExist(TableName string) bool {
	return CD.Db.Migrator().HasTable(TableName)
}
func (CD *CivetData) CreateTableByName(DataBase interface{}, Name string) bool {
	first := CD.Db.Migrator().HasTable(Name)
	if first == true {
		return false
	} else {
		ck := CD.Db.Migrator().HasTable(DataBase)
		if ck == true {
			fmt.Println("存在表")
			err := CD.Db.Migrator().RenameTable(DataBase, Name)
			if err != nil {
				return false
			}
		} else {
			fmt.Println("不存在表")
			CD.CreateTable(DataBase)
			err := CD.Db.Migrator().RenameTable(DataBase, Name)
			if err != nil {
				return false
			}
		}
		return CD.Db.Migrator().HasTable(Name)
	}

}
func (CD *CivetData) CreateTable(DataBase interface{}) bool {
	err := CD.Db.Migrator().CreateTable(DataBase)
	if err != nil {
		return false
	} else {
		return true
	}
}
func CreateCivetData(UserName string, Password string, Host string, Port string) *CivetData {
	return &CivetData{
		conf: SetDataBaseInfo(UserName, Password, Host, Port),
	}
}
func SetDataBaseInfo(UserName string, Password string, Host string, Port string) *DataBase {
	return &DataBase{
		UserName: UserName,
		PassWord: Password,
		Host:     Host,
		Port:     Port,
	}
}
func (CD *CivetData) CreateDataBase(DataBaseName string) bool {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", CD.conf.UserName, CD.conf.PassWord, CD.conf.Host, CD.conf.Port)
	dbc, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return false
	}
	dbc.Exec(fmt.Sprintf("Create DATABASE %s", DataBaseName))
	return true
}
func (CD *CivetData) ConnectOrCreateDataBase(DataBaseName string) bool {
	ConRes := CD.ConnectDataBase(DataBaseName, func(DBN string, RES *int) {
		CD.CreateDataBase(DataBaseName)
		res := CD.ConnectDataBase(DataBaseName, func(s string, i *int) {})
		if res == true {
			*RES = 1
		} else {
			*RES = 2
		}
	})
	if ConRes == true {
		return true
	} else {
		return false
	}
}
func (CD *CivetData) ConnectDataBase(DataBaseName string, Fail func(string, *int)) bool {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", CD.conf.UserName, CD.conf.PassWord, CD.conf.Host, CD.conf.Port, DataBaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Print("the err is ", err)
		var res int = 0
		Fail(DataBaseName, &res)
		switch res {
		case 0:
			return false
		case 1:
			return true
		case 2:
			return false
		default:
			return false
		}
	} else {
		fmt.Print("\r\nConnect ", DataBaseName, " Success\r\n")
		CD.Db = db
		return true
	}
}
