package rpm

import (
	"math"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

func NewVersion(ver string) (version Version) {
	var err error
	splitted := strings.SplitN(ver, ":", 2)
	if len(splitted) == 1 {
		version.epoch = 0
		ver = splitted[0]
	} else {
		epoch := strings.TrimLeftFunc(splitted[0], unicode.IsSpace)

		version.epoch, err = strconv.Atoi(epoch)
		if err != nil {
			version.epoch = 0
		}

		ver = splitted[1]
	}

	index := strings.Index(ver, "-")
	if index >= 0 {
		version.version = ver[:index]
		version.release = ver[index+1:]

	} else {
		version.version = ver
	}

	return version
}

func (v Version) LessThan(v1 Version) bool {
	return v.Compare(v1) < 0
}

func (v Version) Compare(v1 Version) int {
	if reflect.DeepEqual(v, v1) {
		return 0
	}
	if v.epoch > v1.epoch {
		return 1
	} else if v.epoch < v1.epoch {
		return -1
	}
	ret := rpmvercmp(v.version, v1.version)
	if ret != 0 {
		return ret
	}
	return rpmvercmp(v.release, v1.release)
}

// https://github.com/rpm-software-management/rpm/blob/master/lib/rpmvercmp.c#L16
func rpmvercmp(a, b string) int {
	if a == b {
		return 0
	}

	// get alpha/numeric segements
	segsa := alphs.FindAllString(a, -1)
	segsb := alphs.FindAllString(b, -1)
	segs := int(math.Min(float64(len(segsa)), float64(len(segsb))))

	// compare each segment
	for i := 0; i < segs; i++ {
		a := segsa[i]
		b := segsb[i]

		// compare tildes
		if []rune(a)[0] == '~' || []rune(b)[0] == '~' {
			if []rune(a)[0] != '~' {
				return 1
			}
			if []rune(b)[0] != '~' {
				return -1
			}
		}

		if unicode.IsNumber([]rune(a)[0]) {
			// numbers are always greater than alphas
			if !unicode.IsNumber([]rune(b)[0]) {
				// a is numeric, b is alpha
				return 1
			}

			// trim leading zeros
			a = strings.TrimLeft(a, "0")
			b = strings.TrimLeft(b, "0")

			// longest string wins without further comparison
			if len(a) > len(b) {
				return 1
			} else if len(b) > len(a) {
				return -1
			}

		} else if unicode.IsNumber([]rune(b)[0]) {
			// a is alpha, b is numeric
			return -1
		}

		// string compare
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
	}

	// segments were all the same but separators must have been different
	if len(segsa) == len(segsb) {
		return 0
	}

	// If there is a tilde in a segment past the min number of segments, find it.
	if len(segsa) > segs && []rune(segsa[segs])[0] == '~' {
		return -1
	} else if len(segsb) > segs && []rune(segsb[segs])[0] == '~' {
		return 1
	}

	// whoever has the most segments wins
	if len(segsa) > len(segsb) {
		return 1
	}
	return -1
}

func SecondVersionLessFirst(v1, v2 string) bool {
	vv1 := NewVersion(v1)
	vv2 := NewVersion(v2)
	return vv2.LessThan(vv1)
}
