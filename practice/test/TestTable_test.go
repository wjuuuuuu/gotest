package test

import (
	"net/http"
	"testing"
)

// http Get 기능이 콘텐츠를 다운로드할 수 있는지 확인
func TestTable(t *testing.T) {
	tests := []struct { // 익명구조체
		url        string
		statusCode int
	}{
		{"https://www.google.com", http.StatusOK},
		{"http://rss.cnn.com/rss/cnn_topstorie.rss", http.StatusNotFound},
	}

	t.Logf("Given the need to test downloading different content.")
	{
		for i, tt := range tests {
			t.Logf("\tTest: %d\tWhen checking %q for status code %d", i, tt.url, tt.statusCode)
			{
				resp, err := http.Get(tt.url)
				if err != nil {
					t.Fatalf("\t%s\tShould be able to make the Get call : %v", failed, err)
				}
				t.Logf("\t%s\tShould be able to make the Get call.", succeed)

				defer resp.Body.Close()

				if resp.StatusCode == tt.statusCode {
					t.Logf("\t%s\tShould receive a %d status code.", succeed, tt.statusCode)
				} else {
					t.Errorf("\t%s\tShould receive a %d status code %d.", failed, tt.statusCode, resp.StatusCode)
				}
			}
		}
	}
}
