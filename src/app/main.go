//go:generate go run -tags generate gen.go

package main

import (
	"fmt"
	"github.com/zserge/lorca"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
)

// Go types that are bound to the UI must be thread-safe, because each binding
// is executed in its own goroutine. In this simple case we may use atomic
// operations, but for more complex cases one should use proper synchronization.
type counter struct {
	sync.Mutex
	count int
}

func (c *counter) Add(n int) {
	c.Lock()
	defer c.Unlock()
	c.count = c.count + n
}

func (c *counter) Value() int {
	c.Lock()
	defer c.Unlock()
	return c.count
}

type Request map[string]interface{}
type Object map[string]interface{}

var projects []Object
var projectDir = "C:/Users/Default/Documents"
var projectFileName = "game-config-project.list"

type ConfigProject struct {
	name            string
	id              string
	source          string
	dataOutDir      string
	dataFileName    string
	defineOutDir    string
	defineFileName  string
	defineNamespace string
}

func main() {
	var args []string
	if runtime.GOOS == "linux" {
		args = append(args, "--class=Lorca")
	}
	ui, err := lorca.New("", "", 720, 480, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()
	ui.Bind("loadData", loadData)
	ui.Bind("exportExcel", exportExcel)
	ui.Bind("addProject", addProject)
	ui.Bind("modifyProject", modifyProject)
	ui.Bind("deleteProject", deleteProject)
	projects = make([]Object, 0, 5)
	loadData()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	go http.Serve(ln, http.FileServer(FS))
	_ = ui.Load(fmt.Sprintf("http://%s", ln.Addr()))

	//ui.Eval(`
	//	console.log("Hello, world!");
	//	console.log('Multiple values:', [1, false, {"x":5}]);
	//`)
	ui.Eval(fmt.Sprintf(`console.log("listener.addr: %s");`, ln.Addr()))

	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}
}

//返回一个对象数组
func loadData() []Object {
	var data []Object
	f, _ := ioutil.ReadFile(fmt.Sprintf("%s/%s", projectDir, projectFileName))
	_ = ParseMap(f, &data)
	projects = data
	return data
}

//type ConfigProject struct {
//	name            string
//	id              string
//	source          string
//	dataOutDir      string
//	dataFileName    string
//	defineOutDir    string
//	defineFileName  string
//	defineNamespace string
//}
func exportExcel(r Request) Object {
	log.Println("\nexportExcel: ", r)
	if r == nil {
		return Object{"code": 102, "msg": "修改项目失败, 项目参数为空值"}
	}
	e := Run(r["source"].(string), r["dataOutDir"].(string),
		r["dataFileName"].(string), false, r["defineOutDir"].(string),
		r["defineFileName"].(string), r["defineNamespace"].(string))
	if e != nil {
		return Object{"code": 103, "msg": "项目导出失败, " + e.Error()}
	}
	return Object{"code": 0, "msg": "项目导出成功", "data": r}
}

//新建一个项目
//在项目目录下创建[id].gt文件存储项目信息
func addProject(r Request) Object {
	log.Println("\naddProject: ", r)
	if r == nil {
		return Object{"code": 102, "msg": "修改项目失败, 项目参数为空值"}
	}

	appendOrUpdate(Object(r))
	err := updateFile()
	if err != nil {
		return Object{"code": 100, "msg": "创建项目失败," + err.Error()}
	}
	return Object{"code": 0, "msg": "创建项目成功"}
}
func modifyProject(r Request) Object {
	log.Println("\nmodifyProject: ", r)
	if r == nil {
		return Object{"code": 102, "msg": "修改项目失败, 项目参数为空值"}
	}

	appendOrUpdate(Object(r))
	err := updateFile()
	if err != nil {
		return Object{"code": 100, "msg": "修改项目失败," + err.Error()}
	}
	return Object{"code": 0, "msg": "修改项目成功"}
}

func deleteProject(id interface{}) Object {
	log.Println("\ndeleteProject: ", id)

	remove(id)
	err := updateFile()
	if err != nil {
		return Object{"code": 100, "msg": "删除项目失败," + err.Error()}
	}
	return Object{"code": 0, "msg": "删除项目成功"}
}

func updateFile() error {
	EnsureDir(projectDir)
	file, bPrj := GetFillAndBufferWriter(fmt.Sprintf("%s/%s", projectDir, projectFileName))
	_, err := bPrj.WriteString(StringifyMap(projects))
	_ = bPrj.Flush()
	file.Close()
	return err
}

func appendOrUpdate(r Object) {
	for i, project := range projects {
		if project["id"] == r["id"] {
			projects[i] = r
			return
		}
	}
	projects = append(projects, r)
}
func remove(id interface{}) {
	for i, project := range projects {
		if project["id"] == id {
			projects = append(projects[:i], projects[i+1:]...)
		}
	}
}
