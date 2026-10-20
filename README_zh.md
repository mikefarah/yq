> **语言**：本文档为简体中文翻译，[点此查看英文原版 README.md](README.md)。

# yq

![Build](https://github.com/mikefarah/yq/workflows/Build/badge.svg)  ![Docker Pulls](https://img.shields.io/docker/pulls/mikefarah/yq.svg) ![Github Releases (by Release)](https://img.shields.io/github/downloads/mikefarah/yq/total.svg) ![Go Report](https://goreportcard.com/badge/github.com/mikefarah/yq) ![CodeQL](https://github.com/mikefarah/yq/workflows/CodeQL/badge.svg)


一个轻量、可移植的命令行 YAML、JSON、INI 与 XML 处理器。`yq` 的语法类似 [jq](https://github.com/stedolan/jq)（一款流行的 JSON 处理器），但它同样能处理 yaml 文件以及 json、kyaml、xml、ini、properties、csv 和 tsv。它目前还不支持 `jq` 的全部功能——不过已经覆盖了最常见的操作和函数，并且新功能还在持续加入。

yq 由 Go 编写——所以你可以下载对应平台的免依赖二进制文件，直接开始使用！如果你更习惯别的工具，下面列出了可用的多种包管理器，也包括 Docker 与 Podman。

## 快速使用指南

### 基本操作

**读取一个值：**
```bash
yq '.a.b[0].c' file.yaml
```

**从 STDIN 管道读入：**
```bash
yq '.a.b[0].c' < file.yaml
```

**就地更新一个 yaml 文件：**
```bash
yq -i '.a.b[0].c = "cool"' file.yaml
```

**使用环境变量更新：**
```bash
NAME=mike yq -i '.a.b[0].c = strenv(NAME)' file.yaml
```

### 进阶操作

**合并多个文件：**
```bash
# merge two files
yq -n 'load("file1.yaml") * load("file2.yaml")'

# merge using globs (note: `ea` evaluates all files at once instead of in sequence)
yq ea '. as $item ireduce ({}; . * $item )' path/to/*.yml
```

**对一个 yaml 文件进行多处更新：**
```bash
yq -i '
  .a.b[0].c = "cool" |
  .x.y.z = "foobar" |
  .person.name = strenv(NAME)
' file.yaml
```

**在数组中查找并更新某一项：**
```bash
# Note: requires input file - add your file at the end
yq -i '(.[] | select(.name == "foo") | .address) = "12 cat st"' data.yaml
```

**在格式之间转换：**
```bash
# Convert JSON to YAML (pretty print)
yq -Poy sample.json

# Convert YAML to JSON
yq -o json file.yaml

# Convert XML to YAML
yq -o yaml file.xml
```

更多示例请见[配方与示例](https://mikefarah.gitbook.io/yq/recipes)，更多信息请见[文档](https://mikefarah.gitbook.io/yq/)。

也欢迎查看关于[常见疑问](https://github.com/mikefarah/yq/discussions/categories/q-a)以及[有趣点子](https://github.com/mikefarah/yq/discussions/categories/show-and-tell)的讨论

## 安装

### [下载最新的二进制文件](https://github.com/mikefarah/yq/releases/latest)

### wget
使用 wget 下载预编译好的二进制文件，请选择你的平台与架构：

**Linux（示例）：**
```bash
# Set your platform variables (adjust as needed)
VERSION=v4.2.0
PLATFORM=linux_amd64

# Download compressed binary
wget https://github.com/mikefarah/yq/releases/download/${VERSION}/yq_${PLATFORM}.tar.gz -O - |\
  tar xz && sudo mv yq_${PLATFORM} /usr/local/bin/yq

# Or download plain binary
wget https://github.com/mikefarah/yq/releases/download/${VERSION}/yq_${PLATFORM} -O /usr/local/bin/yq &&\
    chmod +x /usr/local/bin/yq
```

**最新版本（Linux AMD64）：**
```bash
wget https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -O /usr/local/bin/yq &&\
    chmod +x /usr/local/bin/yq
```

**可用平台：** `linux_amd64`, `linux_arm64`, `linux_arm`, `linux_386`, `darwin_amd64`, `darwin_arm64`, `windows_amd64`, `windows_386` 等。

### MacOS / Linux 通过 Homebrew：
使用 [Homebrew](https://brew.sh/)
```
brew install yq
```

### Linux 通过 snap：
```
snap install yq
```

#### Snap 注意事项
在 snap 中，`yq` 以[_严格限制（strict confinement）_](https://snapcraft.io/docs/snap-confinement)方式安装，这意味着它无法直接访问 root 拥有的系统文件。要读取这类文件，你可以：

```
sudo cat /etc/myfile | yq '.a.path'
```

而要写入这类系统文件，你可以任选其一，比如使用 [sponge](https://linux.die.net/man/1/sponge)：
```
sudo cat /etc/myfile | yq '.a.path = "value"' | sudo sponge /etc/myfile
```
或者写入一个临时文件：
```
sudo cat /etc/myfile | yq '.a.path = "value"' | sudo tee /etc/myfile.tmp
sudo mv /etc/myfile.tmp /etc/myfile
rm /etc/myfile.tmp
```

### 使用 Docker 或 Podman 运行

#### 单次使用：
```bash
# Docker - process files in current directory
docker run --rm -v "${PWD}":/workdir mikefarah/yq '.a.b[0].c' file.yaml

# Podman - same usage as Docker
podman run --rm -v "${PWD}":/workdir mikefarah/yq '.a.b[0].c' file.yaml
```

**安全提示：** 你可以在 Docker 中以受限权限运行 `yq`：
```bash
docker run --rm --security-opt=no-new-privileges --cap-drop all --network none \
  -v "${PWD}":/workdir mikefarah/yq '.a.b[0].c' file.yaml
```

#### 通过 STDIN 管道传入数据：

你需要向 Docker/Podman 传入 `-i --interactive` 标志：

```bash
# Process piped data
docker run -i --rm mikefarah/yq '.this.thing' < myfile.yml

# Same with Podman
podman run -i --rm mikefarah/yq '.this.thing' < myfile.yml
```

#### 以交互方式运行命令：

```bash
docker run --rm -it -v "${PWD}":/workdir --entrypoint sh mikefarah/yq
```

```bash
podman run --rm -it -v "${PWD}":/workdir --entrypoint sh mikefarah/yq
```

有一个 bash 函数会更省事，这样就不必输入整条 docker 命令：

```bash
yq() {
  docker run --rm -i -v "${PWD}":/workdir mikefarah/yq "$@"
}
```

```bash
yq() {
  podman run --rm -i -v "${PWD}":/workdir mikefarah/yq "$@"
}
```
#### 以 root 身份运行：

`yq` 的容器镜像不再以 root 身份运行 (https://github.com/mikefarah/yq/pull/860)。如果你想在容器镜像里安装更多内容，或者在读写文件时遇到权限问题，你需要采取以下任一方式：


```
docker run --user="root" -it --entrypoint sh mikefarah/yq
```

```
podman run --user="root" -it --entrypoint sh mikefarah/yq
```

或者，在你的 Dockerfile 中：

```
FROM mikefarah/yq

USER root
RUN apk add --no-cache bash
USER yq
```

#### 缺少时区数据
默认情况下，yq 所使用的 alpine 镜像不包含时区数据。如果你想使用 `tz` 运算符，就需要把这些数据加进去：

```
FROM mikefarah/yq

USER root
RUN apk add --no-cache tzdata
USER yq
```

#### 在 SELinux 下使用 Podman

如果你在 SELinux 环境下使用 podman，就需要在卷挂载上加上共享卷标志 `:z`：

```
-v "${PWD}":/workdir:z
```

### GitHub Action
```
  - name: Set foobar to cool
    uses: mikefarah/yq@master
    with:
      cmd: yq -i '.foo.bar = "cool"' 'config.yml'
  - name: Get an entry with a variable that might contain dots or spaces
    id: get_username
    uses: mikefarah/yq@master
    with:
      cmd: yq '.all.children.["${{ matrix.ip_address }}"].username' ops/inventories/production.yml
  - name: Reuse a variable obtained in another step
    run: echo ${{ steps.get_username.outputs.result }}
```

更多信息请见 https://mikefarah.gitbook.io/yq/usage/github-action

### Go 安装：
```
go install github.com/mikefarah/yq/v4@latest
```

## 社区支持的安装方式
这些方式由社区支持 :heart: ——不过，它们可能落后于官方支持的发布版本。

_请注意，Debian 软件包（此前由 @rmescandon 提供支持）已不再维护，请改用其他安装方式。_


### X-CMD
查看 x-cmd 上的 `yq`：https://x-cmd.com/mod/yq

- 即时结果：实时查看 yq 过滤器的输出。
- 错误处理：遇到语法错误？它会显示错误信息，并按最接近的有效过滤器给出结果

感谢 @edwinjhlee！

### Nix

```
nix profile install nixpkgs#yq-go
```

更多信息请见[这里](https://search.nixos.org/packages?channel=unstable&show=yq-go&from=0&size=50&sort=relevance&type=packages&query=yq-go)


### Webi

```
webi yq
```

更多信息请见[webi](https://webinstall.dev/)
由 @adithyasunil26 提供支持 (https://github.com/webinstall/webi-installers/tree/master/yq)

### Arch Linux

```
pacman -S go-yq
```

### Windows:

使用 [Chocolatey](https://chocolatey.org)

[![Chocolatey](https://img.shields.io/chocolatey/v/yq.svg)](https://chocolatey.org/packages/yq)
[![Chocolatey](https://img.shields.io/chocolatey/dt/yq.svg)](https://chocolatey.org/packages/yq)
```
choco install yq
```
由 @chillum 提供支持 (https://chocolatey.org/packages/yq)

使用 [scoop](https://scoop.sh/)
```
scoop install main/yq
```

使用 [winget](https://learn.microsoft.com/en-us/windows/package-manager/)
```
winget install --id MikeFarah.yq
```

### MacPorts:
使用 [MacPorts](https://www.macports.org/)
```
sudo port selfupdate
sudo port install yq
```
由 @herbygillot 提供支持 (https://ports.macports.org/maintainer/github/herbygillot)

### Alpine Linux

Alpine Linux v3.20+（以及 Edge）：
```
apk add yq-go
```

Alpine Linux v3.19 及更早版本：
```
apk add yq
```

由 Tuan Hoang 提供支持 (https://pkgs.alpinelinux.org/packages?name=yq-go)

### Flox:

Flox 可以用来在 Linux、MacOS 上安装 yq，也可以通过 WSL 在 Windows 上安装。

```
flox install yq
```


### MacOS / Linux 通过 gah：
使用 [gah](https://github.com/marverix/gah)

```
gah install yq
```

## 特性
- [包含大量示例的详细文档](https://mikefarah.gitbook.io/yq/)
- 由可移植的 Go 编写，所以你可以下载一个很棒的免依赖二进制文件
- 语法与 `jq` 相似，但同样支持处理 YAML、INI、[JSON](https://mikefarah.gitbook.io/yq/usage/convert) 与 [XML](https://mikefarah.gitbook.io/yq/usage/xml) 文件
- 完整支持多文档 yaml 文件
- 支持 yaml [front matter](https://mikefarah.gitbook.io/yq/usage/front-matter) 块（例如 jekyll/assemble）
- 带颜色的 yaml 输出
- [日期/时间的处理与格式化（含时区支持）](https://mikefarah.gitbook.io/yq/operators/datetime)
- [深层数据结构](https://mikefarah.gitbook.io/yq/operators/traverse-read)
- [键排序](https://mikefarah.gitbook.io/yq/operators/sort-keys)
- 处理 yaml 的 [注释](https://mikefarah.gitbook.io/yq/operators/comment-operators)、[样式](https://mikefarah.gitbook.io/yq/operators/style)、[标签](https://mikefarah.gitbook.io/yq/operators/tag)以及[锚点与别名](https://mikefarah.gitbook.io/yq/operators/anchor-and-alias-operators)。
- [就地更新](https://mikefarah.gitbook.io/yq/v/v4.x/commands/evaluate#flags)
- [用于选择与更新的复杂表达式](https://mikefarah.gitbook.io/yq/operators/select#select-and-update-matching-values-in-map)
- 更新时会保留 yaml 的格式与注释（不过空白字符方面仍有问题）
- [编解码 base64 数据](https://mikefarah.gitbook.io/yq/operators/encode-decode)
- [从其他文件加载内容](https://mikefarah.gitbook.io/yq/operators/load)
- [与 json/ndjson 相互转换](https://mikefarah.gitbook.io/yq/v/v4.x/usage/convert)
- [与 xml 相互转换](https://mikefarah.gitbook.io/yq/v/v4.x/usage/xml)
- [与 hcl（terraform）相互转换](https://mikefarah.gitbook.io/yq/v/v4.x/usage/hcl)
- [与 toml 相互转换](https://mikefarah.gitbook.io/yq/v/v4.x/usage/toml)
- [与 properties 相互转换](https://mikefarah.gitbook.io/yq/v/v4.x/usage/properties)
- [与 csv/tsv 相互转换](https://mikefarah.gitbook.io/yq/usage/csv-tsv)
- [通用 shell 补全脚本（bash/zsh/fish/powershell）](https://mikefarah.gitbook.io/yq/v/v4.x/commands/shell-completion)
- [Reduce](https://mikefarah.gitbook.io/yq/operators/reduce) 可用于合并多个文件、对数组求和，或其他花式用法。
- [Github Action](https://mikefarah.gitbook.io/yq/usage/github-action) 可用于你的自动化流水线（感谢 @devorbitus）

## [用法](https://mikefarah.gitbook.io/yq/)

更详细、更进阶的用法，请见[文档](https://mikefarah.gitbook.io/yq/)。

```
Usage:
  yq [flags]
  yq [command]

Examples:

# yq tries to auto-detect the file format based off the extension, and defaults to YAML if it's unknown (or piping through STDIN)
# Use the '-p/--input-format' flag to specify a format type.
cat file.xml | yq -p xml

# read the "stuff" node from "myfile.yml"
yq '.stuff' < myfile.yml

# update myfile.yml in place
yq -i '.stuff = "foo"' myfile.yml

# print contents of sample.json as idiomatic YAML
yq -P -oy sample.json


Available Commands:
  completion  Generate the autocompletion script for the specified shell
  eval        (default) Apply the expression to each document in each yaml file in sequence
  eval-all    Loads _all_ yaml documents of _all_ yaml files and runs expression once
  help        Help about any command

Flags:
  -C, --colors                          force print with colors
      --csv-auto-parse                  parse CSV YAML/JSON values (default true)
      --csv-separator char              CSV Separator character (default ,)
      --debug-node-info                 debug node info
  -e, --exit-status                     set exit status if there are no matches or null or false is returned
      --expression string               forcibly set the expression argument. Useful when yq argument detection thinks your expression is a file.
      --from-file string                Load expression from specified file.
  -f, --front-matter string             (extract|process) first input as yaml front-matter. Extract will pull out the yaml content, process will run the expression against the yaml content, leaving the remaining data intact
      --header-preprocess               Slurp any header comments and separators before processing expression. (default true)
  -h, --help                            help for yq
  -I, --indent int                      sets indent level for output (default 2)
  -i, --inplace                         update the file in place of first file given.
  -p, --input-format string             [auto|a|yaml|y|json|j|kyaml|ky|props|p|csv|c|tsv|t|xml|x|base64|uri|toml|hcl|h|lua|l|ini|i] parse format for input. (default "auto")
      --lua-globals                     output keys as top-level global variables
      --lua-prefix string               prefix (default "return ")
      --lua-suffix string               suffix (default ";\n")
      --lua-unquoted                    output unquoted string keys (e.g. {foo="bar"})
  -M, --no-colors                       force print with no colors
  -N, --no-doc                          Don't print document separators (---)
  -0, --nul-output                      Use NUL char to separate values. If unwrap scalar is also set, fail if unwrapped scalar contains NUL char.
  -n, --null-input                      Don't read input, simply evaluate the expression given. Useful for creating docs from scratch.
  -o, --output-format string            [auto|a|yaml|y|json|j|kyaml|ky|props|p|csv|c|tsv|t|xml|x|base64|uri|toml|hcl|h|shell|s|lua|l|ini|i] output format type. (default "auto")
  -P, --prettyPrint                     pretty print, shorthand for '... style = ""'
      --properties-array-brackets       use [x] in array paths (e.g. for SpringBoot)
      --properties-separator string     separator to use between keys and values (default " = ")
      --security-disable-env-ops        Disable env related operations.
      --security-disable-file-ops       Disable file related operations (e.g. load)
      --shell-key-separator string      separator for shell variable key paths (default "_")
  -s, --split-exp string                print each result (or doc) into a file named (exp). [exp] argument must return a string. You can use $index in the expression as the result counter. The necessary directories will be created.
      --split-exp-file string           Use a file to specify the split-exp expression.
      --string-interpolation            Toggles strings interpolation of \(exp) (default true)
      --tsv-auto-parse                  parse TSV YAML/JSON values (default true)
  -r, --unwrapScalar                    unwrap scalar, print the value with no quotes, colors or comments. Defaults to true for yaml (default true)
  -v, --verbose                         verbose mode
  -V, --version                         Print version information and quit
      --xml-attribute-prefix string     prefix for xml attributes (default "+@")
      --xml-content-name string         name for xml content (if no attribute name is present). (default "+content")
      --xml-directive-name string       name for xml directives (e.g. <!DOCTYPE thing cat>) (default "+directive")
      --xml-keep-namespace              enables keeping namespace after parsing attributes (default true)
      --xml-proc-inst-prefix string     prefix for xml processing instructions (e.g. <?xml version="1"?>) (default "+p_")
      --xml-raw-token                   enables using RawToken method instead Token. Commonly disables namespace translations. See https://pkg.go.dev/encoding/xml#Decoder.RawToken for details. (default true)
      --xml-skip-directives             skip over directives (e.g. <!DOCTYPE thing cat>)
      --xml-skip-proc-inst              skip over process instructions (e.g. <?xml version="1"?>)
      --xml-strict-mode                 enables strict parsing of XML. See https://pkg.go.dev/encoding/xml for more details.
      --yaml-fix-merge-anchor-to-spec   Fix merge anchor to match YAML spec. Will default to true in late 2025

Use "yq [command] --help" for more information about a command.
```

## 常见问题排查

### 常见问题

**PowerShell 引号问题：**
```powershell
# Use single quotes for expressions
yq '.a.b[0].c' file.yaml

# Or escape double quotes
yq ".a.b[0].c = \"value\"" file.yaml
```

### 获取帮助

- **查看已有 issue**：[GitHub Issues](https://github.com/mikefarah/yq/issues)
- **提问**：[GitHub Discussions](https://github.com/mikefarah/yq/discussions)
- **文档**：[完整文档](https://mikefarah.gitbook.io/yq/)
- **示例**：[配方与示例](https://mikefarah.gitbook.io/yq/recipes)

## 已知问题 / 缺失特性
- `yq` 会尽力保留注释的位置与空白字符，但并不能覆盖所有场景，详见 https://github.com/go-yaml/yaml/tree/v3
- Powershell 对如何给 yq 加引号有它自己的……[看法](https://mikefarah.gitbook.io/yq/usage/tips-and-tricks#quotes-in-windows-powershell)
- "yes"、"no" 在 yaml 1.2 标准中已不再作为布尔值——而 yq 采用的正是该标准。

更多常见问题与解决方案，请见[技巧与窍门](https://mikefarah.gitbook.io/yq/usage/tips-and-tricks)。
