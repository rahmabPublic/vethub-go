package api

import (
	"strings"
	"testing"
)

func TestValidateOwner(t *testing.T) {
	t.Parallel()

	validReq := CreateOwnerRequest{
		FirstName: "William",
		LastName:  "Clown",
		Address:   strPtr("17 Circus Lane"),
		City:      strPtr("Madison"),
		Telephone: strPtr("6085558899"),
		Email:     strPtr("william.clown@example.com"),
	}

	tests := []struct {
		name        string
		req         CreateOwnerRequest
		wantMessage string
	}{
		{name: "valid owner", req: validReq},
		{name: "missing first name", req: CreateOwnerRequest{FirstName: "   ", LastName: "Doe"}, wantMessage: "firstName is required"},
		{name: "first name too long", req: CreateOwnerRequest{FirstName: strings.Repeat("a", 256), LastName: "Doe"}, wantMessage: "firstName must be at most 255 characters"},
		{name: "missing last name", req: CreateOwnerRequest{FirstName: "Jane", LastName: " "}, wantMessage: "lastName is required"},
		{name: "last name too long", req: CreateOwnerRequest{FirstName: "Jane", LastName: strings.Repeat("b", 256)}, wantMessage: "lastName must be at most 255 characters"},
		{name: "address too long", req: CreateOwnerRequest{FirstName: "Jane", LastName: "Doe", Address: strPtr(strings.Repeat("c", 256))}, wantMessage: "address must be at most 255 characters"},
		{name: "city too long", req: CreateOwnerRequest{FirstName: "Jane", LastName: "Doe", City: strPtr(strings.Repeat("d", 256))}, wantMessage: "city must be at most 255 characters"},
		{name: "telephone too long", req: CreateOwnerRequest{FirstName: "Jane", LastName: "Doe", Telephone: strPtr(strings.Repeat("1", 256))}, wantMessage: "telephone must be at most 255 characters"},
		{name: "telephone has non digits", req: CreateOwnerRequest{FirstName: "Jane", LastName: "Doe", Telephone: strPtr("123-456")}, wantMessage: "telephone must contain only digits"},
		{name: "email too long", req: CreateOwnerRequest{FirstName: "Jane", LastName: "Doe", Email: strPtr(strings.Repeat("e", 250) + "@x.com")}, wantMessage: "email must be at most 255 characters"},
		{name: "invalid email", req: CreateOwnerRequest{FirstName: "Jane", LastName: "Doe", Email: strPtr("not-an-email")}, wantMessage: "email must be a valid email address"},
		{name: "empty optional telephone accepted", req: CreateOwnerRequest{FirstName: "Jane", LastName: "Doe", Telephone: strPtr("")}},
		{name: "empty optional email accepted", req: CreateOwnerRequest{FirstName: "Jane", LastName: "Doe", Email: strPtr("")}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			msg := validateOwner(&tt.req)
			if tt.wantMessage == "" {
				if msg != nil {
					t.Fatalf("expected no validation error, got %q", *msg)
				}
				return
			}
			if msg == nil {
				t.Fatalf("expected validation error %q, got nil", tt.wantMessage)
			}
			if *msg != tt.wantMessage {
				t.Fatalf("expected validation error %q, got %q", tt.wantMessage, *msg)
			}
		})
	}
}

func strPtr(v string) *string {
	return &v
}
