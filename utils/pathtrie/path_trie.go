package pathtrie

// Inspired by https://github.com/akitasoftware/akita-libs/blob/main/path_trie/path_trie.go

import (
	"container/list"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ludovic-pourrat/kong-api-catalog-harvester/utils"
	"reflect"
	"strings"
)

type TrieNode struct {
	Children PathToTrieNode

	// Name of the path segment corresponding to this node.
	// E.g. if this node represents /v1/foo/bar,
	// the Name would be "bar" and the Path would be "/v1/foo/bar".
	Name string

	// Unique Id
	Id string

	// Path includes the node's name and uniquely identifies the node in the tree.
	Path string

	// URL as it was.
	URL string

	// Value of the full path
	Value interface{}

	// Operation
	Operation *openapi3.Operation

	// Method
	Method string
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
	id string,
	operation *openapi3.Operation,
	method string,
	val interface{}) bool {
	return pt.InsertMerge(computed,
		url,
		id,
		operation,
		method,
		val, func(existing, newV *interface{}) {
			*existing = *newV
		})
}

// Insert val at path, with path segments separated by PathSeparator.
// Returns true if a new path was created, false if an existing path
// was overwritten.
//
// The merge function is responsible for updating the existing value
// with the new value.
func (pt *PathTrie) InsertMerge(computed string,
	url string,
	id string,
	operation *openapi3.Operation,
	method string,
	val interface{},
	merge ValueMergeFunc) (isNewPath bool) {

	trie := pt.Trie
	isNewPath = true
	// TODO: what about computed that ends with pt.PathSeparator is it different ?
	segments := strings.Split(computed, pt.PathSeparator)
	urls := strings.Split(url, pt.PathSeparator)

	// Traverse the Trie along computed, inserting nodes where necessary.
	for idx, segment := range segments {
		isLastSegment := idx == len(segments)-1
		if node, ok := trie[segment]; ok {
			if isLastSegment {
				// If this is the last computed segment, then this is the node to update.
				// If node value is not empty it means that an existing computed is overwritten
				isNewPath = IsNil(node.Value)
				merge(&node.Value, &val)
			} else {
				trie = node.Children
			}
		} else {
			var newNode *TrieNode
			if len(trie) >= 32 {
				for k := range trie {
					delete(trie, k)
				}
				// TODO merge query params
				paramName := utils.GenerateParamName()
				segments[idx] = paramName
			}
			newNode = pt.createPathTrieNode(operation,
				method,
				utils.GetName(method, strings.Join(segments, "/")),
				segments,
				urls,
				idx,
				isLastSegment,
				val)
			trie[segment] = newNode
			// continue descending.
			trie = newNode.Children

		}
	}

	return isNewPath
}

func (pt *PathTrie) createPathTrieNode(operation *openapi3.Operation,
	method string,
	id string,
	urls []string,
	segments []string,
	idx int,
	isLastSegment bool,
	val interface{}) *TrieNode {
	//fullPathSegments := segments[:idx+1]
	node := &TrieNode{
		Children:  make(PathToTrieNode),
		Name:      segments[idx],
		Path:      strings.Join(segments, pt.PathSeparator),
		Id:        id,
		Operation: operation,
		Method:    method,
		URL:       strings.Join(urls, pt.PathSeparator),
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

func (trie PathToTrieNode) getMatchNodes(segments []string, idx int) []*TrieNode {
	var nodes []*TrieNode

	isLastSegment := idx == len(segments)-1

	for _, node := range trie {
		// Check for node segment match
		if !node.isNameMatch(segments[idx]) {
			continue
		}

		// If this is the last path segment, then return node if it holds a value.
		if isLastSegment {
			if node.Value != nil {
				nodes = append(nodes, node)
			}
			continue
		}

		// Otherwise, continue descending.
		newNodes := node.Children.getMatchNodes(segments, idx+1)
		if len(newNodes) > 0 {
			nodes = append(nodes, newNodes...)
		}
	}

	return nodes
}

func (node *TrieNode) isNameMatch(segment string) bool {
	if utils.IsPathParam(node.Name) {
		return true
	}

	if node.Name == segment {
		return true
	}

	return false
}

func (node *TrieNode) isFullPathMatch(path string) bool {
	return node.Path == path
}
