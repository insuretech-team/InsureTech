package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStorageResourceExtractor(t *testing.T) {
	t.Parallel()

	extract := StorageResourceExtractor()
	tests := []struct {
		name   string
		method string
		path   string
		want   string
	}{
		{name: "upload", method: http.MethodPost, path: "/v1/storage/files", want: "upload"},
		{name: "list", method: http.MethodGet, path: "/v1/storage/files", want: "get"},
		{name: "upload batch", method: http.MethodPost, path: "/v1/storage/files:batch", want: "upload-batch"},
		{name: "upload url", method: http.MethodPost, path: "/v1/storage/files:upload-url", want: "upload-url"},
		{name: "finalize", method: http.MethodPost, path: "/v1/storage/files:finalize", want: "finalize"},
		{name: "get by id", method: http.MethodGet, path: "/v1/storage/files/abc", want: "get"},
		{name: "update", method: http.MethodPatch, path: "/v1/storage/files/abc", want: "update"},
		{name: "delete", method: http.MethodDelete, path: "/v1/storage/files/abc", want: "delete"},
		{name: "download url post", method: http.MethodPost, path: "/v1/storage/files/abc:download-url", want: "download-url"},
		{name: "download url get", method: http.MethodGet, path: "/v1/storage/files/abc:download-url", want: "download-url"},
		{name: "download url post slash", method: http.MethodPost, path: "/v1/storage/files/abc/download-url", want: "download-url"},
		{name: "download url get slash", method: http.MethodGet, path: "/v1/storage/files/abc/download-url", want: "download-url"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(tc.method, tc.path, nil)
			got := extract(req)
			if got != tc.want {
				t.Fatalf("resource mismatch: got=%q want=%q", got, tc.want)
			}
		})
	}
}
