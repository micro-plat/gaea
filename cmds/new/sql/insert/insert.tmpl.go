package insert

const InsertTmpl = `
//Create{{.name|cname}} 添加{{.desc}}
const Insert{{.name|cname}} = 'insert into {{.name}}
({{range $i,$c:=.createcolumns}}{{$c.name}}{{if $c.end}},{{end}}{{end}})
values({{range $i,$c:=.createcolumns}}@{{$c.name}}{{if $c.end}},{{end}}{{end}})'
`
const InsertFunc = `
//Create 添加{{.desc}}
func(d *Db{{.name|cname}}) Create(input Create{{.name|cname}}) error {
	db := d.c.GetRegularDB()

	params := map[string]interface{}{
		{{range $i,$c:=.createcolumns}}"{{$c.name}}": input.{{$c.name|cname}} {{if $c.end}},{{end}}{{end}},
	}
	_, q, a, err := db.Execute(sql.Insert{{.name|cname}}, params)
	if err != nil {
		return fmt.Errorf("添加{{.desc}}数据发生错误(err:%v),sql:%s,参数：%v", err, q, a)
	}
	return nil
}
`
