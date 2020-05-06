package Tools

import (
	"sort"
)

//对[]string 类型的切片进行元素唯一化处理，返回来的元素已经排好序
func SetString(arr []string) []string {
	newArr := make([]string, 0)

	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
			}

		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	sort.Strings(newArr)
	return (newArr)
}

//对[]int 类型的切片进行元素唯一化处理，返回来的元素已经排好序
func SetInt(arr []int) []int {
	newArr := make([]int, 0)

	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
			}

		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	sort.Ints(newArr)
	return (newArr)
}

//c := []interface{}{1,2,3,3,4,4}
//	d := Tools.Set(c)	//[1 2 3 4]
//注意切片时的语法 c := interface{}{1,1,2,2,3,3,4,5}

func Set(arr []interface{}) []interface{} {
	newArr := make([]interface{}, 0)

	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
			}

		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}

	return (newArr)
}

func GetNums(a interface{}, arr []interface{}) int {
	nums := 0

	for _, b := range arr {
		if a == b {
			nums++
		}
	}
	return nums
}

//获取一个切片中出现最多次的元素
func Most(arr []interface{}) interface{} {
	//对arr进行元素唯一化处理，返回arr2
	arr2 := Set(arr)
	//将arr2传入 方法GetNums(a interface{},arr []interface{}) int,返回一个数字，代表这个元素
	most := 0
	var maxnum interface{}
	for _, item := range arr2 {
		itemnums := GetNums(item, arr)
		if itemnums > most {
			most = itemnums
			maxnum = item
		}
	}
	//出现过几次,num接收这个值，然后对比这个值，将数值大的那个元素的下标存入maxnum中
	//返回这个次数最大的元素

	return maxnum
}

///
///slices... []string	传入多个切片
///
func SliceJoin(slices ...[]string) []string {
	//先获取这些数组的len
	count := 0
	for _, item := range slices {
		count += len(item)

	}
	//声明并初始化一个最终用于接受所有合并切片的切片
	result := make([]string, 0, count)
	for _, slice := range slices {
		if slice != nil {
			for _, v := range slice {
				result = append(result, v)
			}
		}

	}

	return result
}
