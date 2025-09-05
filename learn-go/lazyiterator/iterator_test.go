package main

import (
	"fmt"
	"testing"
)

func Test2(t *testing.T) {
	fmt.Println("测试惰性迭代器")

	// 创建一个包含1到10数字的惰性迭代器
	it := NewLazyIterator(func(yield func(interface{}) bool) {
		for i := 1; i <= 10; i++ {
			if !yield(i) {
				return
			}
		}
	})

	// 测试基本迭代
	fmt.Println("原始数据:")
	PrintAll(it)

	// 测试Map功能：将每个数字乘以2
	it = CreateNumberIterator(1, 10)
	mappedIt := it.(*lazyIterator).Map(func(v interface{}) interface{} {
		return v.(int) * 2
	})
	fmt.Println("\n应用Map后（每个数字乘以2）:")
	PrintAll(mappedIt)

	// 测试Filter功能：只保留偶数
	it = CreateNumberIterator(1, 10)
	filteredIt := it.(*lazyIterator).Filter(func(v interface{}) bool {
		return v.(int)%2 == 0 // 只保留偶数
	})
	fmt.Println("\n应用Filter后（只保留偶数）:")
	PrintAll(filteredIt)

	// 测试Take功能：只取前5个元素
	it = CreateNumberIterator(1, 10)
	takenIt := it.(*lazyIterator).Take(5)
	fmt.Println("\n应用Take后（只取前5个元素）:")
	PrintAll(takenIt)

	// 测试链式操作：先过滤出偶数，然后乘以3，最后取前2个
	it = CreateNumberIterator(1, 10)
	chainedIt := it.(*lazyIterator).
		Filter(func(v interface{}) bool {
			return v.(int)%2 == 0
		}).(*lazyIterator).Map(func(v interface{}) interface{} {
		return v.(int) * 3
	}).(*lazyIterator).Take(2)

	fmt.Println("\n链式操作（过滤偶数，乘以3，取前2个）:")
	PrintAll(chainedIt)
}
