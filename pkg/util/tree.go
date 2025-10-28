package util

/**
author:郭健峰
time:2024/05/27
comment:本工具为用于处理树形结构的通用工具，使用需确保原始数据每项有一个id和一个pid(父级id)
comment(en):This tool is a general tool for processing tree structure.
			It is necessary to ensure that each item of original data has an id and a pid(parent id).
*/

import (
	"reflect"
	"sort"
	"strconv"
)

type TreeNode struct {
	ID       string
	PID      string
	Value    map[string]interface{}
	Children []*TreeNode
}

func MapToTree(data []map[string]interface{}) []map[string]interface{} {
	nodeMap := make(map[string]*TreeNode)
	for _, item := range data {
		id := item["id"].(string)
		pid := item["pid"].(string)

		node := &TreeNode{
			ID:    id,
			PID:   pid,
			Value: item,
		}
		nodeMap[id] = node
	}
	root := &TreeNode{
		ID:    "0",
		PID:   "0",
		Value: nil,
	}
	// 创建tree的树形结构
	makeTree(root, nodeMap)
	var _result []map[string]interface{}
	// 创建返回值的标准格式
	_result = handleResult(root)
	return _result
}

func handleResult(root *TreeNode) []map[string]interface{} {
	var result []map[string]interface{}
	for _, child := range root.Children {
		if child.Children != nil {
			_child := handleResult(child) // 递归处理child
			child.Value["child"] = _child
		}
		result = append(result, child.Value)
	}
	return result
}

func makeTree(root *TreeNode, nodes map[string]*TreeNode) {
	for _, node := range nodes {
		if node.PID == root.ID {
			root.Children = append(root.Children, deepCopyTreeNode(node))
			delete(nodes, node.ID)
		}
	}
	for _, child := range root.Children {
		makeTree(child, nodes)
	}
}

// DeepCopyTreeNode performs a deep copy of a TreeNode.
func deepCopyTreeNode(original *TreeNode) *TreeNode {
	if original == nil {
		return nil
	}

	copyNode := &TreeNode{
		ID:    original.ID,
		PID:   original.PID,
		Value: deepCopyMap(original.Value),
	}

	for _, child := range original.Children {
		copyNode.Children = append(copyNode.Children, deepCopyTreeNode(child))
	}

	return copyNode
}

// DeepCopyMap performs a deep copy of a map[string]interface{}.
func deepCopyMap(original map[string]interface{}) map[string]interface{} {
	if original == nil {
		return nil
	}

	copyMap := make(map[string]interface{})
	for key, value := range original {
		copyMap[key] = deepCopyValue(value)
	}
	return copyMap
}

func deepCopyValue(value interface{}) interface{} {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Map:
		originalMap := value.(map[string]interface{})
		return deepCopyMap(originalMap)
	case reflect.Slice:
		originalSlice := reflect.ValueOf(value)
		copySlice := reflect.MakeSlice(originalSlice.Type(), originalSlice.Len(), originalSlice.Cap())
		for i := 0; i < originalSlice.Len(); i++ {
			copySlice.Index(i).Set(reflect.ValueOf(deepCopyValue(originalSlice.Index(i).Interface())))
		}
		return copySlice.Interface()
	default:
		return value
	}
}

func SortMenuList(list []map[string]interface{}) []map[string]interface{} {
	_tree := MapToTree(list)
	tree := HandleMenuList(_tree, 1)
	traverseObjTree(&tree)
	return tree
}

func traverseObjTree(node *[]map[string]interface{}) {
	if node == nil || len(*node) == 0 {
		return
	}

	sort.Slice(*node, func(i, j int) bool {
		sortIStr, sortIExists := (*node)[i]["sort"].(string)
		sortJStr, sortJExists := (*node)[j]["sort"].(string)
		if sortIExists && sortJExists {
			sortI, errI := strconv.Atoi(sortIStr)
			sortJ, errJ := strconv.Atoi(sortJStr)
			if errI == nil && errJ == nil {
				return sortI < sortJ
			}
		}
		return false
	})

	for i := 0; i < len(*node); i++ {
		//暂时屏蔽wx功能
		if (*node)[i]["url"] == "/home/weixin/list" {
			*node = append((*node)[:i], (*node)[i+1:]...)
			i--
			continue
		}
		if (*node)[i] != nil && (*node)[i]["child"] != nil {
			children, ok := (*node)[i]["child"].([]map[string]interface{})
			if ok && len(children) > 0 {
				traverseObjTree(&children)
				(*node)[i]["child"] = children
			}
		}
	}
}

func HandleMenuList(tree []map[string]interface{}, level int) []map[string]interface{} {
	sort.Slice(tree, func(i, j int) bool {
		a, _ := strconv.Atoi(tree[i]["id"].(string))
		b, _ := strconv.Atoi(tree[j]["id"].(string))
		return a < b
	})
	for _, leaf := range tree {
		leaf["level"] = level
		leaf["selected"] = false
		delete(leaf, "pid")
		if leaf["child"] != nil {
			HandleMenuList(leaf["child"].([]map[string]interface{}), level+1)
		}
	}
	return tree
}

func HandleRuleList(ruleList []map[string]interface{}) []string {
	tree := MapToTree(ruleList)
	result := make([]string, 0)
	for _, leaf := range tree {
		if leaf["child"] != nil {
			for _, child := range leaf["child"].([]map[string]interface{}) {
				if child["name"] != nil {
					result = append(result, "admin-"+child["name"].(string))
				}
			}
		}
	}
	return result
}
