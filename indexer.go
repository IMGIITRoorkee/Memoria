package memoria

type Indexer interface {
	Initialize(keys <-chan string)
	Insert(key string)
	Delete(key string)
	Keys(frm string, n int) []string
}
