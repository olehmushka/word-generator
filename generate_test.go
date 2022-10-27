package wordgenerator

import (
	"testing"
	"time"

	slicetools "github.com/olehmushka/golang-toolkit/slice_tools"
)

func TestGenerate(t *testing.T) {
	opts := GenerateOpts{
		BaseName: "celtic_wb",
		BaseWords: []string{
			"Aberaman",
			"Aberangell",
			"Aberarth",
			"Aberavon",
			"Aberbanc",
			"Aberbargoed",
			"Aberbeeg",
			"Abercanaid",
			"Abercarn",
			"Abercastle",
			"Abercegir",
			"Abercraf",
			"Abercregan",
			"Abercych",
			"Abercynon",
			"Aberdare",
			"Aberdaron",
			"Aberdaugleddau",
			"Aberdeen",
			"Aberdulais",
			"Aberdyfi",
			"Aberedw",
			"Abereiddy",
			"Abererch",
			"Abereron",
			"Aberfan",
			"Aberffraw",
			"Aberffrwd",
			"Abergavenny",
			"Abergele",
			"Aberglasslyn",
			"Abergorlech",
			"Abergwaun",
			"Abergwesyn",
			"Abergwili",
			"Abergwynfi",
			"Abergwyngregyn",
			"Abergynolwyn",
			"Aberhafesp",
			"Aberhonddu",
			"Aberkenfig",
			"Aberllefenni",
			"Abermain",
			"Abermaw",
			"Abermorddu",
			"Abermule",
			"Abernant",
			"Aberpennar",
			"Aberporth",
			"Aberriw",
			"Abersoch",
			"Abersychan",
			"Abertawe",
			"Aberteifi",
			"Aberthin",
			"Abertillery",
			"Abertridwr",
			"Aberystwyth",
			"Achininver",
			"Afonhafren",
			"Alisaha",
			"Antinbhearmor",
			"Ardenna",
			"Attacon",
			"Beira",
			"Bhrura",
			"Boioduro",
			"Bona",
			"Boudobriga",
			"Bravon",
			"Brigant",
			"Briganta",
			"Briva",
			"Cambodunum",
			"Cambra",
			"Caracta",
			"Catumagos",
			"Centobriga",
			"Ceredigion",
			"Chalain",
			"Dinn",
			"Diwa",
			"Dubingen",
			"Duro",
			"Ebora",
			"Ebruac",
			"Eburodunum",
			"Eccles",
			"Eighe",
			"Eireann",
			"Ferkunos",
			"Genua",
			"Ghrainnse",
			"Inbhear",
			"Inbhir",
			"Inbhirair",
			"Innerleithen",
			"Innerleven",
			"Innerwick",
			"Inver",
			"Inveraldie",
			"Inverallan",
			"Inveralmond",
			"Inveramsay",
			"Inveran",
			"Inveraray",
			"Inverarnan",
			"Inverbervie",
			"Inverclyde",
			"Inverell",
			"Inveresk",
			"Inverfarigaig",
			"Invergarry",
			"Invergordon",
			"Invergowrie",
			"Inverhaddon",
			"Inverkeilor",
			"Inverkeithing",
			"Inverkeithney",
			"Inverkip",
			"Inverleigh",
			"Inverleith",
			"Inverloch",
			"Inverlochlarig",
			"Inverlochy",
			"Invermay",
			"Invermoriston",
			"Inverness",
			"Inveroran",
			"Invershin",
			"Inversnaid",
			"Invertrossachs",
			"Inverugie",
			"Inveruglas",
			"Inverurie",
			"Kilninver",
			"Kirkcaldy",
			"Kirkintilloch",
			"Krake",
			"Latense",
			"Leming",
			"Lindomagos",
			"Llanaber",
			"Lochinver",
			"Lugduno",
			"Magoduro",
			"Monmouthshire",
			"Narann",
			"Novioduno",
			"Nowijonago",
			"Octoduron",
			"Penning",
			"Pheofharain",
			"Ricomago",
			"Rossinver",
			"Salodurum",
			"Seguia",
			"Sentica",
			"Theorsa",
			"Uige",
			"Vitodurum",
			"Windobona",
		},
		Min:  4,
		Max:  12,
		Dupl: "ndl",
	}

	g := New()
	firstRoundStartTime := time.Now()
	result, err := g.Generate(opts)
	if err != nil {
		t.Errorf("can not generate word (err=%+v)", err)
		return
	}
	firstRoundEndTime := time.Since(firstRoundStartTime)
	if len(result) < opts.Min {
		t.Errorf("incorrect generated word (word=%s)", result)
	}

	secondRoundStartTime := time.Now()
	if _, err = g.Generate(opts); err != nil {
		t.Errorf("can not generate word (err=%+v)", err)
		return
	}
	secondRoundEndTime := time.Since(secondRoundStartTime)
	if firstRoundEndTime.Nanoseconds() < secondRoundEndTime.Microseconds() {
		t.Errorf("second generation words for the same word base can not have greater duration than first")
	}

	iterCount := 1000
	generated := make([]string, 0, iterCount)
	for i := 0; i < iterCount; i++ {
		curr, err := g.Generate(opts)
		if err != nil {
			t.Errorf("unexpected error = %+v", err)
			return
		}
		generated = append(generated, curr)
	}
	uniqueGenerated := slicetools.Unique(generated)
	expectedMinOriginalWordsPercentage := 0.75
	originalPercentage := float64(len(uniqueGenerated)) / float64(iterCount)
	if originalPercentage < expectedMinOriginalWordsPercentage {
		t.Errorf("original words too low percentage (expected=%.2f, actual=%.2f)", expectedMinOriginalWordsPercentage, originalPercentage)
	}
}
