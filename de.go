package main

import (
	"embed"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v2"
)

//go:embed web/*
var web embed.FS

type ColInfo struct {
	Name    string `db:"COLUMN_NAME"`
	Comment string `db:"COLUMN_COMMENT"`
}

type TabInfo struct {
	Comment string
	Col     []ColInfo
}

var (
	db             *sqlx.DB
	sqlPlaceholder string
	// lock sync.Mutex
	tabmap map[string]TabInfo
	conf   struct {
		ConnectionStr string `yaml:"ConnectionStr"`
		Type          string `yaml:"Type"`
		DB            string `yaml:"DB"`
		PoolSize      int    `yaml:"PoolSize"`
	}
	EXCEL_COL = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	// path      = flag.String("p", "/Users/carmel/Desktop/4311.xlsx", "path of excel file")
	// generate  = flag.String("g", "", "automatically generate columns, separated with comma")
	config = flag.String("c", "conf.yml", "configuration path")
)

func dbInit() error {
	var err error
	db, err = sqlx.Connect(conf.Type, conf.ConnectionStr)
	if err != nil {
		return err
	}

	switch conf.Type {
	case "postgres", "pgx", "pq-timeouts", "cloudsqlpostgres", "ql":
		sqlPlaceholder = "$"
	case "mysql", "sqlite3":
		sqlPlaceholder = "?"
	case "oci8", "ora", "goracle", "godror":
		sqlPlaceholder = ":"
	case "sqlserver":
		sqlPlaceholder = "@"
	}

	tabmap = make(map[string]TabInfo)
	var tables []struct {
		Name    string `db:"TABLE_NAME"`
		Comment string `db:"TABLE_COMMENT"`
	}
	if conf.Type == "mysql" {
		db.MustExec(fmt.Sprintf("use %s", conf.DB))

		_ = db.Select(&tables, fmt.Sprintf("SELECT table_name,table_comment FROM information_schema.tables WHERE table_schema='%s'", conf.DB))

		for _, tab := range tables {
			var cols []ColInfo
			_ = db.Select(&cols, fmt.Sprintf("SELECT column_name,column_comment FROM information_schema.columns WHERE table_schema='%s' and table_name='%s'", conf.DB, tab.Name))
			tabmap[tab.Name] = TabInfo{tab.Comment, cols}
		}
	} else {
		_ = db.Select(&tables, "SELECT name FROM sqlite_master WHERE type='table'")

		for _, tab := range tables {
			var cols []ColInfo
			_ = db.Select(&cols, fmt.Sprintf("SELECT name as COLUMN_COMMENT FROM pragma_table_info('%s')", tab.Name))
			tabmap[tab.Name] = TabInfo{tab.Comment, cols}
		}
	}

	return err
}

func main() {
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

	// 自动打开浏览器
	var cmd *exec.Cmd
	url := "http://localhost:7788/app/"
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	}
	if cmd != nil {
		go func() {
			println("Open the default browser after two seconds...")
			println("You can also open " + url + " manually.")

			time.Sleep(time.Second * 2)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			_ = cmd.Run()
		}()
	}

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

	r.PathPrefix("/web/").Handler(http.StripPrefix("/", http.FileServer(http.FS(web))))
	err = http.ListenAndServe(":7788", r)
	if err != nil {
		println(err.Error())
	}

}
