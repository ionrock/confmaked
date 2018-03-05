package maketypes

type Makefile interface {
	GetName()
	Download() error
}

// type URLMakefile struct {
// 	Name string
// 	URL  string
// }

// func (c URLMakefile) Download() error {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return err
// 	}

// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)

// }
