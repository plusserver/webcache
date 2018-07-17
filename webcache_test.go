// dheilema 2018
// webcache tests

package webcache

import (
	"reflect"
	"testing"
	"time"
)

var (
	myPage CachedPage
)

// a basic request flow
func Test_Simple(t *testing.T) {
	content := []byte("Test_Simple")
	p := NewCachedPage(time.Second)
	err := p.StartUpdate()
	if err != nil {
		t.Fatal(err)
	}
	n, err := p.Write(content)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(content) {
		t.Fatalf("Error in .Write(): Got %d expected %s\n", n, len(content))
	}
	p.EndUpdate()
	got := p.Get()
	if !reflect.DeepEqual(got, content) {
		t.Fatalf("Error in .Get(): Got %s expected %s\n", got, content)
	}
}

// the second start should raise an error
func Test_DoubleStart(t *testing.T) {
	p := NewCachedPage(time.Second)
	p.StartUpdate()
	err := p.StartUpdate()
	if err != ErrUpdateInProgress {
		t.Fatalf("Error in .StartUpdate(): Got %s expected %s\n", err, ErrUpdateInProgress)
	}
}

// write without start should raise an error
func Test_ForgottenStart(t *testing.T) {
	content := []byte("Test_ForgottenStart")
	p := NewCachedPage(time.Second)
	_, err := p.Write(content)
	if err != ErrWriteWithoutUpdate {
		t.Fatalf("Error in .Write(): Got %s expected %s\n", err, ErrWriteWithoutUpdate)
	}
}

func Test_Valid(t *testing.T) {
	content := []byte("Test_Valid")
	p := NewCachedPage(time.Second)
	p.StartUpdate()
	p.Write(content)
	p.EndUpdate()
	if !p.Valid() {
		t.Fatalf("Error in .Valid(): Got false expected true\n")
	}
	time.Sleep(time.Second * 2)
	if p.Valid() {
		t.Fatalf("Error in .Valid(): Got true expected false\n")
	}

}

func Test_Clear(t *testing.T) {
	content := []byte("Test_Clear")
	p := NewCachedPage(time.Second)
	p.StartUpdate()
	p.Write(content)
	p.EndUpdate()
	if !p.Valid() {
		t.Fatalf("Error in .Valid(): Got false expected true\n")
	}
	p.Clear()
	if p.Valid() {
		t.Fatalf("Error in .Valid(): Got true expected false\n")
	}

}

// advanced concurrent access to find locking issues
func Test_Caching(t *testing.T) {
	p := NewCachedPage(time.Millisecond * 100)
	getCachedContent(&p)
	for clientc := 1; clientc <= 5; clientc++ {
		go testClient(&p, t)
		for i := 1; i < 5; i++ {
			time.Sleep(time.Second * 1)
			r, u := p.GetStatistics()
			if r == 0 || u == 0 {
				t.Fatalf("Lock problem while testing with %d clients: successful reads=%d, contentet updates=%d\n", clientc, r, u)
			}
			p.ClearStatistics()
		}
	}

}

// make continous requests
func testClient(p *CachedPage, t *testing.T) {
	for {
		c := getCachedContent(p)
		if len(c) != 21 {
			t.Fatalf("Error in testClient(): Len(%s)= %d\n", c, len(c))
		}
	}
}

// get content from cache or update it
func getCachedContent(p *CachedPage) []byte {
	if !p.Valid() {
		if p.StartUpdate() == nil {
			time.Sleep(time.Millisecond * 50) // simulated long running backend function
			p.Write([]byte("Test_Caching " + time.Now().Format("15:04:05")))
			p.EndUpdate()
		}
	}
	return p.Get()
}
