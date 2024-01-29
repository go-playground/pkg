//go:build go1.18
// +build go1.18

package sliceext

import (
	"sort"

	optionext "github.com/go-playground/pkg/v5/values/option"
)

// Retain retains only the elements specified by the function.
//
// This shuffles and returns the retained values of the slice.
func Retain[T any](slice []T, fn func(v T) bool) []T {
	results := make([]T, 0, len(slice))
	for _, v := range slice {
		v := v
		if fn(v) {
			results = append(results, v)
		}
	}
	return results
}

// Filter filters out the elements specified by the function.
//
// This shuffles and returns the retained values of the slice.
func Filter[T any](slice []T, fn func(v T) bool) []T {
	results := make([]T, 0, len(slice))
	for _, v := range slice {
		v := v
		if fn(v) {
			continue
		}
		results = append(results, v)
	}
	return results
}

// Map maps a slice of []T -> []U using the map function.
func Map[T, U any](slice []T, init U, fn func(accum U, v T) U) U {
	if len(slice) == 0 {
		return init
	}
	accum := init
	for _, v := range slice {
		accum = fn(accum, v)
	}
	return accum
}

// Sort sorts the sliceWrapper x given the provided less function.
//
// The sort is not guaranteed to be stable: equal elements
// may be reversed from their original order.
//
// For a stable sort, use SortStable.
func Sort[T any](slice []T, less func(i T, j T) bool) {
	sort.Slice(slice, func(j, k int) bool {
		return less(slice[j], slice[k])
	})
}

// SortStable sorts the sliceWrapper x using the provided less
// function, keeping equal elements in their original order.
func SortStable[T any](slice []T, less func(i T, j T) bool) {
	sort.SliceStable(slice, func(j, k int) bool {
		return less(slice[j], slice[k])
	})
}

// Reduce reduces the elements to a single one, by repeatedly applying a reducing function.
func Reduce[T any](slice []T, fn func(accum T, current T) T) optionext.Option[T] {
	if len(slice) == 0 {
		return optionext.None[T]()
	}
	accum := slice[0]
	for _, v := range slice {
		accum = fn(accum, v)
	}
	return optionext.Some(accum)
}

// Reverse reverses the slice contents.
func Reverse[T any](slice []T) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}
