package sql

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/micro-plat/gaea/cmds"
	"github.com/micro-plat/gaea/cmds/new/sql/oracle"
	"github.com/micro-plat/gaea/cmds/new/util/conf"
	"github.com/micro-plat/gaea/cmds/new/util/md"
	"github.com/micro-plat/gaea/cmds/new/util/path"
	"github.com/urfave/cli"
)

type sqlCmd struct {
}

func NewSqlCmd() cli.Command {
	p := &sqlCmd{}
	return cli.Command{
		Name:   "sql",
		Usage:  "创建SQL语句",
		Flags:  p.geStartFlags(),
		Action: p.action,
	}
}

func (p *sqlCmd) geStartFlags() []cli.Flag {
	return nil
}

func (p *sqlCmd) action(c *cli.Context) (err error) {
	mdFilePath := path.GetMDPath()
	outPath := path.GetModulePath()
	if c.NArg() > 0 {
		mdFilePath = []string{c.Args().Get(0)}
	}
	if c.NArg() > 1 {
		outPath = c.Args().Get(1)
	}
	var tables = make([]*conf.Table, 0, 2)
	for _, p := range mdFilePath {
		t1, err := md.Markdown2Table(p)
		if err != nil {
			cmds.Log.Error(err)
			return err
		}
		tables = append(tables, t1...)
	}

	tmpls, err := oracle.GetTmples(tables)
	if err != nil {
		cmds.Log.Error(err)
		return err
	}
	if len(tmpls) == 0 {
		cmds.Log.Errorf("未找到数据表信息%v", mdFilePath)
		return nil
	}
	if err = p.createFile(outPath, tmpls); err != nil {
		cmds.Log.Error(err)
		return err
	}
	cmds.Log.Infof("SQL生成完成,共生成%d个文件", len(tmpls))
	return nil
}

func (p *sqlCmd) createFile(root string, data map[string]string) error {
	for k, v := range data {
		path := filepath.Join(root, k)
		dir := filepath.Dir(path)
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			err := os.MkdirAll(dir, os.ModeDir|os.ModePerm)
			if err != nil {
				err = fmt.Errorf("创建文件夹%s失败:%v", path, err)
				return err
			}
		}
		if _, err := os.Stat(path); err == nil || os.IsExist(err) {
			cmds.Log.Warn("文件已存在:", path)
			continue
		}
		srcf, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModeAppend|os.ModePerm)
		if err != nil {
			err = fmt.Errorf("无法打开文件:%s(err:%v)", path, err)
			return err
		}
		defer srcf.Close()
		_, err = srcf.WriteString(v)
		if err != nil {
			return err
		}
		cmds.Log.Info("生成文件:", path)
	}
	return nil

}
