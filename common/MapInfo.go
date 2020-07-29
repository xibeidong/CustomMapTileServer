package common

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var DbMySql *sql.DB

type RoadMapInfo struct {
	ID_level uint8
	ID_dir   uint64
	ID_png   uint64
	ID       string
	ImgData  *[]byte
}
type MysqlInfo struct {
	MysqlDataSourceName string
	MapTableName        string
}

func (roadMap *RoadMapInfo) ToMySql(tableName string) bool {
	smt, err := DbMySql.Prepare("insert into " + tableName + " (img,level_id,dir_id,png_id,id) values (?,?,?,?,?)")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer smt.Close()

	_, err = smt.Exec(roadMap.ImgData, roadMap.ID_level, roadMap.ID_dir, roadMap.ID_png, roadMap.ID)

	if err != nil {
		fmt.Println(err)
		return false
	} else {
		//fmt.Println(ret.RowsAffected())
		return true
	}

}
