package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mattn/anko/core"
	"github.com/mattn/anko/env"
	_ "github.com/mattn/anko/packages"
	"github.com/mattn/anko/vm"

	"github.com/carmel/gooxml/document"
	"github.com/carmel/xlsx"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v2"
)

const (
	YES = "Y"
)

type Response struct {
	Status int
	Msg    string
	Result interface{}
}

type resolver func(r *multipart.FileHeader)

func getConf(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var database []string
	if conf.Type == "mysql" {
		_ = db.Select(&database, "SHOW DATABASES")
	} else if conf.Type == "postgres" {
		_ = db.Select(&database, "SELECT datname FROM pg_database WHERE datistemplate=false")
	}

	_ = json.NewEncoder(w).Encode(Response{Status: 1, Msg: "", Result: map[string]interface{}{
		"Conf":     conf,
		"Database": database,
		"TabInfo":  tabmap,
	}})
}

func setConf(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_ = json.NewDecoder(r.Body).Decode(&conf)

	if errCheck(dbInit(), false, w, r) {
		return
	}

	f, _ := os.Create(*config)
	defer f.Close()

	yml, _ := yaml.Marshal(&conf)
	if _, err := f.Write(yml); err != nil {
		_ = json.NewEncoder(w).Encode(Response{Status: -1, Msg: err.Error()})
		return
	}
	_ = json.NewEncoder(w).Encode(Response{Status: 1})
}

func xlsxImp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	table := r.FormValue("Table")
	if table == "" {
		_ = json.NewEncoder(w).Encode(Response{Status: -1, Msg: "Please specify the table name"})
		return
	}
	var total, processed int
	mutiUpload("XLSX", w, r, func(fh *multipart.FileHeader) {

		file, err := fh.Open()
		if errCheck(err, true, w, r) {
			return
		}
		defer file.Close()

		xlsx, err := xlsx.OpenReader(file)
		if errCheck(err, true, w, r) {
			return
		}
		sn := xlsx.GetSheetName(0)
		rows, err := xlsx.GetRows(sn)
		cmt := xlsx.GetComments()[sn]

		col := make([]string, len(cmt))
		for i, c := range cmt {
			col[i] = strings.TrimPrefix(strings.TrimSpace(c.Text), "Author:")
		}

		if errCheck(err, true, w, r) {
			return
		}

		if n := len(rows); n == 0 {
			_ = json.NewEncoder(w).Encode(Response{Status: -1, Msg: fmt.Sprintf("%s is empty.", fh.Filename)})
			return
		} else {
			total += n - 1
		}
		var query string

		autoId := r.FormValue("AutoID")

		if autoId == YES {
			query = fmt.Sprintf("INSERT INTO %s (id,%s)VALUES(%s%s)", table, strings.Join(col, ","), sqlPlaceholder, strings.Repeat(","+sqlPlaceholder, len(col)))
		} else {
			query = fmt.Sprintf("INSERT INTO %s (%s)VALUES(%s)", table, strings.Join(col, ","), strings.TrimPrefix(strings.Repeat(","+sqlPlaceholder, len(col)), ","))
		}

		log.Println(query)

		sem := NewPool(conf.PoolSize, &sync.WaitGroup{})
		for i, row := range rows[1:] {
			log.Printf("[Line %d] being processed\n", i+1)
			sem.Acquire()
			go func(i int, r []string) {
				defer func() {
					sem.Release()
					if err := recover(); err != nil {
						// logger.Printf("第%d行: %+v, 错误: %v", i+1, r, err)
						log.Printf("[Line %d] error: %v\n", i+1, err)
					} else {
						processed++
					}
				}()
				var args []interface{}
				if autoId == YES {
					args = append(args, UUID())
				}
				for _, v := range r {
					args = append(args, v)
				}
				db.MustExec(query, args...)
			}(i+1, row)
		}
		sem.Wait()
	})

	_ = json.NewEncoder(w).Encode(Response{Status: 1, Msg: fmt.Sprintf("%d / %d has processed!", processed, total)})
}

func xlsxExp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 导出指定sql数据
	var sqls []string

	err := json.NewDecoder(r.Body).Decode(&sqls)
	if errCheck(err, false, w, r) {
		return
	}
	var rows *sqlx.Rows
	var index int
	var title []string
	xlsx := xlsx.NewFile()
	for i, sql := range sqls {
		index = 0
		rows, err = db.Queryx(sql)
		if errCheck(err, false, w, r) {
			return
		}

		if i != 0 {
			xlsx.NewSheet(fmt.Sprintf("%s%d", "Sheet", i+1))
		}

		for rows.Next() {
			index++
			if index == 1 {
				title, _ = rows.Columns()
				for n, v := range title {
					if err = xlsx.SetCellValue(fmt.Sprintf("%s%d", "Sheet", i+1), fmt.Sprintf("%s%d", EXCEL_COL[n], 1), v); err != nil {
						log.Println(`SetCellValue`, err)
					}
				}
			}
			rs, _ := rows.SliceScan()
			for n, v := range rs {
				// fmt.Println(fmt.Sprintf("%s%d", "Sheet", i+1), fmt.Sprintf("%s%d", EXCEL_COL[n+1], index))
				if err = xlsx.SetCellValue(fmt.Sprintf("%s%d", "Sheet", i+1), fmt.Sprintf("%s%d", EXCEL_COL[n], index+1), v); err != nil {
					log.Println(`SetCellValue`, err)
				}
			}
		}
	}

	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%d.xlsx", time.Now().Unix()))
	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Add("Content-Transfer-Encoding", "binary")
	err = xlsx.Write(w)
	if errCheck(err, false, w, r) {
		return
	}
}

type Row map[string]interface{}

func docxImp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	table := r.FormValue("Table")
	if table == "" {
		_ = json.NewEncoder(w).Encode(Response{Status: -1, Msg: "Please specify the table name"})
		return
	}

	script := r.FormValue("AnkoScript")
	if script == "" {
		_ = json.NewEncoder(w).Encode(Response{Status: -1, Msg: "Please specify the ank script"})
		return
	}

	e := env.NewEnv()
	core.Import(e)
	_ = e.DefineType("Row", Row{})
	_ = e.Define("AutoID", r.FormValue("AutoID"))
	_ = e.Define("Tab", table)
	_ = e.Define("DB", db)
	_ = e.Define("JsonEncode", EncodeString)
	_ = e.Define("UUID", UUID)

	mutiUpload("DOCX", w, r, func(fh *multipart.FileHeader) {
		file, err := fh.Open()
		if errCheck(err, false, w, r) {
			return
		}
		defer file.Close()
		var doc *document.Document
		doc, err = document.Read(file, fh.Size)
		if errCheck(err, false, w, r) {
			return
		}

		_ = e.Define("PS", doc.Paragraphs())
		_ = e.Define("Doc", fh.Filename)

		_, err = vm.Execute(e, nil, script)
		if errCheck(err, true, w, r) {
			return
		}
	})
	_ = json.NewEncoder(w).Encode(Response{Status: 1, Msg: "Request completed!"})
}

func errCheck(err error, print bool, w http.ResponseWriter, r *http.Request) bool {
	if err != nil {
		if print {
			log.Println(err.Error())
		} else {
			_ = json.NewEncoder(w).Encode(Response{Status: -1, Msg: err.Error()})
		}
		return true
	}
	return false
}
func mutiUpload(formKey string, w http.ResponseWriter, r *http.Request, fn resolver) {
	_ = r.ParseMultipartForm(32 << 20)
	if files, ok := r.MultipartForm.File[formKey]; ok {
		for _, fh := range files {
			fn(fh)
		}
	}
}

func excelExport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	tabname := mux.Vars(r)["table"]
	table := tabmap[tabname]
	filename := table.Comment
	if filename == "" {
		filename = tabname
	}

	var axis string
	var err error
	xlsx := xlsx.NewFile()
	for i, c := range table.Col {
		axis = EXCEL_COL[i] + "1"

		_ = xlsx.AddComment("Sheet1", axis, "", c.Name)
		if c.Comment == "" {
			_ = xlsx.SetCellValue("Sheet1", axis, c.Name)
		} else {
			_ = xlsx.SetCellValue("Sheet1", axis, c.Comment)
		}
	}

	err = excelTransfer(filename, xlsx, w)
	errCheck(err, false, w, r)
}

func excelTransfer(filename string, xlsx *xlsx.File, w http.ResponseWriter) error {
	// w.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s.xlsx", url.QueryEscape(filename)))
	w.Header().Add("Content-Type", "application/octet-stream")
	// w.Header().Set("Content-Transfer-Encoding", "binary")
	return xlsx.Write(w)
}
