package main

import "fmt"

// LazyIterator 惰性迭代器接口
// Next 返回迭代器的下一个元素和一个布尔值，布尔值表示是否还有更多元素
type LazyIterator interface {
	Next() (interface{}, bool)
}

// lazyIterator 实现惰性迭代器
type lazyIterator struct {
	nextFunc func() (interface{}, bool)
}

// Next 返回迭代器的下一个元素
func (it *lazyIterator) Next() (interface{}, bool) {
	return it.nextFunc()
}

// NewLazyIterator 创建新的惰性迭代器
// generator: 生成器函数，用于产生迭代器元素，当yield返回false时停止生成
func NewLazyIterator(generator func(yield func(interface{}) bool)) LazyIterator {
	ch := make(chan interface{})
	done := make(chan struct{})

	// 启动一个goroutine异步执行generator函数
	// defer语句确保在函数退出时关闭channel
	go func() {
		defer close(ch)
		// 调用generator函数，并传入一个yield函数
		// yield函数用于将生成的值发送到channel中
		generator(func(v interface{}) bool {
			// 使用select语句实现非阻塞的channel操作
			// 如果能成功发送到ch，则返回true继续生成
			// 如果done channel有数据，则表示需要终止生成，返回false
			select {
			case ch <- v:
				return true
			case <-done:
				return false
			}
		})
	}()

	return &lazyIterator{
		nextFunc: func() (interface{}, bool) {
			select {
			case v, ok := <-ch:
				// 如果channel已关闭，则关闭done channel
				if !ok {
					close(done)
				}
				return v, ok
			case <-done:
				return nil, false
			}
		},
	}
}

// Map 对迭代器中的每个元素应用函数
// f: 应用于每个元素的转换函数
func (it *lazyIterator) Map(f func(interface{}) interface{}) LazyIterator {
	return NewLazyIterator(func(yield func(interface{}) bool) {
		for {
			v, ok := it.Next()
			if !ok {
				return
			}
			if !yield(f(v)) {
				return
			}
		}
	})
}

// Filter 过滤迭代器中的元素
// predicate: 判断元素是否应该被保留的谓词函数
func (it *lazyIterator) Filter(predicate func(interface{}) bool) LazyIterator {
	return NewLazyIterator(func(yield func(interface{}) bool) {
		for {
			v, ok := it.Next()
			if !ok {
				return
			}
			if predicate(v) {
				if !yield(v) {
					return
				}
			}
		}
	})
}

// Take 取前n个元素
// n: 要取的元素数量
func (it *lazyIterator) Take(n int) LazyIterator {
	return NewLazyIterator(func(yield func(interface{}) bool) {
		count := 0
		// 循环直到取够n个元素或迭代器耗尽
		for count < n {
			v, ok := it.Next()
			if !ok {
				return
			}
			if !yield(v) {
				return
			}
			count++
		}
	})
}

// 辅助函数：创建一个包含指定范围数字的迭代器
func CreateNumberIterator(start, end int) LazyIterator {
	return NewLazyIterator(func(yield func(interface{}) bool) {
		for i := start; i <= end; i++ {
			if !yield(i) {
				return
			}
		}
	})
}

// 辅助函数：打印迭代器中的所有元素
func PrintAll(it LazyIterator) {
	count := 0
	for v, ok := it.Next(); ok; v, ok = it.Next() {
		fmt.Printf("%v ", v)
		count++
	}
	if count == 0 {
		fmt.Println("(空)")
	} else {
		fmt.Println()
	}
}

func main() {
	// 导入fmt包用于打印输出
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
