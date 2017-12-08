package test

import (
	"itchat4go/webservice"
	"testing"
	"fmt"
)

func TestWebService(t *testing.T) {
	webservice.BeginListene()
}

func TestOthers(t *testing.T) {
	letters := []string{"a", "b", "c", "d"}
	for i,_ := range letters {
		if(i % 2== 0) {
			fmt.Println(i)
		}

	}

}