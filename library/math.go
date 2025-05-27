package library

func Factorial(num int) int {
	if num == 1 || num == 0 {
		return 1
	}
	return num * Factorial(num-1)
}
