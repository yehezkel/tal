package tal

import (
	"testing"
)

func TestParseAnnotationSimple(t *testing.T) {
	input := "abcd\x14\x00"
	ann, count, err := parseAnnotation([]byte(input))
	if err != nil {
		t.Fatal(err)
	}

	if count != 6 {
		t.Errorf("Wrong count %d, expecting %d",count,6)
	}

	if string(ann) != "abcd" {
		t.Errorf("Unexpected annotation: %s",string(ann))
	} 
}


func TestParseAnnotation(t *testing.T) {

	type Iteration struct {
		annotation string
		used int
	}

	table := []struct{
		input string
		iterations []Iteration
	}{
		{
			"abc\x14cdb\x14\x00", 
			[]Iteration{
				
				Iteration{"abc",4},
				Iteration{"cdb",5},
				
			},
		},
		{
			"abc\x14\x14\x00", 
			[]Iteration{
				
				Iteration{"abc",4},
				Iteration{"",2},
				
			},
		},
	}

	for _, current := range table {
		input := []byte(current.input)
		
		for _,iteration := range current.iterations {
			ann, count,err :=  parseAnnotation(input)
			if err != nil {
				t.Fatal(err)
			}

			if string(ann) != iteration.annotation {
				t.Errorf("Unexpected parsed annotation: %s expecting %s",string(ann), iteration.annotation)
			}

			if count != iteration.used {
				t.Errorf("Wrong byte processed: %d, expecting %d",count, iteration.used)
			}
			
			input = input[count:]
		}
	}
}