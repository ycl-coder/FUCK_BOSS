package shared_test

import (
	"strings"
	"testing"

	"fuck_boss/backend/internal/domain/shared"
)

func TestNewCity_Valid(t *testing.T) {
	tests := []struct {
		testName string
		code     string
		cityName string
	}{
		{
			testName: "normal city",
			code:     "beijing",
			cityName: "北京",
		},
		{
			testName: "with whitespace",
			code:     "  shanghai  ",
			cityName: "  上海  ",
		},
		{
			testName: "english name",
			code:     "newyork",
			cityName: "New York",
		},
		{
			testName: "code with numbers",
			code:     "city001",
			cityName: "City 001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			city, err := shared.NewCity(tt.code, tt.cityName)
			if err != nil {
				t.Fatalf("NewCity() error = %v, want nil", err)
			}

			// Check that whitespace is trimmed
			expectedCode := strings.TrimSpace(tt.code)
			expectedName := strings.TrimSpace(tt.cityName)

			if city.Code() != expectedCode {
				t.Errorf("City.Code() = %v, want %v", city.Code(), expectedCode)
			}

			if city.Name() != expectedName {
				t.Errorf("City.Name() = %v, want %v", city.Name(), expectedName)
			}
		})
	}
}

func TestNewCity_Invalid(t *testing.T) {
	tests := []struct {
		testName string
		code     string
		cityName string
		wantErr  bool
	}{
		{
			testName: "empty code",
			code:     "",
			cityName: "北京",
			wantErr:  true,
		},
		{
			testName: "empty name",
			code:     "beijing",
			cityName: "",
			wantErr:  true,
		},
		{
			testName: "both empty",
			code:     "",
			cityName: "",
			wantErr:  true,
		},
		{
			testName: "code only whitespace",
			code:     "   ",
			cityName: "北京",
			wantErr:  true,
		},
		{
			testName: "name only whitespace",
			code:     "beijing",
			cityName: "   ",
			wantErr:  true,
		},
		{
			testName: "both only whitespace",
			code:     "   ",
			cityName: "   ",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			_, err := shared.NewCity(tt.code, tt.cityName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCity_Code(t *testing.T) {
	code := "beijing"
	name := "北京"
	city, _ := shared.NewCity(code, name)

	if city.Code() != code {
		t.Errorf("City.Code() = %v, want %v", city.Code(), code)
	}
}

func TestCity_Name(t *testing.T) {
	code := "beijing"
	name := "北京"
	city, _ := shared.NewCity(code, name)

	if city.Name() != name {
		t.Errorf("City.Name() = %v, want %v", city.Name(), name)
	}
}

func TestCity_String(t *testing.T) {
	code := "beijing"
	name := "北京"
	city, _ := shared.NewCity(code, name)

	expected := "City{code: beijing, name: 北京}"
	if city.String() != expected {
		t.Errorf("City.String() = %v, want %v", city.String(), expected)
	}
}

func TestCity_IsZero(t *testing.T) {
	tests := []struct {
		testName string
		city     shared.City
		want     bool
	}{
		{
			testName: "zero value",
			city:     shared.City{},
			want:     true,
		},
		{
			testName: "non-zero value",
			city:     func() shared.City { c, _ := shared.NewCity("beijing", "北京"); return c }(),
			want:     false,
		},
		{
			testName: "only code (via NewCity with empty name should fail, but test zero check)",
			city:     func() shared.City { c, _ := shared.NewCity("beijing", "北京"); return c }(),
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			if got := tt.city.IsZero(); got != tt.want {
				t.Errorf("City.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCity_Equals(t *testing.T) {
	code1 := "beijing"
	name1 := "北京"
	city1, _ := shared.NewCity(code1, name1)
	city2, _ := shared.NewCity(code1, name1)
	city3, _ := shared.NewCity("shanghai", "上海")

	if !city1.Equals(city2) {
		t.Error("City.Equals() = false, want true for same code and name")
	}

	if city1.Equals(city3) {
		t.Error("City.Equals() = true, want false for different cities")
	}

	// Test with different code but same name
	city4, _ := shared.NewCity("beijing2", name1)
	if city1.Equals(city4) {
		t.Error("City.Equals() = true, want false for different codes")
	}

	// Test with same code but different name
	city5, _ := shared.NewCity(code1, "北京市")
	if city1.Equals(city5) {
		t.Error("City.Equals() = true, want false for different names")
	}
}

func TestCity_WhitespaceHandling(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		cityName     string
		expectedCode string
		expectedName string
	}{
		{
			name:         "leading whitespace in code",
			code:         "  beijing",
			cityName:     "北京",
			expectedCode: "beijing",
			expectedName: "北京",
		},
		{
			name:         "trailing whitespace in name",
			code:         "beijing",
			cityName:     "北京  ",
			expectedCode: "beijing",
			expectedName: "北京",
		},
		{
			name:         "both sides whitespace",
			code:         "  shanghai  ",
			cityName:     "  上海  ",
			expectedCode: "shanghai",
			expectedName: "上海",
		},
		{
			name:         "internal whitespace preserved",
			code:         "new york",
			cityName:     "New York",
			expectedCode: "new york",
			expectedName: "New York",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			city, err := shared.NewCity(tt.code, tt.cityName)
			if err != nil {
				t.Fatalf("NewCity() error = %v, want nil", err)
			}

			if city.Code() != tt.expectedCode {
				t.Errorf("City.Code() = %v, want %v", city.Code(), tt.expectedCode)
			}

			if city.Name() != tt.expectedName {
				t.Errorf("City.Name() = %v, want %v", city.Name(), tt.expectedName)
			}
		})
	}
}
