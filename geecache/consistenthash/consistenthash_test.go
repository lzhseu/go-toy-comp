package consistenthash

import (
	"strconv"
	"testing"
)

// 使用这份测试文件需要将编号方法由：
// key := peer + "#" + strconv.Itoa(i)
// 改为：
// key := strconv.Itoa(i) + peer

func TestHashing(t *testing.T) {
	hashmap := New(3, func(data []byte) uint32 {
		i, _ := strconv.Atoi(string(data))
		return uint32(i)
	})

	hashmap.Add("6", "4", "2")

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		if hashmap.Get(k) != v {
			t.Errorf("Asking for %s, should have yields %s", k, v)
		}
	}

	hashmap.Add("8")

	testCases["27"] = "8"

	for k, v := range testCases {
		if hashmap.Get(k) != v {
			t.Errorf("Asking for %s, should have yields %s", k, v)
		}
	}
}
