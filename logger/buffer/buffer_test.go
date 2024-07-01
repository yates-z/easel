package buffer

import "testing"

func Test(t *testing.T) {
	b := New()
	defer b.Free()
	b.WriteString("hello")
	b.WriteByte(',')
	b.Write([]byte(" world"))

	got := b.String()
	want := "hello, world"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
