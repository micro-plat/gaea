package tmpls

//DbHeadTpl Db文件头模板
const DbHeadTpl = `
package %s
{{ $empty := "" -}}
import (
	"fmt"
	"%smodules/const/sql"
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/lib4go/db"
	"github.com/micro-plat/lib4go/types"
)

//Create{{.name|cname}} 创建{{.desc}} 
type Create{{.name|cname}} struct {		
	{{range $i,$c:=.createcolumns}}
	//{{$c.name|cname}} {{$c.desc|cname}}
	{{$c.name|cname}} {{$c.type|cstype}} 'json:"{{$c.name|lower}}" form:"{{$c.name|lower}}" m2s:"{{$c.name|lower}}" valid:"required"'
	{{end}}	
}

//Update{{.name|cname}} 添加{{.desc}} 
type Update{{.name|cname}} struct {	
	{{range $i,$c:=.pk}}
	//{{$c.name|cname}} {{$c.desc|cname}}
	{{$c.name|cname}} {{$c.type|cstype}} 'json:"{{$c.name|lower}}" form:"{{$c.name|lower}}" m2s:"{{$c.name|lower}}" valid:"required"'
	{{end -}}	
	{{range $i,$c:=.updatecolumns}}
	//{{$c.name|cname}} {{$c.desc|cname}}
	{{$c.name|cname}} {{$c.type|cstype}} 'json:"{{$c.name|lower}}" form:"{{$c.name|lower}}" m2s:"{{$c.name|lower}}"'
	{{end}}	
}

//Query{{.name|cname}} 查询{{.desc}} 
type Query{{.name|cname}} struct {		
	{{range $i,$c:=.querycolumns}}
	//{{$c.name|cname}} {{$c.desc|cname}}
	{{$c.name|cname}} {{$c.type|cstype}} 'json:"{{$c.name|lower}}" form:"{{$c.name|lower}}" m2s:"{{$c.name|lower}}"'
	{{end}}
	Pi string 'json:"pi" form:"pi" m2s:"pi" valid:"required"'
	
	Ps string 'json:"ps" form:"ps" m2s:"ps" valid:"required"'
}

//IDb{{.name|cname}}  {{.desc}}接口
type IDb{{.name|cname}} interface {

	//Create 创建
	Create(input *Create{{.name|cname}}) (error)

	//Get 单条查询
	Get({{range $i,$c:=.pk -}}{{$c.name|aname}} {{$c.type|ctype}}{{if $c.end}},{{end}}{{end -}})(db.QueryRow,error)

	//Query 列表查询
	Query(input *Query{{.name|cname}})  (data db.QueryRows, count int, err error)

	//Update 更新
	Update(input *Update{{.name|cname}}) (err error)

	//Delete 删除
	Delete({{range $i,$c:=.pk -}}{{$c.name|aname}} {{$c.type|ctype}}{{if $c.end}},{{end}}{{end -}}) (err error)

	{{if ne .di $empty -}}
	//Get{{.name|cname}}Dictionary 获取数据字典
	Get{{.name|cname}}Dictionary({{if ne .dp $empty -}}t string{{- end}}) (db.QueryRows,error)
	{{- end}}
}

//Db{{.name|cname}} {{.desc}}对象
type Db{{.name|cname}} struct {
	c component.IContainer
}

//NewDb{{.name|cname}} 创建{{.desc}}对象
func NewDb{{.name|cname}}(c component.IContainer) *Db{{.name|cname}} {
	return &Db{{.name|cname}}{
		c: c,
	}
}

{{if ne .di $empty -}}
//Get{{.name|cname}}Dictionary 获取数据字典
func(d *Db{{.name|cname}}) Get{{.name|cname}}Dictionary({{if ne .dp $empty -}}t string{{- end}}) (db.QueryRows,error) {

	db := d.c.GetRegularDB()
	data, _, _, err := db.Query(sql.Get{{.name|cname}}Dictionary, map[string]interface{}{
		{{if ne .dp $empty -}}"{{.dp}}": t,{{- end}}
	})
	if err != nil {
		return nil, fmt.Errorf('获取{{.desc}}数据字典发生错误')
	}
	return data, nil
}
{{- end}}
`

//BaseTpl .
const BaseTpl = `package %s
import "github.com/micro-plat/hydra/component"
type I%s interface {
}
type %s struct {
	c component.IContainer
}
func New%s(c component.IContainer) *%s{
	return &%s{
		c: c,
	}
}
`

const DicTpl = `package sql
{{$empty := "" -}}
{{if ne .di $empty -}}
//Get{{.name|cname}}Dictionary  获取数据字典
const Get{{.name|cname}}Dictionary = 'select {{.di}} as id,{{.dn}} as name from {{.name}}{{.dblink}} where 1=1 
{{if ne .dp $empty}} &{{.dp}} {{end}}'
{{- end}}
`
