package fs

type CacheNode struct {
	name     string
	isDir    bool
	children []*CacheNode
}

var root *CacheNode

func InitCache() error {
	root = &CacheNode{
		name:  "/",
		isDir: true,
	}
	return nil
}
