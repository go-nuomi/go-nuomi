package nuomi

type NuoMiPlugin interface {
	PreRun() error
	Run(string) ([]Result, error)
	ResultToString(*Result) (*string, error)
	GetConfigString() (string, error)
}
