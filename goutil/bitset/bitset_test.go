package bitset

import (
	"fmt"
	"testing"
)

var f = fmt.Printf

func Test_bit(t *testing.T) {
	var s BitSet
	s.Add(0)
	s.Add(3)
	s.Add(63)
	s.Add(100)
	s.Add(64)
	s3 := s.Copy()
	s.Add(10)
	s.Add(11)
	s.Add(128)

	s3.Add(1)
	s3.Add(2)
	s3.Add(65)
	s3.Add(101)
	fmt.Printf("s数量: %d  s3数量: %d\n", s.Len(), s3.Len())

	fmt.Printf("s: \n%ss3: \n%s", &s, s3)

	t1 := s.Copy()
	t1.UnionWith(s3)
	fmt.Printf("合并union : \n%s", t1)
	t2 := s.Copy()
	t2.DifferenceWith(s3)
	fmt.Printf("s与s3的差集: \n%s", t2)
	t3 := s.Copy()
	t3.IntersectWith(s3)
	fmt.Printf("交集intersect : \n%s", t3)
	t4 := s.Copy()
	t4.SymmetricDifference(s3)
	fmt.Printf("并差集：元素出现在A但没有出现在B，或者出现在B没有出现在A  : \n%s", t4)
	//fmt.Printf("%b  s  \n", s.words)
	//fmt.Printf("%b  s3 \n", s3.words)
	//s.SymmetricDifference(s3)
	//fmt.Printf("%b   合并后 \n", s.words)
	fmt.Printf("t1数量: %d  t2数量: %d  t3数量: %d \n", t1.Len(), t2.Len(), t3.Len())
}

func Test_2(t *testing.T) {
	i := uint64(1) << 63
	println(i)

}
