package pathtrie

// Inspired by https://github.com/akitasoftware/akita-libs/blob/main/path_trie/path_trie.go

import (
	"container/list"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ludovic-pourrat/kong-api-catalog-harvester/utils"
	"reflect"
	"strings"
	"sync"
)

var mu sync.Mutex

type TrieNode struct {
	Children PathToTrieNode

	// Name of the path segment corresponding to this node.
	// E.g. if this node represents /v1/foo/bar,
	// the Name would be "bar" and the Path would be "/v1/foo/bar".
	Name string

	// Path includes the node's name and uniquely identifies the node in the tree.
	Path string

	// URL as it was.
	URL string

	// Value of the full path
	Value interface{}

	// Operation
	Operations map[string]*openapi3.Operation
}

type PathToTrieNode map[string]*TrieNode

type PathTrie struct {
	Trie          PathToTrieNode
	PathSeparator string
}

type ValueMergeFunc func(existing, newV *interface{})

func IsNil(a interface{}) bool {
	return a == nil || (reflect.ValueOf(a).Kind() == reflect.Ptr && reflect.ValueOf(a).IsNil())
}

func New() PathTrie {
	return PathTrie{
		Trie:          make(PathToTrieNode),
		PathSeparator: "/",
	}
}

// Insert val at path, with path segments separated by PathSeparator.
// Returns true if a new path was created, false if an existing path
// was overwritten.
func (pt *PathTrie) Insert(computed string,
	url string,
	operation *openapi3.Operation,
	method string,
	val interface{}) bool {
	isNewPath := pt.InsertMerge(strings.Split(computed, pt.PathSeparator),
		strings.Split(url, pt.PathSeparator),
		map[string]*openapi3.Operation{method: operation},
		val, func(existing, newV *interface{}) {
			*existing = *newV
		})
	s := pt.Print()
	fmt.Println(s)
	return isNewPath
}

// Insert val at path, with path segments separated by PathSeparator.
// Returns true if a new path was created, false if an existing path
// was overwritten.
//
// The merge function is responsible for updating the existing value
// with the new value.
func (pt *PathTrie) InsertMerge(segments []string,
	urls []string,
	operations map[string]*openapi3.Operation,
	val interface{},
	merge ValueMergeFunc) (isNewPath bool) {

	trie := pt.Trie
	isNewPath = true
	// TODO: what about computed that ends with pt.PathSeparator is it different ?

	// Traverse the Trie along computed, inserting nodes where necessary.
	for idx, segment := range segments {
		isLastSegment := idx == len(segments)-1

		if node, ok := trie[segment]; ok {
			if len(urls) == idx+1 {
				for k, v := range operations {
					node.Operations[k] = v
				}
			}
			if isLastSegment {
				// If this is the last computed segment, then this is the node to update.
				// If node value is not empty it means that an existing computed is overwritten
				isNewPath = IsNil(node.Value)
				merge(&node.Value, &val)
			} else {
				trie = node.Children
			}
		} else {
			var children []*TrieNode
			if len(trie) >= 2 {
				for k, v := range trie {
					for _, child := range v.Children {
						children = append(children, child)
					}
					delete(trie, k)
				}
				// TODO merge query params
				if !utils.IsPathParam(segment) {
					segment = utils.GenerateParamName(idx, segments)
					segments[idx] = segment
					for k, _ := range operations {
						operations[k].OperationID = utils.GetName(k, segments)
					}
				}
				// TODO update (computed, url, id) in parents
			}
			newOperations := operations
			if len(urls) != idx+1 {
				newOperations = make(map[string]*openapi3.Operation)
			}
			newNode := pt.createPathTrieNode(newOperations,
				segments[:idx+1],
				urls[:idx+1],
				idx,
				isLastSegment,
				val)
			trie[segment] = newNode
			// merge children
			for _, child := range children {
				paths := strings.Split(child.Path, pt.PathSeparator)
				paths[idx] = segment
				for k, v := range child.Operations {
					v.OperationID = utils.GetName(k, paths)
				}
				child.Path = strings.Join(paths, "/")
				pt.InsertMerge(paths,
					strings.Split(child.URL, pt.PathSeparator),
					child.Operations,
					val,
					merge)
			}
			// continue descending.
			trie = newNode.Children

		}
	}

	return isNewPath
}

func (pt *PathTrie) createPathTrieNode(operations map[string]*openapi3.Operation,
	segments []string,
	urls []string,
	idx int,
	isLastSegment bool,
	val interface{}) *TrieNode {
	//fullPathSegments := segments[:idx+1]
	node := &TrieNode{
		Children:   make(PathToTrieNode),
		Name:       segments[idx],
		Path:       strings.Join(segments, pt.PathSeparator),
		Operations: operations,
		URL:        strings.Join(urls, pt.PathSeparator),
	}
	if isLastSegment {
		node.Value = val
	}

	return node
}

// Nodes returns a list of all graph nodes
func (pt *PathTrie) Nodes() []*TrieNode {
	//track the visited nodes
	var nodes []*TrieNode
	// queue of the nodes to visit
	queue := list.New()
	// add the root node to the map of the visited nodes
	root := pt.Trie[""]
	nodes = append(nodes, root)
	queue.PushBack(root)
	for queue.Len() > 0 {
		qnode := queue.Front()
		// iterate through all of its friends
		// mark the visited nodes; enqueue the non-visted
		for _, node := range qnode.Value.(*TrieNode).Children {
			nodes = append(nodes, node)
			queue.PushBack(node)
		}
		queue.Remove(qnode)
	}
	return nodes
}

func (pt *PathTrie) Print() string {
	s := ""
	for _, node := range pt.Nodes() {
		s += node.Name + "\n"
	}
	return s
}
