package slices

import "slices"

// Contains checks if slice contains given element.
func Contains[E comparable](slice []E, element E) bool {
	return slices.Contains(slice, element)
}
