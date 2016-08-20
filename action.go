package apimock

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultActionFileName = "actions.json"
)

type ActionManager interface {
	DoAction(*http.Request) (*MockResponse, error)
	CheckActions()
}

type actionManager struct {
	root   string
	depth  int
	uriMap map[string]string
	cache  map[string]map[string]*actionDefinition
}

type MockResponse struct {
	Status int
	Body   []byte
}

type actionDefinition struct {
	Status int
	Body   interface{}
}

func NewActionManager(root string) (ActionManager, error) {

	absPath, _ := filepath.Abs(root)
	parts := strings.Split(absPath, "/")
	depth := len(parts)
	actionManager := &actionManager{}
	actionManager.root = absPath
	actionManager.depth = depth
	actionManager.uriMap = make(map[string]string)
	actionManager.cache = make(map[string]map[string]*actionDefinition)
	filepath.Walk(absPath, actionManager.loadFile)

	return actionManager, nil
}

func (mgr *actionManager) loadFile(path string, info os.FileInfo, err error) error {

	if info == nil {
		return errors.New(path + " is not found.")
	}

	if info.IsDir() {
		return nil
	}

	if info.Name() != defaultActionFileName {
		fmt.Println("this is not action file.:" + path)
		return errors.New("this is not action file.:" + path)
	}

	// split URI Path from file path
	parts := strings.Split(path, "/")
	subParts := parts[mgr.depth : len(parts)-1]
	subPath := "/" + strings.Join(subParts, "/")
	mgr.uriMap[strings.ToLower(subPath)] = path

	return nil
}

func (mgr *actionManager) getActions(filePath string) (map[string]*actionDefinition, error) {
	actions, ok := mgr.cache[filePath]
	if ok {
		return actions, nil
	}

	jsonValue, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data map[string]map[string]interface{}
	if err := json.Unmarshal(jsonValue, &data); err != nil {
		panic(err)
		return nil, err
	}

	actions = make(map[string]*actionDefinition)
	for httpMethod, action := range data {
		status, ok := action["status"].(float64)
		if !ok {
			status = -1
		}

		body, ok := action["body"].(interface{})
		if !ok {
			body = nil
		}

		actions[httpMethod] = &actionDefinition{
			Status: int(status),
			Body:   body,
		}
	}

	mgr.cache[filePath] = actions
	return actions, nil
}

func (mgr *actionManager) parseActions(actionFile, httpMethod string) (*MockResponse, error) {
	mock := &MockResponse{
		Status: http.StatusOK,
	}

	actions, err := mgr.getActions(actionFile)
	if err != nil {
		return nil, err
	}

	action, ok := actions[httpMethod]
	if !ok {
		mock.Status = http.StatusMethodNotAllowed
		return mock, nil
	}

	mock.Status = action.Status
	if action.Body == nil {
		mock.Body = []byte("")
	} else {
		body, err := json.Marshal(&action.Body)
		if err != nil {
			return nil, err
		}
		mock.Body = body
	}

	return mock, nil
}

func (mgr *actionManager) DoAction(r *http.Request) (*MockResponse, error) {

	path := r.URL.Path
	// pathをparseしてActionManagerからレスポンスを受け取る
	actionFile, ok := mgr.uriMap[strings.ToLower(path)]
	// NotFound
	if !ok {
		mock := &MockResponse{
			Status: http.StatusNotFound,
		}
		return mock, nil
	}

	mock, err := mgr.parseActions(actionFile, r.Method)
	if err != nil {
		return nil, err
	}
	return mock, nil
}

func (mgr *actionManager) CheckActions() {

	for path, file := range mgr.uriMap {
		fmt.Println("Path:", path)
		actions, err := mgr.getActions(file)
		if err != nil {
			fmt.Printf("  Failed to parse file¥n")
			continue
		}
		for k, v := range actions {
			fmt.Printf("  %v=>\n", k)
			fmt.Printf("    Status:%v\n", v.Status)
			jsonValue, _ := json.Marshal(v.Body)
			fmt.Printf("    Body:%v\n", string(jsonValue))
		}
	}
}
