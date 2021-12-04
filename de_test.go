package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/axgle/mahonia"
	"github.com/carmel/gooxml/document"
	"github.com/gorilla/mux"
	"github.com/mattn/anko/env"
	"github.com/mattn/anko/vm"
	"golang.org/x/text/encoding/unicode"
	"gopkg.in/yaml.v2"
)

func TestServer(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()
	c, err := ioutil.ReadFile(*config)
	if err != nil {
		log.Fatalln(err)
	}

	err = yaml.Unmarshal(c, &conf)
	if err != nil {
		log.Fatalln(err)
	}

	if err = dbInit(); err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	r := mux.NewRouter().StrictSlash(true)

	//初始化路由
	r.HandleFunc("/xlsximp", xlsxImp).Methods("POST")
	r.HandleFunc("/app", func(w http.ResponseWriter, rq *http.Request) {
		http.Redirect(w, rq, "/web/index.htm", http.StatusMovedPermanently)
	}).Methods("GET")
	r.HandleFunc("/docximp", docxImp).Methods("POST")
	r.HandleFunc("/xlsxexp", xlsxExp).Methods("POST")
	r.HandleFunc("/getconf", getConf).Methods("GET")
	r.HandleFunc("/exporttemplate/{table}", excelExport).Methods("POST")
	r.HandleFunc("/setconf", setConf).Methods("POST")

	r.Use(mux.CORSMethodMiddleware(r))

	r.PathPrefix("/web").Handler(http.StripPrefix("/", http.FileServer(http.FS(web))))
	err = http.ListenAndServe(":7788", r)
	if err != nil {
		println(err.Error())
	}
}

func TestMultipleCSV(t *testing.T) {
	defer db.Close()
	var record []string
	var val []interface{}
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	if err := filepath.Walk("C:/Users/Vector/Desktop/R", func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fi.IsDir() && strings.HasSuffix(fi.Name(), ".csv") {
			f, _ := os.Open(path)
			r := csv.NewReader(decoder.Reader(f))

			for {
				val = val[0:0]
				record, err = r.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					fmt.Println("Error:", err)
					break
				}
				var c, v []string
				for i, rc := range strings.Fields(record[0]) {
					c = append(c, "c"+strconv.Itoa(i+1))
					v = append(v, "?")
					val = append(val, rc)
				}
				sql := fmt.Sprintf("INSERT INTO attendance(%s) VALUES (%s)", strings.Join(c, ","), strings.Join(v, ","))
				fmt.Println(sql, val)
				db.MustExec(sql, val...)
			}
		}
		return nil
	}); err != nil {
		fmt.Println(err)
	}
}

func ConvertToString(src string, srcCode string, tagCode string) (string, error) {
	if len(src) == 0 || len(srcCode) == 0 || len(tagCode) == 0 {
		return "", errors.New("input arguments error")
	}
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)

	return result, nil
}

func TestTrim(t *testing.T) {
	fmt.Println(strings.TrimPrefix(strings.TrimSpace("Author:carmel"), "Author:"))
}

func TestAnko(t *testing.T) {
	col := []map[string]interface{}{
		{"c1": 1, "c2": "c2"},
	}
	e := env.NewEnv()
	_ = e.DefineType("result", map[string]interface{}{})
	_ = e.Define("col", col)
	_ = e.Define("println", log.Println)

	value, err := vm.Execute(e, nil, `
res = make(result)
res["name"] = "张三"
res["age"] = 12
println(res)
col+=res
`)
	if err == nil {
		fmt.Printf("Return value type: %v\n", reflect.TypeOf(value))
		if res, ok := value.([]map[string]interface{}); ok {
			fmt.Println(res)
		} else {
			fmt.Println("Unexcept type of result")
		}
		fmt.Println(col)
	} else {
		fmt.Println(err)
	}

	_ = e.Define("col", value)
	value, err = vm.Execute(e, nil, `
	res = make(result)
	res["name"] = "张三"
	res["age"] = 12
	col+=res
	`)
	if err == nil {
		fmt.Printf("Return value type: %v\n", reflect.TypeOf(value))
		if res, ok := value.([]map[string]interface{}); ok {
			fmt.Println(res)
		} else {
			fmt.Println("Unexcept type of result")
		}
	} else {
		fmt.Println(err)
	}
}

func TestDocxImp(t *testing.T) {
	doc, err := document.Open("qust.docx")
	if err != nil {
		log.Fatalf("error opening document: %s", err)
	}

	raw, _ := ioutil.ReadFile("docx.ank")
	script := string(raw)
	e := env.NewEnv()
	_ = e.DefineType("Row", Row{})
	_ = e.Define("println", log.Println)
	_ = e.Define("AutoID", YES)
	_ = e.Define("Tab", "qust")
	_ = e.Define("JsonEncode", EncodeString)
	_ = e.Define("UUID", UUID)

	_ = e.Define("PS", doc.Paragraphs())
	_ = e.Define("Doc", "demo_qust")

	_, err = vm.Execute(e, nil, script)
	if err != nil {
		log.Println(err)
	}
}

func TestRegex(t *testing.T) {
	text := `122．处理客户投诉包含( ABCDE )等方法。
  A、回复投诉
  B、记录投诉
  C、判断投诉
  D、分析投诉
  E、反馈投诉
123．网络采购核心因素是管理者对于成本有( ABCD )。
  A、减少开支
  B、自动化
  C、用成本较低的技术替代昂贵的软件
  D、降低维护、支持、能耗成本
  E、预估
`
	stemRe := regexp.MustCompile(`^(\d{1,3})．`)
	var re map[string]*regexp.Regexp = map[string]*regexp.Regexp{
		"开始": regexp.MustCompile(`^单项选择题`),
		"选项": regexp.MustCompile(` *[A-E]、(\S+)\s*`),
		"答案": regexp.MustCompile(`\( *([A-E]{1,5}) *\)`),
	}
	replaceText := "(   )"

	fmt.Println("是否匹配：", stemRe.MatchString(text))
	stem := strings.SplitN(text, "\n", 2)[0]
	fmt.Printf("题干：%q\n", stem)
	fmt.Println("序号：", stemRe.FindStringSubmatch(text)[1])
	var option []string
	for _, o := range re["选项"].FindAllStringSubmatch(text, -1) {
		option = append(option, string(o[1]))
	}
	fmt.Println("选项：", EncodeString(option))
	var answer []string
	for _, o := range re["答案"].FindAllStringSubmatch(text, -1) {
		answer = append(answer, string(o[1]))
	}
	fmt.Println("答案：", EncodeString(answer))
	fmt.Println("擦除答案：", re["答案"].ReplaceAllString(stem, replaceText))

}

func TestJson(t *testing.T) {
	js := []string{"a", "b", "1", "2"}
	str := EncodeString(js)
	fmt.Println(str)
	var des []string
	fmt.Println(Decode(str, &des), des)
}
