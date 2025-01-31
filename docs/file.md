# file
文件相关文档

## 命名
- 文件长度上限为255字符，路径长度无上限
- 正常文件名称中不能包含以下字符：< > / \ | : * ?，开头和结尾不能有空格。特殊文件除外
- 区分大小写
- 当上传文件、从回收站复原复原等场景下，遇到重名文件时，会在文件名称(不包括扩展名)后面加上序号
- filename与filename_gbk字段任何时候都应对应，无论是否为特殊文件

## 目录
- 目录的层级最多为30层，超出30层则会报错
- 根目录的父目录id为null，且仅有根目录的父目录id为null
- 根目录的文件名为空字符串，且仅有根目录的文件名为空字符串

## 特殊文件
- 回收站目录：名称为：`:recycle.bin`，在根目录下
- 根目录：名称为空字符串

## 删除文件
- 删除文件后，将该文件放入回收站，并重命名为一个随机的字符串
- 恢复文件时，如果原目录已删除，则需要自己选择一个目录
- 彻底删除后，会删除数据库中的记录。如果文件大于10MB，则会检查数据库中是否存在相同文件，不存在则删除Minio中的该文件
- 删除文件后，不会恢复容量；彻底删除后，才会恢复容量

## 分享文件
- 当分享文件时，需要选择哪些用户可见(白名单)或哪些用户不可见(黑名单)。自己不需要设置，始终可见
- 其他用户转存分享的文件时，会递归复制该文件及其所有子文件到指定目录
- 当其他用户转存后，修改分享的文件不会同步到其他用户。其他用户转存前修改，则会同步
- 分享的文件不允许删除和移动，但是可以重命名，其子文件可以删除和移动

## 文件直链
- 创建直链后，可以公开访问该文件，不需要鉴权
- 文件夹不能创建直链
- 直链文件在浏览器上访问时，默认不会触发下载
- 直链URL为：/file/link/:link_name

## 文件上传
- 所有的文件上传接口均为单文件上传，多文件上传可以多次调用接口
- 如果文件已存在，则无需重新上传

### 小文件上传
- 仅有不大于20MB的文件才能进行小文件上传，大于20MB则报错
- 直接将整个文件上传即可

### 大文件上传
- 大文件上传步骤为：预上传、分片上传、完成上传。如果中断，可以获取上传进度，并进行断点续传。
- 预上传
    - 将文件基本信息发送到后端
    - 后端通过uuid生成上传唯一标识，并保存到redis中，过期时间为1天
    - 返回唯一标识给前端
- 分片上传
    - 每个分片顺序上传，在请求头加上分片字节范围
    - 分片大小建议为5-10MB。不得大于20MB，否则报错
    - 每上传完成一个分片，需要更新redis中的进度
- 完成上传
    - 后端校验文件基本信息，校验通过则写入minio，否则返回上传失败

### 上传目录
- 仅提供单文件上传和创建目录的接口，上传目录则可以多次调用接口
- 如果文件所在目录已经创建，则可以并发上传，并发量建议3-10。建议采用深度优先搜索目录上传
