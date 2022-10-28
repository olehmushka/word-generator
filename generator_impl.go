package wordgenerator

import (
	"fmt"
	"strings"

	list "github.com/olehmushka/golang-toolkit/list"
	randomTools "github.com/olehmushka/golang-toolkit/random_tools"
	sliceTools "github.com/olehmushka/golang-toolkit/slice_tools"
	stringTools "github.com/olehmushka/golang-toolkit/string_tools"
	"github.com/olehmushka/golang-toolkit/wrapped_error"
)

type generator struct {
	chains *list.FIFOUniqueList[*ChainBasePair]
}

func New() Generator {
	return &generator{
		chains: newChains(),
	}
}

type GenerateOpts struct {
	BaseName  string
	BaseWords []string
	Min       int
	Max       int
	Dupl      string
}

func (g *generator) Generate(opts GenerateOpts) (string, error) {
	if opts.BaseName == "" {
		return "", wrapped_error.NewBadRequestError(nil, "Please define a base")
	}

	var err error
	var data Chain
	for i := 0; i < 20; i++ {
		d, isFound := g.chains.FindOne(func(_, curr, _ *ChainBasePair) bool {
			return curr.BaseName == opts.BaseName
		})
		if !isFound {
			if err := g.updateChains(opts.BaseName, opts.BaseWords); err != nil {
				return "", err
			}
			continue
		}
		data = d.Chain
		break
	}

	if len(data) == 0 {
		return "", wrapped_error.NewBadRequestError(nil, fmt.Sprintf("base_name %s is incorrect! [1]", opts.BaseName))
	}

	val, ok := data[""]
	if !ok {
		return "", wrapped_error.NewBadRequestError(nil, fmt.Sprintf("base_name %s is incorrect! [2]", opts.BaseName))
	}

	v := val
	var cur, w string
	cur, err = sliceTools.RandomValueOfSlice(randomTools.RandFloat64, v)
	if err != nil {
		return "", err
	}

	for i := 0; i < 20; i++ {
		if cur == "" {
			// end of word
			if len(w) < opts.Min {
				cur = ""
				w = ""
				if val, ok := data[""]; ok {
					v = val
				}
			} else {
				break
			}
		} else {
			if len(w)+len(cur) > opts.Max {
				// word too long
				if len(w) < opts.Min {
					w += cur
				}
				break
			} else {
				if val, ok := data[stringTools.GetLastStrChar(cur)]; ok {
					v = val
				}
			}
		}

		w += cur
		cur, err = sliceTools.RandomValueOfSlice(randomTools.RandFloat64, v)
		if err != nil {
			return "", err
		}
	}

	// parse word to get a final name
	l := stringTools.GetLastStrChar(w) // last letter
	if l == "'" || l == " " || l == "-" {
		w = w[0 : len(w)-1] // not allow some characters at the end
	}
	var name string
	for i, c := range w {
		var nextC, afterNextC string
		if len(w) > i+1 {
			nextC = string(w[i+1])
		}
		if len(w) > i+2 {
			afterNextC = string(w[i+2])
		}
		// duplication is not allowed
		if (nextC != "" && string(c) == nextC) && strings.Contains(opts.Dupl, string(c)) {
			continue
		}
		if len(name) == 0 {
			name += strings.ToUpper(string(c))
			continue
		}
		// remove space after hyphen
		if stringTools.GetLastStrChar(name) == "-" && string(c) == " " {
			continue
		}
		// capitalize letter after space
		if stringTools.GetLastStrChar(name) == " " {
			name += strings.ToUpper(string(c))
		}
		// capitalize letter after hyphen
		if stringTools.GetLastStrChar(name) == "-" {
			name += strings.ToUpper(string(c))
		}
		// "ae" => "e"
		if string(c) == "a" && (nextC != "" && nextC == "e") {
			continue
		}
		// remove three same letters in a row
		if i+2 < len(name) && (nextC != "" && string(c) == nextC) && (afterNextC != "" && string(c) == afterNextC) {
			continue
		}

		name += string(c)
	}

	// join the word if any part has only 1 letter
	if hasInsideStrWordsLessThan(name, 2) {
		name = makeInsideWordsCapitalized(name)
	}

	if len(name) < 2 {
		return "", wrapped_error.NewInternalServerError(nil, fmt.Sprintf("name is too short (name=%s)", name))
	}

	return name, nil
}

func (g *generator) updateChains(baseName string, baseWords []string) error {
	chains, err := UpdateChain(baseName, baseWords, g.chains)
	if err != nil {
		return wrapped_error.NewInternalServerError(err, "can not updated chains")
	}
	g.chains = chains

	return nil
}
