package tal

import (
	"testing"
	"time"
)

func TestParseAnnotationSimple(t *testing.T) {
	input := "abcd\x14\x00"
	ann, count, err := parseAnnotation([]byte(input))
	if err != nil {
		t.Fatal(err)
	}

	if count != 6 {
		t.Errorf("Wrong count %d, expecting %d", count, 6)
	}

	if string(ann) != "abcd" {
		t.Errorf("Unexpected annotation: %s", string(ann))
	}
}

func TestParseAnnotation(t *testing.T) {

	type Iteration struct {
		annotation string
		used       int
	}

	table := []struct {
		input      string
		iterations []Iteration
	}{
		{
			"abc\x14cdb\x14\x00",
			[]Iteration{

				Iteration{"abc", 4},
				Iteration{"cdb", 5},
			},
		},
		{
			"abc\x14\x14\x00",
			[]Iteration{

				Iteration{"abc", 4},
				Iteration{"", 2},
			},
		},
	}

	for _, current := range table {
		input := []byte(current.input)

		for _, iteration := range current.iterations {
			ann, count, err := parseAnnotation(input)
			if err != nil {
				t.Fatal(err)
			}

			if string(ann) != iteration.annotation {
				t.Errorf("Unexpected parsed annotation: %s expecting %s", string(ann), iteration.annotation)
			}

			if count != iteration.used {
				t.Errorf("Wrong byte processed: %d, expecting %d", count, iteration.used)
			}

			input = input[count:]
		}
	}
}

func TestTimeStampNoSign(t *testing.T) {

	input := "120\x151\x14test\x14\x00"
	_, _, err := parseStamp([]byte(input))
	if err == nil {
		t.Errorf("No error given but expecting: %s", invalidChar)
		return
	}

	if err != invalidChar {
		t.Errorf("Non expected error: %s expecting %s", err, invalidChar)
	}
}

func TestTimeIncompleted(t *testing.T) {

	table := []string{
		"+",
		"+8\x15",
		"+8\x151",
	}

	for _, input := range table {
		_, _, err := parseStamp([]byte(input))
		if err == nil {
			t.Errorf("No error given but expecting: %s", incompleteAnn)
			continue
		}

		if err != incompleteAnn {
			t.Errorf("Non expected error: %s expecting %s", err, incompleteAnn)
		}
	}
}

func TestTimeInvalid(t *testing.T) {

	table := []string{
		".12\x14123\x14\x00",
		"-1\x00",
	}

	for _, input := range table {
		_, _, err := parseStamp([]byte(input))
		if err == nil {
			t.Errorf("No error given but expecting: %s", invalidChar)
			continue
		}

		if err != invalidChar {
			t.Errorf("Non expected error: %s expecting %s", err, invalidChar)
		}
	}
}

func TestTimeStampBadNumber(t *testing.T) {

	table := []string{
		"-ab\x14123\x14\x00",
		"-1\x15ab\x14\x00",
	}

	for _, input := range table {
		_, _, err := parseStamp([]byte(input))
		if err == nil {
			t.Errorf("No error given but expecting strconv.ParseFloat erro")
		}
	}
}

func TestParse(t *testing.T) {

	type FlatTal struct {
		onset      string
		duration   string
		annotation string
	}

	type Iteration struct {
		input   string
		results []FlatTal
	}

	table := []Iteration{
		Iteration{
			"+120\x151\x14test\x14\x00",
			[]FlatTal{
				FlatTal{"120s", "1s", "test"},
			},
		},

		Iteration{
			"+120\x14test\x14\x00",
			[]FlatTal{
				FlatTal{"120s", "0s", "test"},
			},
		},

		Iteration{
			"+120\x14test\x14test2\x14\x00",
			[]FlatTal{
				FlatTal{"120s", "0s", "test"},
				FlatTal{"120s", "0s", "test2"},
			},
		},

		Iteration{
			"+120.3\x150.5\x14test\x14test2\x14\x00\x00",
			[]FlatTal{
				FlatTal{"120.3s", "0.5s", "test"},
				FlatTal{"120.3s", "0.5s", "test2"},
			},
		},
		Iteration{
			"+120\x14test\x14\x00+120.3\x150.5\x14test\x14test2\x14\x00",
			[]FlatTal{
				FlatTal{"120s", "0s", "test"},
				FlatTal{"120.3s", "0.5s", "test"},
				FlatTal{"120.3s", "0.5s", "test2"},
			},
		},
	}

	for _, current := range table {

		list, err := Parse([]byte(current.input))
		if err != nil {
			t.Fatal(err)
		}

		if len(list) != len(current.results) {
			t.Errorf("Unexpected results count: %d expecting %d", len(list), len(current.results))
		}

		for i, ann := range list {

			rawExpected := current.results[i]
			expOnset, _ := time.ParseDuration(rawExpected.onset)
			expDuration, _ := time.ParseDuration(rawExpected.duration)

			if ann.Stamp.OnSet != expOnset {
				t.Errorf("Unexpected onset on %d got: %s expecting %s", i, ann.Stamp.OnSet, expOnset)
			}

			if ann.Stamp.Duration != expDuration {
				t.Errorf("Unexpected duration on %d got: %s expecting %s", i, ann.Stamp.Duration, expDuration)
			}

			if rawExpected.annotation != string(ann.Annotation) {
				t.Errorf("Unexpected annotation on %d got: %s expecting %s", i, ann.Annotation, rawExpected.annotation)
			}
		}
	}
}
