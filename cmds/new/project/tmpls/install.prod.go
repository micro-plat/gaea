package tmpls

const installProdTmpl = `// +build prod

package main

import (
	'fmt'
	'os'
	'path/filepath'
	'strings'

	'github.com/micro-plat/hydra/component'
)

//install 用于配置生产环境参数，这些参数使用以#开头的变量命名，当执行 install 命令时会引导安装人员进行参数设置
func (s *{{.projectName|lName}}) install() {
	s.IsDebug = false	

	s.Conf.SetInput("#db_connection_string", "数据库连接串", "username/password@host")
	s.Conf.API.SetMainConf("{'address':':#host_port'}")	
	s.Conf.API.SetSubConf('header', "
				{
					'Access-Control-Allow-Origin': '*', 
					'Access-Control-Allow-Methods': 'GET,POST,PUT,DELETE,PATCH,OPTIONS', 
					'Access-Control-Allow-Credentials': 'true'
				}
			")

	s.Conf.API.SetSubConf('auth', "
		{
			'jwt': {
				'exclude': ['#exclude_url'],
				'expireAt': 36000,
				'mode': 'HS512',
				'name': 'sso',
				'secret': 'ef1a8839cb511780903ff6d5d79cf8f8'
			}
		}
		")

	s.Conf.Plat.SetVarConf('db', 'db', "{			
			'provider':'ora',
			'connString':'#db_connection_string',
			'maxOpen':200,
			'maxIdle':10,
			'lifeTime':600		
	}")

	s.Conf.Plat.SetVarConf('cache', 'cache', "
		{
			'proto':'redis',
			'addrs':[
					#redis_server
			],
			'db':1,
			'dial_timeout':10,
			'read_timeout':10,
			'write_timeout':10,
			'pool_size':100
	}
		")

	//自定义安装程序
	s.Conf.API.Installer(func(c component.IContainer) error {
		if !s.Conf.Confirm('创建数据库表结构,添加基础数据?') {
			return nil
		}
		path, err := getSQLPath()
		if err != nil {
			return err
		}
		sqls, err := s.Conf.GetSQL(path)
		if err != nil {
			return err
		}
		db, err := c.GetDB()
		if err != nil {
			return err
		}
		for _, sql := range sqls {
			if sql != '' {
				if _, q, _, err := db.Execute(sql, map[string]interface{}{}); err != nil {
					if !strings.Contains(err.Error(), 'ORA-00942') {
						s.Conf.Log.Errorf('执行SQL失败： %v %s\n', err, q)
					}
				}
			}
		}
		return nil
	})

}

//getSQLPath 获取getSQLPath
func getSQLPath() (string, error) {
	gopath := os.Getenv('GOPATH')
	if gopath == '' {
		return '', fmt.Errorf('未配置环境变量GOPATH')
	}
	path := strings.Split(gopath, ';')
	if len(path) == 0 {
		return '', fmt.Errorf('环境变量GOPATH配置的路径为空')
	}
	return filepath.Join(path[0], 'src/{{.projectName}}/modules/const/sql'), nil
}
`
