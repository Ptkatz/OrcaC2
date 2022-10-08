package portcrackopt

type Options struct {
	// input
	Input     string
	InputFile string
	Module    string
	User      string
	Pass      string
	UserFile  string
	PassFile  string
	// config
	Threads  int
	Timeout  int
	Delay    int
	CrackAll bool
	// output
	OutputFile string
	NoColor    bool
	// debug
	Silent bool
	Debug  bool

	Targets  []string
	UserDict []string
	PassDict []string
}
