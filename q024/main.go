package main

func merge(arr1 []int, arr2 []int) []int {
	res := []int{}

	i, j := 0, 0

	for i < len(arr1) && j < len(arr2) {

		if arr1[i] < arr2[j] {
			res = append(res, arr1[i])
			i++
		} else {
			res = append(res, arr2[j])
			j++
		}

	}

}

func partition()
