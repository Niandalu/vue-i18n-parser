## Vue i18n parser

### 适用项目

使用 yml 和 sfc 来做翻译的 vue 项目。具体格式，参考 `testdata/` 文件夹下的内容

### 主要功能

1. 将 yml、sfc 中的待翻译 key-value 进行提取，保存成 csv。无视未修改的文件。
2. 将翻译后的 csv 解析，重新自动写入到各个文件中。可检测冲突，如果出现冲突，则停止写入

### 安装

```bash
git clone git@github.com:Niandalu/vue-i18n-parser.git
make build # 会在当前文件夹生成一个名为 vip 的二进制文件
```

### 使用

#### 收集待翻译内容整理成 csv

```bash
# 在项目的根目录下
vip collect --diff .
```

#### 将翻译后的内容写回到各个文件

```bash
# 在项目的根目录下
vip feed <翻译后的文件>
```

### 已知问题

- 不能使用 `yes`, `no` 等等 yaml 会解析成非 `string` 的字符串作为 key
- csv 看起来只支持 unix 换行符，这个 bug 未来会修


### License

[MIT](http://opensource.org/licenses/MIT)
