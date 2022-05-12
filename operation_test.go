package apidoc

import (
	"testing"
)

func TestOperation_ParseMatchResponseComment(t *testing.T) {
	// responsePattern := regexp.MustCompile(`^(\d+)\s+([\w\.\d_]+\{.*\}|[\w\.\d_\[\]]+)[^"]*(.*)?`)
	tests := []struct {
		name        string
		commentLine string
		wantLen     int
		matches     []string
	}{
		{
			name:        "match response",
			commentLine: `200 common.Response`,
			wantLen:     4,
			matches:     []string{"200", `common.Response`, ""},
		},
		{
			name:        "match response with comment",
			commentLine: `200 common.Response "正常返回"`,
			wantLen:     4,
			matches:     []string{"200", `common.Response`, `"正常返回"`},
		},
		{
			name:        "match response replace data",
			commentLine: `200 common.Response{code=0,msg="success",data=RegisterResponse}`,
			wantLen:     4,
			matches:     []string{"200", `common.Response{code=0,msg="success",data=RegisterResponse}`, ""},
		},
		{
			name:        "match response replace msg with space",
			commentLine: `200 common.Response{code=0,msg="success error",data=RegisterResponse}`,
			wantLen:     4,
			matches:     []string{"200", `common.Response{code=0,msg="success error",data=RegisterResponse}`, ""},
		},
		{
			name:        "match response replace with description",
			commentLine: `200 common.Response{code=0,msg="success error",data=RegisterResponse} "成功返回"`,
			wantLen:     4,
			matches:     []string{"200", `common.Response{code=0,msg="success error",data=RegisterResponse}`, `"成功返回"`},
		},
		{
			name:        "match response array",
			commentLine: `200 []common.Response`,
			wantLen:     4,
			matches:     []string{"200", `[]common.Response`, ``},
		},
		{
			name:        "match response array with description",
			commentLine: `200 []common.Response "测试"`,
			wantLen:     4,
			matches:     []string{"200", `[]common.Response`, `"测试"`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := responsePattern.FindStringSubmatch(tt.commentLine)
			matchesLen := len(matches)
			if matchesLen != tt.wantLen {
				// for _, m := range matches {
				// 	fmt.Println(m)
				// }

				t.Errorf("%s len(matches) = %v, wantLen %v, matches = %v", t.Name(), matchesLen, tt.wantLen, matches)
			}
			for i, m := range matches {
				if i == 0 {
					continue
				}
				index := i - 1
				if m != tt.matches[index] {
					t.Errorf("%s match = %v, want = %v", t.Name(), m, tt.matches[index])
				}
			}
		})
	}
}
