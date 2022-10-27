package wordgenerator

import (
	"fmt"
	"regexp"
	"strings"

	sliceTools "github.com/olehmushka/golang-toolkit/slice_tools"
	stringTools "github.com/olehmushka/golang-toolkit/string_tools"
	"github.com/olehmushka/golang-toolkit/wrapped_error"
)

type Chain map[string][]string

func AppendByBase(chains map[string]Chain, base string, el Chain) map[string]Chain {
	if len(chains) == 0 {
		chains = make(map[string]Chain)
	}

	chain, ok := chains[base]
	if !ok {
		chains[base] = el
		return chains
	}
	chains[base] = mergeChains(chain, el)
	return chains
}

func CalculateChain(baseWords []string) (Chain, error) {
	chain := make(Chain)
	for _, n := range baseWords {
		name := strings.ToLower(strings.TrimSpace(n))
		isMatch, err := regexp.MatchString("[^\u0000-\u007f]", name)
		if err != nil {
			return nil, wrapped_error.NewBadRequestError(nil, fmt.Sprintf("[CalculateChain] can not match basic chars and en rules can be applied (err = %+v)", err))
		}
		basic := !isMatch // basic chars and English rules can be applied

		// split word into pseudo-syllables
		var (
			i        = -1
			syllable = ""
		)
		for i < len(name) {
			var (
				prev = stringTools.GetCharByIndex(name, i, "") // pre-onset letter
				v    = 0                                       // 0 if no vowels in syllable
			)

			for c := i + 1; stringTools.GetCharByIndex(name, c, "") != "" && len(syllable) < 5; c++ {
				var (
					that = stringTools.GetCharByIndex(name, c, "")
					next = stringTools.GetCharByIndex(name, c+1, "") // next char
				)
				syllable += that
				if syllable == " " || syllable == "-" { // syllable starts with space or hyphen
					break
				}
				if next == "" || next == " " || next == "-" { // no need to check
					break
				}

				if stringTools.Vowel(that) { // check if letter is vowel
					v = 1
				}

				// do not split some diphthongs
				if that == "y" && next == "e" { // 'ye'
					continue
				}
				if basic {
					// English-like
					if that == "o" && next == "o" { // 'oo'
						continue
					}
					if that == "e" && next == "e" { // 'ee'
						continue
					}
					if that == "a" && next == "e" { // 'ae'
						continue
					}
					if that == "c" && next == "h" { // 'ch'
						continue
					}
				}

				// two same vowels in a row
				if stringTools.Vowel(that) && stringTools.Vowel(next) && that == next {
					break
				}

				// syllable has vowel and additional vowel is expected soon
				if afterNext := stringTools.GetCharByIndex(name, c+2, ""); v > 0 && afterNext != "" && stringTools.Vowel(afterNext) {
					break
				}
			}
			if _, ok := chain[prev]; !ok {
				chain[prev] = make([]string, 0, 1)
			}
			chain[prev] = append(chain[prev], syllable)

			// ================
			// before next iter
			i += stringTools.GetLen(syllable, 1)
			syllable = ""
		}
	}

	return chain, nil
}

func UpdateChain(baseName string, baseWords []string, chains map[string]Chain) (map[string]Chain, error) {
	if len(baseWords) == 0 {
		return nil, wrapped_error.NewBadRequestError(nil, fmt.Sprintf("base words is empty (base_name=%s)", baseName))
	}
	c, err := CalculateChain(baseWords)
	if err != nil {
		return nil, err
	}

	return AppendByBase(chains, baseName, c), nil
}

func mergeChains(c1, c2 Chain) Chain {
	out := make(Chain, (len(c1)+len(c2))/2)
	for key, value := range c1 {
		out[key] = value
	}
	for key, value := range c2 {
		if v, ok := out[key]; ok {
			out[key] = sliceTools.Unique(sliceTools.Merge(v, value))
		}
	}

	return out
}
