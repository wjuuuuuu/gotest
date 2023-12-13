package test

import (
	"net/http"
	"testing"
)

const (
	succeed = "\u2713"
	failed  = "\u2717"
)

func TestBasic(t *testing.T) {
	url := "http://www.google.com/"
	statusCode := 200

	t.Log("Given the need to test downloading content.") // 로그에 그냥 찍어줌
	{
		t.Logf("\tTest 0:\tWhen checking %q for status code %d", url, statusCode)
		{
			resp, err := http.Get(url)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to make the Get caLL : %v", failed, err) // 로그를 찍어주고 os.Exit(1) 호출하는 것처럼 종료
			}
			t.Logf("\t%s\tShould be able to make the Get call.", succeed)

			defer resp.Body.Close()

			if resp.StatusCode == statusCode {
				t.Logf("\t%s\tShould receive a %d status code.", succeed, statusCode)
			} else {
				t.Errorf("\t%s\tShould receive a %d status code : %d", failed, statusCode, resp.StatusCode) // 테스트 실패를 알리고 코드는 계속 진행함
			}
			t.Logf("Test Complete")

		}
	}
}
