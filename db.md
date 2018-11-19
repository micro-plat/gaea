### 1. 供货商信息[sys_supplier_info]

| 字段名         | 类型          | 默认值  | 为空  |   约束    | 描述               |
| ------------- | ------------ | :----: | :--: | :-------: | :---------------- |
| sp_id         | number(10)   |        |  否  | PK,IS,SEQ | 供货商编号          |
| sp_name       | varchar2(64) |        |  否  |    IUS    | 供货商名称          |
| balance       | number(20,5) |   0    |  否  |    IS     | 账户余额            |
| conpon_prd_mq | varchar2(32) |        |  否  |    IS     | 制券队列            |
| status        | number(1)    |        |  否  |    IUS    | 状态(0 上架 1 下架)  |

### 2. 供货商信息2[sys_supplier_info2]

| 字段名         | 类型          | 默认值  | 为空  |   约束    | 描述               |
| ------------- | ------------ | :----: | :--: | :-------: | :---------------- |
| sp_id         | number(10)   |        |  否  | PK,IS,SEQ | 供货商编号          |
| sp_name       | varchar2(64) |        |  否  |    IUS    | 供货商名称          |
| balance       | number(20,5) |   0    |  否  |    S     | 账户余额            |
| conpon_prd_mq | varchar2(32) |        |  否  |    S     | 制券队列            |
| status        | number(1)    |        |  否  |    US    | 状态(0 上架 1 下架)  |