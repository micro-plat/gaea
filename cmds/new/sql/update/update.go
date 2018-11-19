package insert

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/micro-plat/gaea/cmds/new/sql/conf"
)

func translate(c string, input interface{}) (string, error) {
	var tmpl = template.New("update").Funcs(makeFunc())
	np, err := tmpl.Parse(c)
	if err != nil {
		return "", err
	}
	buff := bytes.NewBufferString("")
	if err := np.Execute(buff, input); err != nil {
		return "", err
	}
	return buff.String(), nil
}

//GetTmples .
func GetTmples(tbs []*conf.Table, path string, filters []string) (out map[string]string, err error) {
	out = make(map[string]string, len(tbs))
	for _, tb := range tbs {
		if len(filters) > 0 {
			e := false
			for _, f := range filters {
				if strings.Contains(tb.Name, f) {
					e = true
					break
				}
			}
			if !e {
				continue
			}
		}
		columns := make([]map[string]interface{}, 0, len(tb.CNames))
		for i, v := range tb.CNames {
			//获取可更新的数据的字段
			s := strings.Replace(tb.Cons[i], "UNQ", "", -1)
			if strings.Contains(s, "U") {
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
		columns[len(columns)-1]["end"] = false
		input := map[string]interface{}{
			"name":    tb.Name,
			"desc":    tb.Desc,
			"columns": columns,
			"path":    path,
			"pk":      getPks(tb),
		}
		content, err := translate(updateTmpl, input)
		if err != nil {
			return nil, err
		}
		out[fmt.Sprintf("%s.go", tb.Name)] = strings.Replace(content, "'", "`", -1)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func makeFunc() map[string]interface{} {
	return map[string]interface{}{
		"cname": fGetCName,
		"ctype": fGetType,
		"lname": fGetLastName,
	}
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
	if strings.HasPrefix(n, "nvarchar") {
		return "string"
	} else if strings.HasPrefix(n, "number") {
		if strings.Contains(n, ",") {
			return "float64"
		}
		return "int64"
	} else if strings.HasPrefix(n, "date") {
		return "time.Time"
	}
	return "string"
}
func fGetLastName(n string) string {
	names := strings.Split(strings.Trim(n, "/"), "/")
	return names[len(names)-1]
}

func getPks(tb *conf.Table) []string {
	out := make([]string, 0, 1)
	for i, v := range tb.Cons {
		if strings.Contains(v, "PK") {
			out = append(out, tb.CNames[i])
		}
	}
	return out
}
