package tinygit

import "testing"

func TestArrayBounds(t *testing.T) {
	arr := make([]byte, 40)

	t.Log(len(arr[len(arr)-20:]))
}
