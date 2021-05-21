package chainer

import (
	"encoding/binary"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	msg := NewChain()
	testcase := `d:1
g:16
gs:16
+b:00 00 00 00 00 00 00 00 00 00 00 00 00
bs:00 00 00 00 00 00 00 00 00 00
+bs:00 00 00 00 00 00 00 00 00
s:username
+s:password
ss:start1
+ss:start2`
	lines := strings.Split(testcase, "\n")
	msg.Content = append(msg.Content, lines...)
	genBytes := "+gs:16"
	msgGen := NewChain()
	msgGen.Content = append(msgGen.Content, genBytes)
	t.Run("Marshal", func(t *testing.T) {
		bts, err := msg.Marshal()
		assert.NoError(t, err)
		assert.Equal(t, 120, len(bts), "Length should be 37")
	})

	t.Run("Marshal+gs", func(t *testing.T) {
		bts, err := msgGen.Marshal()
		log.Println(bts)
		assert.NoError(t, err)
		size := binary.LittleEndian.Uint32(bts[4:8])
		assert.Equal(t, int(size), len(bts)-8, "Length should be "+strconv.Itoa(int(size)))
	})
}

func TestLoad(t *testing.T) {
	testcases := []string{`d:1
g:16
gs:16
+b:00 00 00 00 00 00 00 00 00 00 00 00 00
bs:00 00 00 00 00 00 00 00 00 00
+bs:00 00 00 00 00 00 00 00 00
s:username
+s:password
ss:start1
+ss:start2`,
		`d:1
g:16
gs:16
+b:00 00 00 00 00 00 00 00 00 00 00 00 00
bs:00 00 00 00 00 00 00 00 00 00
+bs:00 00 00 00 00 00 00 00 00
s:username
+s:password
ss:start1
+ss:start2

d:1
g:16
gs:16
+b:00 00 00 00 00 00 00 00 00 00 00 00 00
bs:00 00 00 00 00 00 00 00 00 00
+bs:00 00 00 00 00 00 00 00 00
s:username
+s:password
ss:start1
+ss:start2`,
	}
	expected := []int{1, 2}
	for i := range testcases {
		t.Run("Split chains", func(t *testing.T) {
			chains, err := LoadChains(strings.NewReader(testcases[i]))
			assert.NoError(t, err)
			assert.Equal(t, expected[i], len(chains), "Length should be 1")
		})
	}

}
