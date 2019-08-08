package bitset

import "fmt"

/**
位数组, 标志位数组极简用法
索引从0开始
@author kidd

*/

type BitSet struct {
	words []uint64
}

const bitLength = 64

// 返回第N个组, 第M位
func getWordAndBit(x int) (int, uint) {
	word, bit := x/bitLength, uint(x%bitLength)
	return word, bit
}

// 设置位
func (s *BitSet) Add(x int) {
	word, bit := getWordAndBit(x)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= uint64(1) << bit
}

// 检查位
func (s *BitSet) Has(x int) bool {
	word, bit := getWordAndBit(x)
	return len(s.words) >= word && s.words[word]&(1<<bit) != 0
}

// 合并另一个位数组
func (s *BitSet) UnionWith(t *BitSet) {
	for x, tword := range t.words {
		if x < len(s.words) {
			s.words[x] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// 与另一个bit数组的交集
func (s *BitSet) IntersectWith(t *BitSet) {
	if len(s.words) <= len(t.words) {
		for x := range s.words {
			s.words[x] &= t.words[x]
		}
	} else {
		for x := range s.words {
			if x < len(t.words) {
				s.words[x] &= t.words[x]
			} else {
				s.words[x] = 0
			}
		}
	}
}

// 差集, 元素出现在s集合，未出现在t集合
func (s *BitSet) DifferenceWith(t *BitSet) {
	for x := range s.words {
		if x < len(t.words) {
			a := s.words[x] ^ t.words[x] // 先求异或,再并集
			s.words[x] &= a
		}
	}
}

// 并差集：元素出现在A但没有出现在B，或者出现在B没有出现在A
func (s *BitSet) SymmetricDifference(t *BitSet) {
	p := t
	if len(s.words) > len(t.words) {
		p = s
	}
	for x := 0; x < len(p.words); x++ {
		if x < len(s.words) && x < len(t.words) {
			s.words[x] = s.words[x] ^ t.words[x]
		} else if x >= len(s.words) && x < len(t.words) {
			s.words[x] = t.words[x]
		}
	}
}

// return the number of elements
func (s *BitSet) Len() int {
	return len(s.words)
}

// remove x from the set
func (s *BitSet) Remove(x int) {
	word, bit := getWordAndBit(x)
	if word > len(s.words) {
		return
	}
	u := (uint64(1)<<(bitLength-1) + (uint64(1)<<(bitLength-1) - 1)) ^ uint64(1)<<bit
	s.words[word] &= u
}

// remove all elements from the set
func (s *BitSet) Clear() {
	var w = []uint64{}
	s.words = w
}

// return a copy of the set
func (s *BitSet) Copy() *BitSet {
	var w = BitSet{} // deepcopy
	for len(w.words) < len(s.words) {
		w.words = append(w.words, 0)
	}
	for x, val := range s.words {
		w.words[x] = val
	}
	return &w
}

func (s *BitSet) String() string {
	out := ""
	for x := range s.words {
		out += fmt.Sprintf("%03d   %64b\n", x, s.words[x])
	}
	return out
}
