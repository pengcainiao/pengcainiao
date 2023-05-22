package middleware

import (
	"context"
	"fmt"
	"github.com/pengcainiao/zero/core/logx"
	"net/http"
	"net/url"
	"strings"
)

var (
	tireRoute = newRouter()
)

type node struct {
	path     string           // 路由路径
	part     string           // 路由中由'/'分隔的部分， 比如路由/hello/:name，那么part就是hello和:name
	children map[string]*node // 子节点
	isWild   bool             // 是否精确匹配，true代表当前节点是通配符，模糊匹配
}

type router struct {
	root map[string]*node // 路由树根节点，每个http方法都是一个路由树
}

func newRouter() *router {
	r := &router{root: make(map[string]*node)}
	r.init()
	return r
}

func (r *router) init() {
	r.addRoute(http.MethodPost, "/:version/auth/pc")             //Deprecated
}


func (n *node) String() string {
	if n == nil {
		return "NOT MATCHED \n"
	}
	return fmt.Sprintf("node{path=%s, part=%s, isWild=%t} \n", n.path, n.part, n.isWild)
}

// parsePath 分隔路由为part字典
// 比如路由/home/:name将被分隔为["home", ":name"]
func parsePath(path string) (parts []string) {
	// 将path以"/"分隔为parts
	var par []string
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		u, _ := url.Parse(path)
		par = strings.Split(u.Path, "/")
		//if par[0] == "" {
		//	par[0] = u.Host
		//}

	} else {
		par = strings.Split(path, "/")
	}
	for _, p := range par {
		if p != "" {
			parts = append(parts, p)
			// 如果part是以通配符*开头的
			if p[0] == '*' {
				break
			}
		}
	}
	return
}

// addRoute 绑定路由到handler
// g.Get() 会调用addRoute方法将path添加到路由树上面
func (r *router) addRoute(method, path string) {
	parts := parsePath(path)
	if _, ok := r.root[method]; !ok {
		r.root[method] = &node{children: make(map[string]*node)}
	}
	root := r.root[method]
	// 将parts插入到路由树
	var tempPath string
	for i := 0; i < len(parts); i++ {
		if i == len(parts)-1 {
			tempPath = path
		}
		var part = parts[i]
		if root.children[part] == nil {
			root.children[part] = &node{
				part:     part,
				path:     tempPath,
				children: make(map[string]*node),
				isWild:   part[0] == ':' || part[0] == '*'}
		}
		root = root.children[part]
	}
	root.path = path
}

// getRoute 获取路由树节点以及路由变量
// method用来判断属于哪一个方法路由树，path用来获取路由树节点和参数
func (r *router) getRoute(method, path string) *node {
	searchParts := parsePath(path)
	if len(searchParts) == 0 {
		return nil
	}

	var (
		ok          bool
		matchedNode *node
	)
	if matchedNode, ok = r.root[method]; !ok {
		logx.NewTraceLogger(context.Background()).Debug().Str("method", method).Msg("查询不到指定方法")
		return nil
	}
	// 在该方法的路由树上查找该路径
	// 查找child是否等于part
	for _, part := range searchParts {
		for _, child := range matchedNode.children {
			var fullPath string
			if child.part == part || child.isWild {
				if len(searchParts) > 0 {
					fullPath += "/" + child.part
					if n := r.search(child, searchParts[1:], fullPath); n != nil {
						return n
					}
				} else {
					return child
				}
			}
		}
	}
	return nil
}

func (r *router) search(child *node, part []string, fullPath string) *node {
	if child == nil {
		return nil
	}
	if child.path == fullPath && len(part) == 0 {
		return child
	}
	var tempPath = fullPath
	for _, n := range child.children {
		if n.path == fullPath {
			return child
		}
		if len(part) > 0 {
			if n.part == part[0] || n.isWild {
				// 可继续向下
				fullPath = tempPath + "/" + n.part
				if v := r.search(n, part[1:], fullPath); v != nil {
					return v
				}
			}
		}
	}
	return nil
}

// Handle 用来绑定路由和handlerFunc
func (r *router) Handle(method, path string) bool {
	if path == "" {
		return false
	}
	// 获取路由树节点和动态路由中的参数
	node := r.getRoute(method, path)
	return node != nil
}
