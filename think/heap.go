package think

type thinkHeap []thinker

func (h *thinkHeap) Len() int           { return len(*h) }
func (h *thinkHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }
func (h *thinkHeap) Less(i, j int) bool { return (*h)[i].next.Before((*h)[j].next) }
func (h *thinkHeap) Push(v interface{}) {
	*h = append(*h, v.(thinker))
}
func (h *thinkHeap) Pop() (v interface{}) {
	l := len(*h) - 1
	v, *h = (*h)[l], (*h)[:l]
	return
}
