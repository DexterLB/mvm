// Code generated by "stringer -type=ShowType"; DO NOT EDIT

package imdb

import "fmt"

const _ShowType_name = "UnknownMovieSeriesEpisode"

var _ShowType_index = [...]uint8{0, 7, 12, 18, 25}

func (i ShowType) String() string {
	if i < 0 || i >= ShowType(len(_ShowType_index)-1) {
		return fmt.Sprintf("ShowType(%d)", i)
	}
	return _ShowType_name[_ShowType_index[i]:_ShowType_index[i+1]]
}
