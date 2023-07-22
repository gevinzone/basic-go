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
	res := Shrink(src[:len(src)-1])
	return res, nil
}

func calCapacity(c, l int) (int, bool) {
	if c <= 64 {
		return c, false
	}
	if c > 2048 && (c/l >= 2) {
		factor := 0.625
		return int(float32(c) * float32(factor)), true
	}
	if c <= 2048 && (c/l >= 4) {
		return c / 2, true
	}
	return c, false
}

func Shrink[T any](src []T) []T {
	c, l := cap(src), len(src)
	n, changed := calCapacity(c, l)
	if !changed {
		return src
	}
	s := make([]T, 0, n)
	s = append(s, src...)
	return s
}
