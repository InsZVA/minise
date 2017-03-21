package DictionaryGenerator

import (
	"os"
	"github.com/agonopol/go-stem"
	"bufio"
)

var tire *tireTreeNode

type tireTreeNode struct {
	stop bool
	next [26]*tireTreeNode
}

func (t *tireTreeNode) run(char byte) *tireTreeNode {
	if char >= 'A' && char <= 'Z' {
		char = char - 'A' + 'a'
	}
	if t.next[char - 'a'] != nil {
		return t.next[char - 'a']
	}
	t.next[char - 'a'] = &tireTreeNode{}
	return t.next[char - 'a']
}

func (t *tireTreeNode) get(char byte) *tireTreeNode {
	if char > 'A' && char < 'Z' {
		char = char - 'A' + 'a'
	}
	if t.next[char - 'a'] != nil {
		return t.next[char - 'a']
	}
	return nil
}

func init() {
	tire = &tireTreeNode{}
	f, err := os.OpenFile("./stopwords.txt", os.O_RDONLY, 0666)
	if err != nil { panic(err) }
	defer f.Close()
	reader := bufio.NewReader(f)

	t := tire
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			// EOF
			break
		}
		stem := stemmer.Stem([]byte(s))
		for i := 0; i < len(stem); i++ {
			if stem[i] == '\r' || stem[i] == '\n' {
				break
			}
			t = t.run(stem[i])
		}
		t.stop = true
		t = tire
	}
}

// if word is a stop word return ""
// else return its stem
func Parse(word string) string {
	stem := stemmer.Stem([]byte(word))
	t := tire
	for i := 0; i < len(stem); i++ {
		if stem[i] == '\r' || stem[i] == '\n' {
			break
		}
		t = t.get(stem[i])
		if t == nil { return string(stem) }
	}
	if t.stop == true { return "" }
	return string(stem)
}