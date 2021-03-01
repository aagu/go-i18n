# go-i18n

## 是什么

一个为项目添加本地化和国际化支持的golang库，该库很大程度上受[nicksnyder/go-i18n](https://github.com/nicksnyder/go-i18n)的启发

## 怎么用

在项目中引入依赖，并声明一些`Message`对象，`ID`必须唯一，`Text`中写入需要翻译的字符串，支持go template语法

```go
var (
	jan = translation.Message{ID: "January", Text: "January"}
	greeting = translation.Message{ID: "Hello", Text: "Hello"}
	format = translation.Message{ID: "DayOfMonth", Text: "The {{.Day}}(th) day of {{.Month}}"}
	ways = translation.Message{ID: "TwoWay", Text: "One way is to {{.One}}, the other is to {{.Other}}"}
	roma = translation.Message{ID: "roma", Text: "roma"}
	paris = translation.Message{ID: "paris", Text: "paris"}
	text = translation.Message{ID: "text", Text: "{{.}}"}
)
```

使用go-i18n extract 命令生成原语言的json文件
```shell
go-i18n extract -o ./i18n main.go
```

准备若干种语言对应的翻译，文件名需要符合`golang.org/x/text/language`中对语言的定义，并且以.json结尾

en.json
```json
[
  {
    "id":"TwoWay",
    "text":"One way is to {{.One}}, the other is to {{.Other}}"
  },
  {
    "id":"January",
    "text":"January"
  },
  {
    "id":"Hello",
    "text":"Hello"
  },
  {
    "id":"roma",
    "text":"roma"
  },
  {
    "id":"paris",
    "text":"paris"
  },
  {
    "id":"DayOfMonth",
    "text":"The {{.Day}}(th) day of {{.Month}}"
  },
  {
    "id":"text",
    "text":"{{.}}"
  }
]
```

对于每个需要翻译的条目，id是唯一的，text内的文字支持go-template模板语法。翻译完成后将所有文件放到一个文件夹中。

在main函数中加载翻译
```go
translation.LoadTranslations(`C:\Users\aagui\IdeaProjects\go-i18n\i18n`)
```

可以设置默认语言，这样`Message`在调用`String`方法时便会根据默认语言输出文本。
```go
translation.SetDefaultLocale(lang.English)
```

也可以用`Translate`方法指定输出的语言。为`Message`对象设置的输出语言，会传递给其所引用的`Message`对象。否则，默认语言会被传递。
```go
fmt.Println(greeting.Translate(lang.SimplifiedChinese))
fmt.Println(format.FormatTranslate(lang.SimplifiedChinese, translation.TemplateData{"Day": 2, "Month": jan}))
fmt.Println(ways.Format(translation.TemplateData{"One": roma, "Other": paris}))
```