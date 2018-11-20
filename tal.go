package tal

import (
	"errors"
	"strconv"
	"time"
)

var (
	incompleteAnn = errors.New("Incomplete Annotation")
	unexpectedEnd = errors.New("Unexpected end of tal")
	invalidChar   = errors.New("Invalid Char")
)

const (
	onSET_PLUS  = '+'
	onSET_MINUS = '-'

	TOKEN_ONSET      = '\x15'
	TOKEN_ANNOTATION = '\x14'
	TOKEN_END        = '\x00'
)

type TimeStamp struct {
	OnSet    time.Duration
	Duration time.Duration
}

type Tal struct {
	Stamp      TimeStamp
	Annotation []byte
}

func Parse(sample []byte) ([]Tal, error) {

    var result []Tal
    i, l := 0, len(sample)

    for  i < l {

        stamp, j, err := parseStamp(sample)
        if err != nil {
            return result, err
        }

        i += j
        for ;; {

            ann, j, err := parseAnnotation(sample[i:])
            if err != nil {
                return result, err
            }

            i += j

            result = append(result, Tal{
                stamp,
                ann,
            })

            if j > len(ann) + 1 {
                break
            }


        }
    }

    return result, nil
}

func parseStamp(sample []byte) (TimeStamp, int, error) {

	i := 0
	result := TimeStamp{}

	if sample[i] != onSET_PLUS && sample[i] != onSET_MINUS {
		return result, 0, invalidChar
	}

	i, end := nextToken(sample)
	if end {
		return result, i, incompleteAnn
	}

	token := sample[i]

	if token == TOKEN_END {
		return result, i, invalidChar
	}

	onSet, err := strconv.ParseFloat(string(sample[:i]), 64)
	if err != nil {
		return result, i, err
	}

	result.OnSet = time.Duration(onSet * float64(time.Second))

	//default value in case duration is not relevant
	duration := 0.0
	if token == TOKEN_ONSET {

		i++
		s := i

		i, end = nextToken(sample[s:])
		i = s + i
		if end {
			return result, i, incompleteAnn
		}

		if sample[i] != TOKEN_ANNOTATION {
			return result, i, invalidChar
		}

		duration, err = strconv.ParseFloat(string(sample[s:i]), 64)
		if err != nil {
			return result, i, err
		}

	}

	result.Duration = time.Duration(duration * float64(time.Second))

	return result, i + 1, nil
}

//abc\x14dec\x14\x00
//abc\x14\x00\x00
func parseAnnotation(sample []byte) ([]byte, int, error) {

	l := len(sample)

	pos, end := nextToken(sample)
	if end {
		return sample, pos, incompleteAnn
	}

	token := sample[pos]
	ann := sample[:pos]

	if token == TOKEN_END || token == TOKEN_ONSET {
		return sample, pos, invalidChar
	}

	pos++

	for pos < l && sample[pos] == TOKEN_END {
		pos++
	}

	return ann, pos, nil
	//return []byte{},0,nil
}

func nextToken(input []byte) (pos int, end bool) {

	l := len(input)

	for pos < l &&
		input[pos] != TOKEN_END &&
		input[pos] != TOKEN_ANNOTATION &&
		input[pos] != TOKEN_ONSET {

		pos++
	}

	end = (pos == l)
	return
}
