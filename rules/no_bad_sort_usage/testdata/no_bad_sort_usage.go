package fixtures

import "sort"

func badIntSliceUsage() {
	var xs []int
	xs = sort.IntSlice(xs) // MATCH /suspicious sort.IntSlice usage, maybe sort.Ints was intended/
}

func badFloat64SliceUsage() {
	var xs []float64
	xs = sort.Float64Slice(xs) // MATCH /suspicious sort.Float64Slice usage, maybe sort.Float64s was intended/
}

func badStringSliceUsage() {
	var xs []string
	xs = sort.StringSlice(xs) // MATCH /suspicious sort.StringSlice usage, maybe sort.Strings was intended/
}

func goodIntsSorted() {
	xs := []int{3, 1, 2}
	sort.Ints(xs)
	_ = xs
}

func goodFloat64sSorted() {
	xs := []float64{3.0, 1.0, 2.0}
	sort.Float64s(xs)
	_ = xs
}

func goodStringsSorted() {
	xs := []string{"c", "a", "b"}
	sort.Strings(xs)
	_ = xs
}

func goodIntSliceUsedAsInterface() {
	xs := []int{3, 1, 2}
	var s sort.Interface = sort.IntSlice(xs)
	sort.Sort(s)
}

func goodFloat64SliceUsedAsInterface() {
	xs := []float64{3.0, 1.0, 2.0}
	var s sort.Interface = sort.Float64Slice(xs)
	sort.Sort(s)
}

func goodStringSliceUsedAsInterface() {
	xs := []string{"c", "a", "b"}
	var s sort.Interface = sort.StringSlice(xs)
	sort.Sort(s)
}
