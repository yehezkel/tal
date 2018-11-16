package tal

import (
	"time"
	"errors"
	"strconv"
)

var (
	incompleteAnn = errors.New("Incomplete Annotation")
	unexpectedEnd = errors.New("Unexpected end of tal")
)

const (
	onSET = '\x15'
	dURATION = '\x14'
	end = '\x00'
	onSET_PLUS = '+' 
	onSET_MINUS = '-' 
)


type TimeStamp struct {
	OnSet time.Duration
	Duration time.Duration
}


type Tal struct {
	Stamp TimeStamp
	Annotation []byte
}

func Parse(sammple []byte) ([]Tal, error) {
	return []Tal{}, nil
}


func parseStamp(sammple []byte) (TimeStamp, int, error) {

	l,i := len(sample), 0

	if sample[i] != onSET_PLUS && sample[i] != onSET_MINUS {
		//error
	}

	for i < l && sample[i] != onSET && sample[i] != dURATION{
		
		c := sample[i]
		if c != 45 && c < 48 && c > 57 {
			//error
		}
		i++
	}

	if i == l {
		//error
	}

	onSet, err := strconv.ParseFloat(string(sample[:i]), 64)
	if err != nil {
		//error
	}

	if sample[i] == onSET {

		for i < l && sample[i] != dURATION{
		
			c := sample[i]
			if c != 45 && c < 48 && c > 57 {
				//error
			}
			i++
		}

		if i == l {
			//error
		}

				

	}





	return TimeStamp{}, 0, nil
}

//abc\x14dec\x14\x00
//abc\x14\x00\x00
func parseAnnotation(sample []byte) ([]byte, int, error) {

	i, l := 0, len(sample)
	for i < l &&  sample[i] != dURATION && sample[i] != end {
		i++
	}

	if i == l {
		return sample, i, incompleteAnn
	}

	ann := sample[:i]

	if sample[i] == dURATION {
		i++
	}

	for i < l && sample[i] == end {
		i++
	}

	return ann,i,nil 
	//return []byte{},0,nil
}