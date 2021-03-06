package tmpls

//InsertMysqlTmpl mysql insert sql 模板
const InsertMysqlTmpl = `
//Insert{{.name|cname}} 添加{{.desc}}
const Insert{{.name|cname}} = 'insert into {{.name}}{{.dblink}}
({{range $i,$c:=.createcolumns}}{{$c.name}}{{if $c.end}},{{end}}{{end}})
values({{range $i,$c:=.createcolumns}}@{{$c.name}}{{if $c.end}},{{end}}{{end}})'
`

//InsertOracleTmpl oracle insert sql 模板
const InsertOracleTmpl = `
//Create{{.name|cname}} 添加{{.desc}}
const Insert{{.name|cname}} = 'insert into {{.name}}{{.dblink}}
({{range $i,$c:=.seqs}}{{$c.name}},{{end}}{{range $i,$c:=.createcolumns}}{{$c.name}}{{if $c.end}},{{end}}{{end}})
values({{range $i,$c:=.seqs}}{{$c.seqname}}.nextval,{{end}}{{range $i,$c:=.createcolumns}}@{{$c.name}}{{if $c.end}},{{end}}{{end}})'
`

//InsertFunc mysql insert 函数模板
const InsertFunc = `
//Create 添加{{.desc}}
func(d *Db{{.name|cname}}) Create(input *Create{{.name|cname}}) error {

	db := d.c.GetRegularDB()
	lastInsertID, affectedRow, q, a, err := db.Executes(sql.Insert{{.name|cname}}, map[string]interface{}{
		{{range $i,$c:=.createcolumns -}}
		"{{$c.name}}": input.{{$c.name|cname}},
		{{end -}}
	})
	fmt.Println(lastInsertID, affectedRow)
	if err != nil {
		return fmt.Errorf("添加{{.desc}}数据发生错误(err:%v),sql:%s,参数：%v,lastInsertID:%v,受影响的行数：%v", err, q, a, lastInsertID,affectedRow)
	}
	return nil
}
`
