// Code generated by "stringer -type=ItemType"; DO NOT EDIT

package imdb

import "fmt"

const _ItemType_name = "UnknownAnyMovieSeriesEpisode"

var _ItemType_index = [...]uint8{0, 7, 10, 15, 21, 28}

func (i ItemType) String() string {
	if i < 0 || i >= ItemType(len(_ItemType_index)-1) {
		return fmt.Sprintf("ItemType(%d)", i)
	}
	return _ItemType_name[_ItemType_index[i]:_ItemType_index[i+1]]
}
