package data

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/micro-plat/gaea/cmds"
	"github.com/micro-plat/gaea/cmds/new/module/tmpls"
	"github.com/micro-plat/gaea/cmds/new/util/conf"
)

//GetInputData 获取模板数据
func getInputData(tb *conf.Table) map[string]interface{} {
	input := map[string]interface{}{
		"name":          tb.Name,
		"desc":          tb.Desc,
		"createcolumns": getCreateColumns(tb),
		"querycolumns":  getQueryColumns(tb),
		"updatecolumns": getUpdateColumns(tb),
		"selectcolumns": getSelectColumns(tb),
		"pk":            getPks(tb),
		"seqs":          getSeqs(tb),
	}

	return input
}

func getCreateColumns(tb *conf.Table) []map[string]interface{} {
	columns := make([]map[string]interface{}, 0, len(tb.CNames))

	for i, v := range tb.CNames {
		if tb.Cons[i] == "" || tb.Cons[i] == "-" {
			panic("数据表没有指定约束")
		}
		if strings.Contains(tb.Cons[i], "I") && !strings.Contains(tb.Cons[i], "SEQ") {
			row := map[string]interface{}{
				"name": v,
				"desc": tb.Descs[i],
				"type": tb.Types[i],
				"len":  tb.Lens[i],
				"end":  i != len(tb.CNames)-1,
			}
			columns = append(columns, row)
		}

	}
	if len(columns) > 0 {
		columns[len(columns)-1]["end"] = false
	}
	return columns
}

func getQueryColumns(tb *conf.Table) []map[string]interface{} {
	columns := make([]map[string]interface{}, 0, len(tb.CNames))

	for i, v := range tb.CNames {
		if tb.Cons[i] == "" || tb.Cons[i] == "-" {
			panic("数据表没有指定约束")
		}
		if strings.Contains(tb.Cons[i], "Q") && !strings.Contains(tb.Cons[i], "SEQ") {
			row := map[string]interface{}{
				"name": v,
				"desc": tb.Descs[i],
				"type": tb.Types[i],
				"len":  tb.Lens[i],
				"end":  i != len(tb.CNames)-1,
			}
			columns = append(columns, row)
		}

	}
	if len(columns) > 0 {
		columns[len(columns)-1]["end"] = false
	}
	return columns
}

func getUpdateColumns(tb *conf.Table) []map[string]interface{} {
	columns := make([]map[string]interface{}, 0, len(tb.CNames))

	for i, v := range tb.CNames {
		if tb.Cons[i] == "" || tb.Cons[i] == "-" {
			panic("数据表没有指定约束")
		}
		if strings.Contains(tb.Cons[i], "U") && !strings.Contains(tb.Cons[i], "SEQ") && !strings.Contains(tb.Cons[i], "PK") {
			row := map[string]interface{}{
				"name": v,
				"desc": tb.Descs[i],
				"type": tb.Types[i],
				"len":  tb.Lens[i],
				"end":  i != len(tb.CNames)-1,
			}
			columns = append(columns, row)
		}

	}
	if len(columns) > 0 {
		columns[len(columns)-1]["end"] = false
	}
	return columns
}

func getSelectColumns(tb *conf.Table) []map[string]interface{} {
	columns := make([]map[string]interface{}, 0, len(tb.CNames))

	for i, v := range tb.CNames {
		if tb.Cons[i] == "" || tb.Cons[i] == "-" {
			panic("数据表没有指定约束")
		}
		if strings.Contains(tb.Cons[i], "S") {
			row := map[string]interface{}{
				"name": v,
				"desc": tb.Descs[i],
				"type": tb.Types[i],
				"len":  tb.Lens[i],
				"end":  i != len(tb.CNames)-1,
			}
			columns = append(columns, row)
		}

	}
	if len(columns) > 0 {
		columns[len(columns)-1]["end"] = false
	}
	return columns
}
func getPks(tb *conf.Table) []map[string]interface{} {
	columns := make([]map[string]interface{}, 0, len(tb.CNames))

	for i, v := range tb.CNames {
		if strings.Contains(tb.Cons[i], "PK") {
			row := map[string]interface{}{
				"name": v,
				"desc": tb.Descs[i],
				"type": tb.Types[i],
				"len":  tb.Lens[i],
				"end":  i != len(tb.CNames)-1,
			}
			columns = append(columns, row)
		}
	}
	if len(columns) > 0 {
		columns[len(columns)-1]["end"] = false
	}
	return columns
}

func getSeqs(tb *conf.Table) []map[string]interface{} {
	columns := make([]map[string]interface{}, 0, len(tb.CNames))

	for i, v := range tb.CNames {
		if strings.Contains(tb.Cons[i], "SEQ") {
			row := map[string]interface{}{
				"name":    v,
				"seqname": fmt.Sprintf("seq_%s_%s", fGetNName(tb.Name), getFilterName(tb.Name, v)),
				"desc":    tb.Descs[i],
				"type":    tb.Types[i],
				"len":     tb.Lens[i],
				"end":     i != len(tb.CNames)-1,
			}
			columns = append(columns, row)
		}
	}
	if len(columns) > 0 {
		columns[len(columns)-1]["end"] = false
	}
	return columns
}

func fGetNName(n string) string {
	items := strings.Split(n, "_")
	if len(items) <= 1 {
		return n
	}
	return strings.Join(items[1:], "_")
}

func getFilterName(t string, f string) string {
	text := make([]string, 0, 1)
	tb := strings.Split(t, "_")
	fs := strings.Split(f, "_")
	for _, v := range fs {
		ex := false
		for _, k := range tb {
			if v == k {
				ex = true
				break
			}
		}
		if !ex {
			text = append(text, v)
		}
	}
	if len(text) == 0 {
		return "id"
	}
	return strings.Join(text, "_")

}

func makeFunc() map[string]interface{} {
	return map[string]interface{}{
		"aname": fGetAName,
		"cname": fGetCName,
		"ctype": fGetType,
		"lname": fGetLastName,
		"lower": fToLower,
	}
}

func fGetAName(n string) string {
	items := strings.Split(n, "_")
	nitems := make([]string, 0, len(items))
	for k, i := range items {
		if k == 0 {
			nitems = append(nitems, i)
		}
		if k > 0 {
			nitems = append(nitems, strings.ToUpper(i[0:1])+i[1:])
		}

	}
	return strings.Join(nitems, "")
}

func fGetCName(n string) string {
	items := strings.Split(n, "_")
	nitems := make([]string, 0, len(items))
	for _, i := range items {
		nitems = append(nitems, strings.ToUpper(i[0:1])+i[1:])
	}
	return strings.Join(nitems, "")
}
func fGetType(n string) string {
	switch {
	case strings.Contains(n, "varchar"):
		return "string"
	case strings.Contains(n, "number"):
		if strings.Contains(n, ",") {
			return "float64"
		}
		var i, j int
		for k, v := range n {
			if v == '(' {
				i = k
			}
			if v == ')' {
				j = k
			}
		}
		ii, _ := strconv.Atoi(n[i+1 : j])
		if ii < 10 {
			return "int"
		}
		return "int64"

	case strings.Contains(n, "date"):
		return "time.Time"
	default:
		return "string"
	}
}
func fGetLastName(n string) string {
	names := strings.Split(strings.Trim(n, "/"), "/")
	return names[len(names)-1]
}

func fToLower(s string) string {
	return strings.ToLower(s)
}

//translate  .
func translate(tag string, tplName string, input interface{}) (string, error) {
	var tmpl = template.New(tag).Funcs(makeFunc())
	np, err := tmpl.Parse(tplName)
	if err != nil {
		return "", err
	}
	buff := bytes.NewBufferString("")
	if err := np.Execute(buff, input); err != nil {
		return "", err
	}
	return buff.String(), nil
}

//GetTmples 获取模板
//@tag 模板标签
//@tplName 模板名字
//@tbs 表结构体
//@filters 过滤字段
//@makeFunc 是否生成函数
//@return out 返回数据
func GetTmples(tag, tplName string, tbs []*conf.Table, filters []string, makeFunc bool, modulePath string) (out map[string]map[string]string, err error) {
	out = map[string]map[string]string{}
	for _, tb := range tbs {
		if len(filters) > 0 {
			e := false
			for _, f := range filters {
				if strings.EqualFold(tb.Name, f) {
					e = true
					break
				}
			}
			if !e {
				continue
			}
		}
		//获取模板数据
		input := getInputData(tb)
		//翻译模板
		content, err := translate(tag, tplName, input)

		if err != nil {
			return nil, err
		}
		if makeFunc { //生成函数
			c := make(map[string]string)
			if strings.Contains(modulePath, "sql") || !strings.Contains(modulePath, "modules") {
				modulePath = "modules"
			}
			c[fmt.Sprintf(modulePath+"/%s.go", strings.Replace(tb.Name, "_", "/", -1))] = strings.Replace(content, "'", "`", -1)
			head, err := translate("head", tmpls.DbHeadTpl, input)
			if err != nil {
				return nil, err
			}
			c["head"] = strings.Replace(head, "'", "`", -1)
			out[fmt.Sprintf(modulePath+"/%s.go", strings.Replace(tb.Name, "_", "/", -1))] = c
		} else { //生成sql
			c := make(map[string]string)

			modulePath = "modules/const/sql"

			c[fmt.Sprintf(modulePath+"/%s.go", strings.Replace(tb.Name, "_", ".", -1))] = strings.Replace(content, "'", "`", -1)
			out[fmt.Sprintf(modulePath+"/%s.go", strings.Replace(tb.Name, "_", ".", -1))] = c
		}
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

//createFile
//创建并生成文件
func createFile(add bool, data map[string]map[string]string) error {
	for k, v := range data {
		_, ok := v["head"]
		if ok { //生成函数文件头
			_, err := os.Stat(k)
			//不存在则创建
			if err != nil {
				add = true
				os.MkdirAll(path.Dir(k), os.ModePerm)
				f, err := os.Create(k)
				if err != nil {
					err = fmt.Errorf("无法创建文件:%s(err:%v)", k, err)
					return err
				}
				defer f.Close()
				m := strings.Split(k, "/")
				absPath, _ := filepath.Abs(k)
				i := strings.Index(absPath, "src")
				j := strings.Index(absPath, "modules")
				_, err = f.WriteString(fmt.Sprintf(v["head"], m[len(m)-2], absPath[i+4:j]))
				if err != nil {
					return err
				}
				cmds.Log.Info("写入crud函数头部文件成功:", k)
			}
		} else { //生成sql文件头
			_, err := os.Stat(k)
			//不存在则创建
			if err != nil {
				add = true
				os.MkdirAll(path.Dir(k), os.ModePerm)
				f, err := os.Create(k)
				if err != nil {
					err = fmt.Errorf("无法创建文件:%s(err:%v)", k, err)
					return err
				}
				defer f.Close()
				_, err = f.WriteString("package sql")
				if err != nil {
					return err
				}
				cmds.Log.Info("写入sql头部文件成功:", k)
			}
		}
		if !add {
			return fmt.Errorf("文件已经存在：%s，未输入 -add 不执行任何操作", k)
		}
		srcf, err := os.OpenFile(k, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModeAppend|os.ModePerm)
		if err != nil {
			err = fmt.Errorf("无法打开文件:%s(err:%v)", k, err)
			return err
		}
		defer srcf.Close()
		_, err = srcf.WriteString(v[k])
		if err != nil {
			return err
		}
		cmds.Log.Info("写入文件成功:", k)

	}
	return nil

}

//CreateModulesFile 创建 modules 文件
func CreateModulesFile(add, cover bool, tmpls map[string]map[string]string) (err error) {
	if cover {
		for k := range tmpls {
			cmds.Log.Warnf("覆盖文件：%s", k)
			os.Remove(k)
		}
	}
	//创建文件
	if err = createFile(add, tmpls); err != nil {
		cmds.Log.Error(err)
		return err
	}
	return nil
}
