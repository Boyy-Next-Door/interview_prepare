package leetcode_test

import (
	"fmt"
	"testing"
)

func TestSolve(t *testing.T) {
	// ret := romanToInt("MCMXCIV")
	// ret := longestCommonPrefix([]string{"fl1232123123ower","fl123123123ow","fl123123123ight"})
	// ret  := threeSum([]int {-1,0,1,2,-1,-4,-2,-3,3,0,4})
	// threeSumClosest([]int {0,0,0}, 1)
	// letterCombinations("23")
	// fourSum([]int{1,0,-1,0,-2,2}, 0)
	// removeNthFromEnd([]int{1,0,-1,0,-2,2}, 0)
	// head := &ListNode{
	// 	Val: 1,
	// 	Next: &ListNode{
	// 		Val: 5,
	// 		Next: &ListNode{
	// 			Val: 8, 
	// 			Next: &ListNode{
	// 				Val: 9,
	// 				Next: nil,
	// 			},
	// 		},
	// 	},
	// }

	// head2 := &ListNode{
	// 	Val: 3,
	// 	Next: &ListNode{
	// 		Val: 4,
	// 		Next: &ListNode{
	// 			Val: 7, 
	// 			Next: &ListNode{
	// 				Val: 10,
	// 				Next: nil,
	// 			},
	// 		},
	// 	},
	// }

	// removeNthFromEnd(head, 1)
	// isValid("[][]({[]})")
	// mergeTwoLists(head, head2)
	// generateParenthesis(4)

	// fmt.Println(ret)

	// mergeKLists([]*ListNode{
	// 	head, head2,
	// })

	// fmt.Println(divide(4, 1))

	// fmt.Println(findSubstring("barfoofoobarthefoobarman", []string{"foo","bar","the"}))

} 

func solve(dividend, divisor int) int {
    if dividend < divisor {
		return 0
	} 
    times := 1 
    a, b := dividend, divisor

    for a > b {
        times <<= 1
        b <<= 1
    }

    if a < b {
		if times == 1 {
			return 0
		} else {
			return times >> 1 + solve(a - b >> 1, divisor)
		}
    } 
	return times
}


func choose(lists []*ListNode, curr *ListNode) (bool, *ListNode) {
	var minValNode *ListNode = nil
	minValNodeIdx := 0

	for i := 0; i < len(lists); i ++ {
		if lists[i] != nil && (minValNode == nil || lists[i].Val < minValNode.Val) {
			minValNode = lists[i]
			minValNodeIdx = i
		}
	}

	curr.Next = minValNode
	curr = curr.Next

	success := minValNode != nil
	if success {
		lists[minValNodeIdx] = minValNode.Next 
	}
	return success, curr
}


/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
 func swapPairs(head *ListNode) *ListNode {
	if head == nil || head.Next == nil{
		return head
	}

	newHead := head.Next
	curr := head
	pre := &ListNode{
	}
	for curr != nil{
		next := curr.Next

		if next != nil {
			curr.Next = next.Next
			pre.Next = next
			next.Next = curr

			pre = curr
			curr = curr.Next
		} else {
			// 有一个节点单出来了  啥也不用干
			return newHead
		}
	}

	return newHead
 }

func generateParenthesis(n int) []string {
	stack := []string{"("}
	ret := []string{}
	handle(stack, n - 1, "(", &ret)
	return ret
}

func handle(stack []string, n int, curr string, ret *[]string) {
	if n == 0 {
		for i := 0; i < len(stack); i ++ {
			curr += ")"
		}
		(*ret) = append((*ret), curr)
		return
	} else {
		// 当前还有左括号可以加入stack
		// pop一个(
		if len(stack) > 0 {
			handle(stack[0:len(stack)-1], n, curr + ")", ret)
		}
		// push一个(
		stack = append(stack, "(")
		handle(stack, n - 1, curr + "(", ret)
	}
}
/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
 func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	if list1 == nil {
		return list2
	}
	if list2 == nil {
		return list1 
	}
	ret := &ListNode{
		Next: nil,
	}
	head := ret

	for list1 != nil && list2 != nil {
		if list1.Val < list2.Val {
			ret.Next = list1
			list1 = list1.Next
		} else {
			ret.Next = list2
			list2 = list2.Next
		}
		ret = ret.Next
	}

	if list1 != nil {
		ret.Next = list1
	} else if list2 != nil {
		ret.Next = list2
	}


	return head.Next
 }

func isValid(s string) bool {
	if len(s) % 2 == 1 || len(s) == 0 {
		return false 
	}
	stack := []byte{}
	braceMap := map[byte]byte {
		'(': ')',
		'{': '}',
		'[': ']',
	}
	for i:=0; i < len(s); i++ {
		if len(stack) == 0 || s[i] != stack[len(stack) - 1] {
			right, exist := braceMap[s[i]]
			if !exist {
				return false
			}
			stack = append(stack, right)
		} else {
			stack = stack[0:len(stack) - 1]
		}
	}

	return len(stack) == 0
}

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */

type ListNode struct {
	Val int
	Next *ListNode
}

func removeNthFromEnd(head *ListNode, n int) *ListNode {
	curr := head
	target := head
	var pre *ListNode
	diff := 0

	for {
		if diff == n - 1 {
			if curr.Next == nil {
				// pre 连到 curr
				if pre == nil {
					return target.Next
				} else {
					pre.Next = target.Next
					return head
				}
			} else {
				pre = target
				target = target.Next 
			}
		} else {
			// 还没达到目标差距 target不动
			diff += 1
		}
		curr = curr.Next
	}

	// 当curr为nil时  此时的
	return nil
}


func fourSum(nums []int, target int) [][]int {
	sort(&nums)
	var res [][]int

	for i := 0; i < len(nums) - 3; i ++ {
		if i > 0 && nums[i - 1] == nums[i] {
			continue
		}

		for j := i + 1; j < len(nums) - 2; j ++ {
			if j > i + 1 && nums[j - 1] == nums[j] {
				continue
			}

			l, r := j + 1, len(nums) - 1

			if nums[i] + nums[j] > target && nums[l] >= 0 {
				continue
			}

			newTarget := target - nums[i] - nums[j]

			for l < r {
				if nums[l] + nums[r] == newTarget {
					res = append(res, []int{nums[i], nums[j], nums[l], nums[r]})

					// l右  r左
					l += 1
					for l < r {
						if nums[l - 1] == nums[l]  {
							l += 1
							continue
						}
						break
					}
					r -= 1
					for l < r {
						if nums[r + 1] == nums[r]  {
							r -= 1
							continue
						}
						break
					}
				}

				if nums[l] + nums[r] < newTarget {
					// 小了  l往右走
					l += 1
					for l < r {
						if nums[l - 1] == nums[l]  {
							l += 1
							continue
						}
						break
					}
				} else if nums[l] + nums[r] > newTarget {
					// 大了  r往左走
					r -= 1
					for l < r {
						if nums[r + 1] == nums[r]  {
							r -= 1
							continue
						}
						break
					}
				}
			}
		}
	}

	fmt.Println(res)
	return res
}

func letterCombinations(digits string) []string {
	if digits == "" {
		return []string{}
	}
	
	// res[i]  表示 从第i个按键开始的所有组合  是一个数组
	// res[0] 就是结果 
	// res[len - 1] = 3或4个字符 （对应按键上的四个字符）
	// res[i] = digits[i]上的3或4个字符分别追加上res[i + 1]中的所有元素

	keyMap := map[string][]string {
		"2": {"a", "b", "c"}, 
		"3": {"d", "e", "f"}, 
		"4": {"g", "h", "i"}, 
		"5": {"j", "k", "l"}, 
		"6": {"m", "n", "o"}, 
		"7": {"p", "q", "r", "s"}, 
		"8": {"t", "u", "v"}, 
		"9": {"w", "x", "y", "z"}, 
	}
	res := make([][]string, len(digits))

	res[len(digits) - 1] = keyMap[string(digits[len(digits) - 1])]
	
	for i := len(digits) - 2; i > -1 ;i -- {
		for j := range res[i + 1] {
			currKeyList := keyMap[string(digits[i])]
			for k:=0; k < len(currKeyList); k ++ {
				res[i] = append(res[i], currKeyList[k] + res[i + 1][j])
			}
		}
	}


	// fmt.Println(res)
	return res[0]
}

func threeSumClosest(nums []int, target int) int {
	minDiff := 999999
	minDiffVal := 0
	sort(&nums)
	for i:=0; i < len(nums) - 2; i ++ {
		l, r := i + 1, len(nums) - 1



		for l < r {
			currDiff := abs(nums[i] + nums[l] + nums[r] - target)
			// 找到一组可行解 
			if abs(currDiff) < minDiff {
				minDiff = currDiff
				minDiffVal = nums[i] + nums[l] + nums[r]

			}
			if nums[i] + nums[l] + nums[r] > target {
				// 大了  r往左走
				r -= 1
				for r > l {
					if nums[r] == nums[r + 1] {
						r -= 1
					} else {
						break
					}
				}
			} else if nums[i] + nums[l] + nums[r] < target {
				// 小了  l往右走
				l += 1
				for l < len(nums) {
					if nums[l] == nums[l - 1] {
						l += 1
					} else {
						break
					}
				}
			} else {
				// 已经完美相等  没必要再找了
				return target
			}
		}

		for i + 1 < len(nums) && nums[i + 1] == nums [i] {
			i += 1
		}
	}

	fmt.Println(minDiff)
	fmt.Println(minDiffVal)


	return minDiffVal
}

func abs(a int) int {
	if a > 0 {
		return a
	}
	return -1 * a
}

func threeSum(nums []int) [][]int {
	sort(&nums)
	result := make([][]int, 0)
	fmt.Println(nums)
	l, r := 1, len(nums) - 1
	for i:=0; i < len(nums) && nums[i] <= 0; i ++ {
		l = i + 1
		r = len(nums) - 1
		for l < r {
			// 找到一组可行解
			if nums[l] + nums[r] + nums[i] == 0 {
				result = append(result, []int{nums[i], nums[l], nums[r]})
				// l往右走 这里要跳过与可行解中l相同的元素
				for l + 1 < len(nums) {
					if nums[l + 1] == nums[l] {
						l += 1
					} else {
						l += 1
						break
					}
				}
				// l变大了  r一定要变小  往左走
				for r > l {
					if nums[r-1] == nums[r] {
						r -= 1
					} else {
						r -= 1
						break
					}
					
				}
			} else if nums[l] + nums[r] + nums[i] > 0 {
				// 已经大于0了 只能r往左走
				for r > l {
					if nums[r-1] == nums[r] {
						r -= 1
					} else {
						r -= 1
						break
					}
				}
			} else if nums[l] + nums[r] + nums[i] < 0 {
				// 小于0 l可以往右走 r不动
				for l + 1< len(nums) {
					if nums[l + 1] == nums[l] {
						l += 1
					} else {
						l += 1
						break
					}
				}
			}
		}
		for i + 1 < len(nums) && nums[i + 1] == nums[i] {
			i += 1
		}
	}


	return result
} 

func sort(list *[]int) {
	for i:=0; i < len(*list) - 1; i ++ {
		for j := i + 1; j < len(*list); j ++ {
			if (*list)[i] > (*list)[j] {
				tmp := (*list)[i]
				(*list)[i] = (*list)[j]
				(*list)[j] = tmp
			}
		}
	}
	fmt.Println(list)
}

func longestCommonPrefix(strs []string) string {
	minLen := minLength(strs)
	ret := ""
	for i:=0; i < minLen; i ++ {
		currChar := strs[0][i]

		for j := 1; j < len(strs); j ++ {
			if strs[j][i] != currChar {
				return ret
			}
		}
		ret += string(currChar)
	}

	return ret
}

func minLength(strs []string) int {
	min := 200
	for i := range strs {
		if len(strs[i]) < min {
			min =  len(strs[i])
		}
	}
	return min
}


func romanToInt(s string) int {
	valMap := map[byte]int{
		'I': 1,
		'V': 5,
		'X': 10,
		'L': 50,
		'C': 100,
		'D': 500,
		'M': 1000,
	}

	sum := 0
	for i :=len(s) - 1; i > -1; i -- {
		// currChar := string(s[i])
		if i < len(s) - 1 && valMap[s[i]] < valMap[s[i+1]] {
			sum -= valMap[s[i]]
		} else {
			sum += valMap[s[i]]
		}
	}

	return sum
}