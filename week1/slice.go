package week1

import "errors"

var ErrIndexOutOfRange = errors.New("下标超出范围")

// DeleteAt 删除指定位置的元素
// 如果下标不是合法的下标，返回 ErrIndexOutOfRange
func DeleteAt[T any](src []T, index int) ([]T, error) {
	if src == nil || index < 0 || index >= len(src) {
		return nil, ErrIndexOutOfRange
	}
	for i := index; i < len(src)-1; i++ {
		src[i] = src[i+1]
	}
	return src[:len(src)-1], nil
}
