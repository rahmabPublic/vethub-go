package api

import "testing"

func TestValidateOwnerPicture_AcceptsSupportedImageTypesWithinLimit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		contentType string
		sizeBytes   int64
	}{
		{name: "jpeg", contentType: "image/jpeg", sizeBytes: 1024},
		{name: "jpg with charset", contentType: "image/jpeg; charset=binary", sizeBytes: 2048},
		{name: "png", contentType: "image/png", sizeBytes: 5 * 1024 * 1024},
		{name: "webp", contentType: "image/webp", sizeBytes: 4096},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if msg := validateOwnerPicture(tt.contentType, tt.sizeBytes); msg != nil {
				t.Fatalf("expected valid picture, got validation error: %q", *msg)
			}
		})
	}
}

func TestValidateOwnerPicture_RejectsUnsupportedOrInvalidFiles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		contentType string
		sizeBytes   int64
		wantMessage string
	}{
		{
			name:        "unsupported type gif",
			contentType: "image/gif",
			sizeBytes:   1024,
			wantMessage: "picture must be a JPG, PNG, or WebP image",
		},
		{
			name:        "empty content type",
			contentType: "",
			sizeBytes:   1024,
			wantMessage: "picture must be a JPG, PNG, or WebP image",
		},
		{
			name:        "empty file",
			contentType: "image/png",
			sizeBytes:   0,
			wantMessage: "picture file must not be empty",
		},
		{
			name:        "too large",
			contentType: "image/webp",
			sizeBytes:   5*1024*1024 + 1,
			wantMessage: "picture file size must be at most 5 MB",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			msg := validateOwnerPicture(tt.contentType, tt.sizeBytes)
			if msg == nil {
				t.Fatalf("expected validation error %q, got nil", tt.wantMessage)
			}
			if *msg != tt.wantMessage {
				t.Fatalf("expected validation error %q, got %q", tt.wantMessage, *msg)
			}
		})
	}
}
