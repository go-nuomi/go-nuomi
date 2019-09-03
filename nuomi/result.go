package nuomi

type ResultStatus int

const (
	StatusFound  ResultStatus = iota
	StatusMissed ResultStatus = iota
)

type Result struct {
	Entity     string
	StatusCode int
	Status     ResultStatus
	Extra      string
	Size       *int64
}

//转成string，这里后来需要入库
func (r *Result) ToString(n *NuoMi) (string, error) {
	s, err := n.plugin.ResultToString(r)
	if err != nil {
		return "", nil
	}
	return *s, nil
}
