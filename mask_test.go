package mask

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"

	"reflect"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func Example() {
	rand.Seed(12345)
	type Address struct {
		PostCode string `mask:"zero"`
	}
	type User struct {
		ID      string
		Name    string `mask:"filled"`
		Age     int    `mask:"random100"`
		Address Address
	}

	user := User{
		ID:   "123456",
		Name: "Usagi",
		Age:  3,
		Address: Address{
			PostCode: "123-4567",
		},
	}
	maskUser, err := Mask(user)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v", maskUser)

	// Output:
	// {ID:123456 Name:***** Age:83 Address:{PostCode:}}
}

func ExampleRegisterMaskField() {
	rand.Seed(12345)
	type User2 struct {
		ID      string
		Name    string
		Age     int
		ExtData map[string]string
	}
	user := User2{
		ID:   "123456",
		Name: "Usagi",
		Age:  3,
		ExtData: map[string]string{
			"ID":       "123456",
			"Favorite": "Cat",
		},
	}

	RegisterMaskField("ID", "zero")
	RegisterMaskField("Age", "random100")
	RegisterMaskField("Name", "filled4")
	RegisterMaskField("Favorite", "filled6")
	maskUser, err := Mask(user)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v", maskUser)

	// Output:
	// {ID: Name:**** Age:83 ExtData:map[Favorite:****** ID:]}
}

type BenchTarget struct {
	I  int    `mask:"zero"`
	S  string `mask:"filled"`
	M  map[string]string
	SS []string  `mask:"filled"`
	IS []int     `mask:"rondom100"`
	FS []float64 `mask:"rondom100"`
	B  *BenchTarget2
}

type BenchTarget2 struct {
	I  int       `mask:"random100"`
	S  string    `mask:"fixed"`
	SS []string  `mask:"filled"`
	IS []int     `mask:"rondom100"`
	FS []float64 `mask:"rondom100"`
}

func BenchmarkMask(b *testing.B) {
	RegisterMaskField("Hoge", MaskTypeFixed)
	RegisterMaskField("Bob", MaskTypeFilled+"4")
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v := BenchTarget{
			I: 1,
			S: "Hello World",
			M: map[string]string{
				"Hoge": "Fuga",
				"Bob":  "Alica",
			},
			SS: []string{
				"One",
				"Two",
				"Three",
			},
			IS: []int{
				1,
				2,
				3,
			},
			FS: []float64{
				1.0,
				2.0,
				3.0,
			},
			B: &BenchTarget2{
				I: 2,
				S: "Hello World2",
				SS: []string{
					"One",
					"Two",
					"Three",
				},
				IS: []int{
					1,
					2,
					3,
				},
				FS: []float64{
					1,
					2,
					3,
				},
			},
		}
		Mask(v)
	}
}

func TestMask_PrimitiveType(t *testing.T) {
	tests := map[string]struct {
		input any
		want  any
	}{
		"string": {
			input: "ヤハッ！",
			want:  "ヤハッ！",
		},
		"empty string": {
			input: "",
			want:  "",
		},
		"true": {
			input: true,
			want:  true,
		},
		"false": {
			input: false,
			want:  false,
		},
		"int16": {
			input: int16(2019),
			want:  int16(2019),
		},
		"int32": {
			input: int32(20190122),
			want:  int32(20190122),
		},
		"int64": {
			input: int64(20190122),
			want:  int64(20190122),
		},
		"int": {
			input: 20190122,
			want:  20190122,
		},
		"zero int": {
			input: 0,
			want:  0,
		},
		"float32": {
			input: 20190122,
			want:  20190122,
		},
		"float64": {
			input: 20190122,
			want:  20190122,
		},
		"zero float64": {
			input: 0,
			want:  0,
		},
		"byte": {
			input: byte(12),
			want:  byte(12),
		},
		"uint": {
			input: uint(202),
			want:  uint(202),
		},
		"string pointer": {
			input: convertStringPtr("ヤハッ！"),
			want:  convertStringPtr("ヤハッ！"),
		},
		"int16 pointer": {
			input: convertInt16Ptr(2019),
			want:  convertInt16Ptr(2019),
		},
		"int32 pointer": {
			input: convertInt32Ptr(20190122),
			want:  convertInt32Ptr(20190122),
		},
		"int64 pointer": {
			input: convertInt64Ptr(20190122),
			want:  convertInt64Ptr(20190122),
		},
		"int pointer": {
			input: convertIntPtr(20190122),
			want:  convertIntPtr(20190122),
		},
		"float32 pointer": {
			input: convertFloat32Ptr(20190122),
			want:  convertFloat32Ptr(20190122),
		},
		"float64 pointer": {
			input: convertFloat64Ptr(20190122),
			want:  convertFloat64Ptr(20190122),
		},
		"nil string pointer": {
			input: (*string)(nil),
			want:  (*string)(nil),
		},
		"nil int pointer": {
			input: (*int)(nil),
			want:  (*int)(nil),
		},
		"nil float64 pointer": {
			input: (*float64)(nil),
			want:  (*float64)(nil),
		},
	}

	for name, tt := range tests {
		t.Run(defaultTestCase(name), func(t *testing.T) {
			defer cleanup(t)
			got, err := Mask(tt.input)
			if assert.NoError(t, err) {
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Error(diff)
				}
			}
		})
		t.Run(newMaskerTestCase(name), func(t *testing.T) {
			m := newMasker()
			got, err := m.Mask(tt.input)
			if assert.NoError(t, err) {
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Error(diff)
				}
			}
		})
	}
}

func TestMask_CompositeType(t *testing.T) {
	type stringTest struct {
		Usagi string
	}
	type stringPtrTest struct {
		Usagi *string
	}
	type stringSliceTest struct {
		Usagi []string
	}
	type stringArrayTest struct {
		Usagi [2]string
	}
	type stringSlicePtrTest struct {
		Usagi *[]string
	}
	type intTest struct {
		Usagi int
	}
	type intPtrTest struct {
		Usagi *int
	}
	type intSliceTest struct {
		Usagi []int
	}
	type intArrayTest struct {
		Usagi [3]int
	}
	type intSlicePtrTest struct {
		Usagi *[]int
	}
	type uintTest struct {
		Usagi uint
	}
	type uintSliceTest struct {
		Usagi []uint
	}
	type uintArrayTest struct {
		Usagi [3]uint
	}
	type float64Test struct {
		Usagi float64
	}
	type float64PtrTest struct {
		Usagi *float64
	}
	type float64SliceTest struct {
		Usagi []float64
	}
	type float64ArrayTest struct {
		Usagi [3]float64
	}
	type float64SlicePtrTest struct {
		Usagi *[]float64
	}
	type boolTest struct {
		Usagi bool
	}
	type boolPtrTest struct {
		Usagi *bool
	}
	type byteTest struct {
		Usagi byte
	}
	type byteSliceTest struct {
		Usagi []byte
	}
	type byteArrayTest struct {
		Usagi [4]byte
	}
	type mapStringToStringTest struct {
		Usagi map[string]string
	}
	type mapStringToStringSliceTest struct {
		Usagi []map[string]string
	}
	type mapStringToIntTest struct {
		Usagi map[string]int
	}
	type mapStringToFloat64Test struct {
		Usagi map[string]float64
	}
	type mapIntToStringTest struct {
		Usagi map[int]string
	}
	type mapIntToIntTest struct {
		Usagi map[int]int
	}
	type mapIntToFloat64Test struct {
		Usagi map[int]float64
	}
	type mapStructToStringTest struct {
		Usagi map[stringTest]string
	}
	type mapStructToIntTest struct {
		Usagi map[stringTest]int
	}
	type mapStructToFloat64Test struct {
		Usagi map[stringTest]float64
	}
	type mapStructToStructTest struct {
		Usagi map[stringTest]stringTest
	}
	type structTest struct {
		StringTest      stringTest
		StringSliceTest stringSliceTest
	}
	type structSliceTest struct {
		SliceTest []stringTest
	}
	type unexportedTest struct {
		usagi string
	}
	type anyTest struct {
		Usagi any
	}
	type anyPtrTest struct {
		Usagi *any
	}
	type anySliceTest struct {
		Usagis []any
	}

	tests := map[string]struct {
		input any
		want  any
	}{
		"string fields": {
			input: &stringTest{Usagi: "ヤハッ！"},
			want:  &stringTest{Usagi: "ヤハッ！"},
		},
		"string empty fields": {
			input: &stringTest{},
			want:  &stringTest{Usagi: ""},
		},
		"string slice": {
			input: []string{"ハァ？", "ウラ", "フゥン"},
			want:  []string{"ハァ？", "ウラ", "フゥン"},
		},
		"string array": {
			input: [2]string{"ハァ？", "ウラ"},
			want:  [2]string{"ハァ？", "ウラ"},
		},
		"nil string slice": {
			input: ([]string)(nil),
			want:  ([]string)(nil),
		},
		"string slice ptr": {
			input: convertStringSlicePtr([]string{"ハァ？", "ウラ", "フゥン"}),
			want:  convertStringSlicePtr([]string{"ハァ？", "ウラ", "フゥン"}),
		},
		"nil string slice ptr": {
			input: (*[]string)(nil),
			want:  (*[]string)(nil),
		},
		"map string to string": {
			input: map[string]string{"うさぎ": "ハァ？", "うさぎ2": "ウラ", "うさぎ3": "フゥン"},
			want:  map[string]string{"うさぎ": "ハァ？", "うさぎ2": "ウラ", "うさぎ3": "フゥン"},
		},
		"nil map string to string": {
			input: (map[string]string)(nil),
			want:  (map[string]string)(nil),
		},
		"string ptr fields": {
			input: &stringPtrTest{Usagi: convertStringPtr("ヤハッ！")},
			want:  &stringPtrTest{Usagi: convertStringPtr("ヤハッ！")},
		},
		"nil string ptr fields": {
			input: &stringPtrTest{},
			want:  &stringPtrTest{Usagi: nil},
		},
		"string slice fields": {
			input: &stringSliceTest{Usagi: []string{"ハァ？", "ウラ", "フゥン"}},
			want:  &stringSliceTest{Usagi: []string{"ハァ？", "ウラ", "フゥン"}},
		},
		"string array fields": {
			input: &stringArrayTest{Usagi: [2]string{"ハァ？", "ウラ"}},
			want:  &stringArrayTest{Usagi: [2]string{"ハァ？", "ウラ"}},
		},
		"nil string slice fields": {
			input: &stringSliceTest{},
			want:  &stringSliceTest{Usagi: ([]string)(nil)},
		},
		"string slice ptr fields": {
			input: &stringSlicePtrTest{Usagi: convertStringSlicePtr([]string{"ハァ？", "ウラ", "フゥン"})},
			want:  &stringSlicePtrTest{Usagi: convertStringSlicePtr([]string{"ハァ？", "ウラ", "フゥン"})},
		},
		"nil string slice ptr fields": {
			input: &stringSlicePtrTest{},
			want:  &stringSlicePtrTest{Usagi: (*[]string)(nil)},
		},
		"int slice": {
			input: []int{20190122, 20200501, 20200501},
			want:  []int{20190122, 20200501, 20200501},
		},
		"int array": {
			input: [3]int{20190122, 20200501, 20200501},
			want:  [3]int{20190122, 20200501, 20200501},
		},
		"nil int slice": {
			input: ([]int)(nil),
			want:  ([]int)(nil),
		},
		"int slice ptr": {
			input: convertIntSlicePtr([]int{20190122, 20200501, 20200501}),
			want:  convertIntSlicePtr([]int{20190122, 20200501, 20200501}),
		},
		"nil int slice ptr": {
			input: (*[]int)(nil),
			want:  (*[]int)(nil),
		},
		"int fields": {
			input: &intTest{Usagi: 20190122},
			want:  &intTest{Usagi: 20190122},
		},
		"zero int fields": {
			input: &intTest{},
			want:  &intTest{Usagi: 0},
		},
		"int ptr fields": {
			input: &intPtrTest{Usagi: convertIntPtr(20190122)},
			want:  &intPtrTest{Usagi: convertIntPtr(20190122)},
		},
		"nil int ptr fields": {
			input: &intPtrTest{},
			want:  &intPtrTest{Usagi: nil},
		},
		"int slice fields": {
			input: &intSliceTest{Usagi: []int{20190122, 20200501, 20200501}},
			want:  &intSliceTest{Usagi: []int{20190122, 20200501, 20200501}},
		},
		"int array fields": {
			input: &intArrayTest{Usagi: [3]int{20190122, 20200501, 20200501}},
			want:  &intArrayTest{Usagi: [3]int{20190122, 20200501, 20200501}},
		},
		"nil int slice fields": {
			input: &intSliceTest{},
			want:  &intSliceTest{Usagi: ([]int)(nil)},
		},
		"int slice ptr fields": {
			input: &intSlicePtrTest{Usagi: convertIntSlicePtr([]int{20190122, 20200501, 20200501})},
			want:  &intSlicePtrTest{Usagi: convertIntSlicePtr([]int{20190122, 20200501, 20200501})},
		},
		"nil int slice ptr fields": {
			input: &intSlicePtrTest{},
			want:  &intSlicePtrTest{Usagi: (*[]int)(nil)},
		},
		"uint fields": {
			input: &uintTest{Usagi: 20190122},
			want:  &uintTest{Usagi: 20190122},
		},
		"uint slice fields": {
			input: &uintSliceTest{Usagi: []uint{20190122, 20200501, 20200501}},
			want:  &uintSliceTest{Usagi: []uint{20190122, 20200501, 20200501}},
		},
		"uint array fields": {
			input: &uintArrayTest{Usagi: [3]uint{20190122, 20200501, 20200501}},
			want:  &uintArrayTest{Usagi: [3]uint{20190122, 20200501, 20200501}},
		},
		"float64 slice": {
			input: []float64{20190122, 20200501, 20200501},
			want:  []float64{20190122, 20200501, 20200501},
		},
		"float64 array": {
			input: [3]float64{20190122, 20200501, 20200501},
			want:  [3]float64{20190122, 20200501, 20200501},
		},
		"nil float64 slice": {
			input: ([]float64)(nil),
			want:  ([]float64)(nil),
		},
		"float64 slice ptr": {
			input: convertFloat64SlicePtr([]float64{20190122, 20200501, 20200501}),
			want:  convertFloat64SlicePtr([]float64{20190122, 20200501, 20200501}),
		},
		"nil float64 slice ptr": {
			input: (*[]float64)(nil),
			want:  (*[]float64)(nil),
		},
		"float64 fields": {
			input: &float64Test{Usagi: 20190122},
			want:  &float64Test{Usagi: 20190122},
		},
		"zero float64 fields": {
			input: &float64Test{},
			want:  &float64Test{Usagi: 0},
		},
		"float64 ptr fields": {
			input: &float64PtrTest{Usagi: convertFloat64Ptr(20190122)},
			want:  &float64PtrTest{Usagi: convertFloat64Ptr(20190122)},
		},
		"nil float64 ptr fields": {
			input: &float64PtrTest{},
			want:  &float64PtrTest{Usagi: nil},
		},
		"float64 slice fields": {
			input: &float64SliceTest{Usagi: []float64{20190122, 20200501, 20200501}},
			want:  &float64SliceTest{Usagi: []float64{20190122, 20200501, 20200501}},
		},
		"float64 array fields": {
			input: &float64ArrayTest{Usagi: [3]float64{20190122, 20200501, 20200501}},
			want:  &float64ArrayTest{Usagi: [3]float64{20190122, 20200501, 20200501}},
		},
		"nil float64 slice fields": {
			input: &float64SliceTest{},
			want:  &float64SliceTest{Usagi: ([]float64)(nil)},
		},
		"float64 slice ptr fields": {
			input: &float64SlicePtrTest{Usagi: convertFloat64SlicePtr([]float64{20190122, 20200501, 20200501})},
			want:  &float64SlicePtrTest{Usagi: convertFloat64SlicePtr([]float64{20190122, 20200501, 20200501})},
		},
		"nil float64 slice ptr fields": {
			input: &float64SlicePtrTest{},
			want:  &float64SlicePtrTest{Usagi: (*[]float64)(nil)},
		},
		"bool fields": {
			input: &boolTest{Usagi: true},
			want:  &boolTest{Usagi: true},
		},
		"zero bool fields": {
			input: &boolTest{},
			want:  &boolTest{Usagi: false},
		},
		"bool ptr fields": {
			input: &boolPtrTest{Usagi: convertBoolPtr(true)},
			want:  &boolPtrTest{Usagi: convertBoolPtr(true)},
		},
		"nil bool ptr fields": {
			input: &boolPtrTest{},
			want:  &boolPtrTest{Usagi: (*bool)(nil)},
		},
		"byte slice": {
			input: []byte{10, 11, 12},
			want:  []byte{10, 11, 12},
		},
		"byte array": {
			input: [4]byte{10, 11, 12, 13},
			want:  [4]byte{10, 11, 12, 13},
		},
		"byte fields": {
			input: byteTest{Usagi: 11},
			want:  byteTest{Usagi: 11},
		},
		"byte slice fields": {
			input: byteSliceTest{Usagi: []byte{11, 12, 13}},
			want:  byteSliceTest{Usagi: []byte{11, 12, 13}},
		},
		"byte array fields": {
			input: byteArrayTest{Usagi: [4]byte{11, 12, 13, 14}},
			want:  byteArrayTest{Usagi: [4]byte{11, 12, 13, 14}},
		},
		"map string to string fields": {
			input: &mapStringToStringTest{Usagi: map[string]string{"うさぎ": "ハァ？", "うさぎ2": "ウラ", "うさぎ3": "フゥン"}},
			want:  &mapStringToStringTest{Usagi: map[string]string{"うさぎ": "ハァ？", "うさぎ2": "ウラ", "うさぎ3": "フゥン"}},
		},
		"map string to string slice fields": {
			input: &mapStringToStringSliceTest{Usagi: []map[string]string{{"うさぎ": "ハァ？", "うさぎ2": "ウラ", "うさぎ3": "フゥン"}}},
			want:  &mapStringToStringSliceTest{Usagi: []map[string]string{{"うさぎ": "ハァ？", "うさぎ2": "ウラ", "うさぎ3": "フゥン"}}},
		},
		"nil map string to string fields": {
			input: &mapStringToStringTest{},
			want:  &mapStringToStringTest{Usagi: (map[string]string)(nil)},
		},
		"map string to int fields": {
			input: &mapStringToIntTest{Usagi: map[string]int{"うさぎ": 20190122, "うさぎ2": 20190122, "うさぎ3": 20190122}},
			want:  &mapStringToIntTest{Usagi: map[string]int{"うさぎ": 20190122, "うさぎ2": 20190122, "うさぎ3": 20190122}},
		},
		"map string to float64 fields": {
			input: &mapStringToFloat64Test{Usagi: map[string]float64{"うさぎ": 20190122, "うさぎ2": 20190122, "うさぎ3": 20190122}},
			want:  &mapStringToFloat64Test{Usagi: map[string]float64{"うさぎ": 20190122, "うさぎ2": 20190122, "うさぎ3": 20190122}},
		},
		"map int to string fields": {
			input: &mapIntToStringTest{Usagi: map[int]string{201901221: "ハァ？", 201901222: "ウラ", 201901223: "フゥン"}},
			want:  &mapIntToStringTest{Usagi: map[int]string{201901221: "ハァ？", 201901222: "ウラ", 201901223: "フゥン"}},
		},
		"map int to int fields": {
			input: &mapIntToIntTest{Usagi: map[int]int{1: 201901221, 2: 201901222, 3: 201901223}},
			want:  &mapIntToIntTest{Usagi: map[int]int{1: 201901221, 2: 201901222, 3: 201901223}},
		},
		"map int to float64 fields": {
			input: &mapIntToFloat64Test{Usagi: map[int]float64{1: 201901221, 2: 201901222, 3: 201901223}},
			want:  &mapIntToFloat64Test{Usagi: map[int]float64{1: 201901221, 2: 201901222, 3: 201901223}},
		},
		"map struct to string fields": {
			input: &mapStructToStringTest{Usagi: map[stringTest]string{{Usagi: "ウサギ１"}: "ハァ？", {Usagi: "ウサギ２"}: "ウラ", {Usagi: "ウサギ３"}: "フゥン"}},
			want:  &mapStructToStringTest{Usagi: map[stringTest]string{{Usagi: "ウサギ１"}: "ハァ？", {Usagi: "ウサギ２"}: "ウラ", {Usagi: "ウサギ３"}: "フゥン"}},
		},
		"map struct to int fields": {
			input: &mapStructToIntTest{Usagi: map[stringTest]int{{Usagi: "ウサギ１"}: 201901221, {Usagi: "ウサギ２"}: 201901222, {Usagi: "ウサギ３"}: 201901223}},
			want:  &mapStructToIntTest{Usagi: map[stringTest]int{{Usagi: "ウサギ１"}: 201901221, {Usagi: "ウサギ２"}: 201901222, {Usagi: "ウサギ３"}: 201901223}},
		},
		"map struct to float64 fields": {
			input: &mapStructToFloat64Test{Usagi: map[stringTest]float64{{Usagi: "ウサギ１"}: 201901221, {Usagi: "ウサギ２"}: 201901222, {Usagi: "ウサギ３"}: 201901223}},
			want:  &mapStructToFloat64Test{Usagi: map[stringTest]float64{{Usagi: "ウサギ１"}: 201901221, {Usagi: "ウサギ２"}: 201901222, {Usagi: "ウサギ３"}: 201901223}},
		},
		"map struct to struct fields": {
			input: &mapStructToStructTest{Usagi: map[stringTest]stringTest{{Usagi: "ウサギ１"}: {Usagi: "Rabbit1"}, {Usagi: "ウサギ２"}: {Usagi: "Rabbit2"}, {Usagi: "ウサギ３"}: {Usagi: "Rabbit3"}}},
			want:  &mapStructToStructTest{Usagi: map[stringTest]stringTest{{Usagi: "ウサギ１"}: {Usagi: "Rabbit1"}, {Usagi: "ウサギ２"}: {Usagi: "Rabbit2"}, {Usagi: "ウサギ３"}: {Usagi: "Rabbit3"}}},
		},
		"struct fields": {
			input: &structTest{
				StringTest:      stringTest{Usagi: "ヤハッ！"},
				StringSliceTest: stringSliceTest{Usagi: []string{"ハァ？", "ウラ", "フゥン"}},
			},
			want: &structTest{
				StringTest:      stringTest{Usagi: "ヤハッ！"},
				StringSliceTest: stringSliceTest{Usagi: []string{"ハァ？", "ウラ", "フゥン"}},
			},
		},
		"zero struct fields": {
			input: &structTest{},
			want: &structTest{
				StringTest:      stringTest{},
				StringSliceTest: stringSliceTest{},
			},
		},
		"struct slice fields": {
			input: &structSliceTest{
				SliceTest: []stringTest{
					{Usagi: "ハァ？"}, {Usagi: "ウラ"}, {Usagi: "フゥン"},
				},
			},
			want: &structSliceTest{
				SliceTest: []stringTest{{
					Usagi: "ハァ？"}, {Usagi: "ウラ"}, {Usagi: "フゥン"},
				},
			},
		},
		"nil struct slice fields": {
			input: &structSliceTest{},
			want:  &structSliceTest{SliceTest: ([]stringTest)(nil)},
		},
		"unexported fields": {
			input: &unexportedTest{usagi: "ヤハッ！"},
			want:  &unexportedTest{},
		},
		"string in any fields": {
			input: anyTest{Usagi: "hoge"},
			want:  anyTest{Usagi: "hoge"},
		},
		"int in any fields": {
			input: anyTest{Usagi: 2},
			want:  anyTest{Usagi: 2},
		},
		"string struct in any fields": {
			input: anyTest{Usagi: stringTest{"hoge"}},
			want:  anyTest{Usagi: stringTest{"hoge"}},
		},
		"int struct in any fields": {
			input: anyTest{Usagi: intTest{2}},
			want:  anyTest{Usagi: intTest{2}},
		},
		"nil in any fields": {
			input: anyTest{Usagi: nil},
			want:  anyTest{Usagi: nil},
		},
		"string slice in any slice fields": {
			input: anySliceTest{Usagis: []any{"hoge", "fuga"}},
			want:  anySliceTest{Usagis: []any{"hoge", "fuga"}},
		},
		"string struct slice in any slice fields": {
			input: anySliceTest{Usagis: []any{stringTest{"hoge"}, stringTest{"fuga"}}},
			want:  anySliceTest{Usagis: []any{stringTest{"hoge"}, stringTest{"fuga"}}},
		},
		"int slice in any slice fields": {
			input: anySliceTest{Usagis: []any{3, 3}},
			want:  anySliceTest{Usagis: []any{3, 3}},
		},
		"int struct slice in any slice fields": {
			input: anySliceTest{Usagis: []any{intTest{3}, intTest{3}}},
			want:  anySliceTest{Usagis: []any{intTest{3}, intTest{3}}},
		},
		"nil in any slice fields": {
			input: anySliceTest{Usagis: nil},
			want:  anySliceTest{Usagis: nil},
		},
		"string struct in any pointer fields": {
			input: anyPtrTest{Usagi: convertAnyPtr(stringTest{"hoge"})},
			want:  anyPtrTest{Usagi: convertAnyPtr(stringTest{"hoge"})},
		},
		"int struct in any pointer fields": {
			input: anyPtrTest{Usagi: convertAnyPtr(intTest{2})},
			want:  anyPtrTest{Usagi: convertAnyPtr(intTest{2})},
		},
	}

	for name, tt := range tests {
		for _, cache := range []bool{true, false} {
			t.Run(defaultTestCase(name), func(t *testing.T) {
				defer cleanup(t)
				defaultMasker.Cache(cache)
				got, err := Mask(tt.input)
				assert.Nil(t, err)
				if diff := cmp.Diff(tt.want, got, allowUnexported(tt.input)); diff != "" {
					t.Error(diff)
				}
			})
			t.Run(newMaskerTestCase(name), func(t *testing.T) {
				m := newMasker()
				m.Cache(cache)
				got, err := m.Mask(tt.input)
				assert.Nil(t, err)
				if diff := cmp.Diff(tt.want, got, allowUnexported(tt.input)); diff != "" {
					t.Error(diff)
				}
			})
		}
	}
}

func TestMask_SameStruct(t *testing.T) {
	// Caching the struct type in sync.Map.
	// If there are different fields with the same struct name in the same package, it will result in an error.
	t.Skip()
	type sameStructNameTest struct {
		Usagi string
	}
	createSameStruct := func(value int) any {
		type sameStructNameTest struct {
			Usagi int
		}
		return sameStructNameTest{value}
	}

	t.Run(defaultTestCase("same struct name"), func(t *testing.T) {
		defer cleanup(t)
		{
			input := sameStructNameTest{"Rabbit"}
			got, err := Mask(input)
			assert.Nil(t, err)
			if diff := cmp.Diff(input, got); diff != "" {
				t.Error(diff)
			}
		}
		{
			input := createSameStruct(2)
			got, err := Mask(input)
			assert.Nil(t, err)
			if diff := cmp.Diff(input, got); diff != "" {
				t.Error(diff)
			}
		}
	})
	t.Run(newMaskerTestCase("same struct name"), func(t *testing.T) {
		m := newMasker()
		{
			input := sameStructNameTest{"Rabbit"}
			got, err := m.Mask(input)
			assert.Nil(t, err)
			if diff := cmp.Diff(input, got); diff != "" {
				t.Error(diff)
			}
		}

		{
			input := createSameStruct(2)
			got, err := m.Mask(input)
			assert.Nil(t, err)
			if diff := cmp.Diff(input, got); diff != "" {
				t.Error(diff)
			}
		}
	})
}

func TestMask_SameAnonynousStruct(t *testing.T) {
	t.Run(defaultTestCase("same anonymous struct name"), func(t *testing.T) {
		defer cleanup(t)
		{
			input := struct {
				Usagi string
			}{
				Usagi: "Rabbit",
			}
			got, err := Mask(input)
			assert.Nil(t, err)
			if diff := cmp.Diff(input, got); diff != "" {
				t.Error(diff)
			}
		}
		{
			input := struct {
				A int
			}{
				A: 2,
			}
			got, err := Mask(input)
			assert.Nil(t, err)
			if diff := cmp.Diff(input, got); diff != "" {
				t.Error(diff)
			}
		}
	})
	t.Run(newMaskerTestCase("same anonymous struct name"), func(t *testing.T) {
		m := newMasker()
		{
			input := struct {
				Usagi string
			}{
				Usagi: "Rabbit",
			}
			got, err := m.Mask(input)
			assert.Nil(t, err)
			if diff := cmp.Diff(input, got); diff != "" {
				t.Error(diff)
			}
		}

		{
			input := struct {
				A int
			}{
				A: 2,
			}
			got, err := m.Mask(input)
			assert.Nil(t, err)
			if diff := cmp.Diff(input, got); diff != "" {
				t.Error(diff)
			}
		}
	})
}

func TestString(t *testing.T) {
	tests := map[string]struct {
		tag   string
		input string
		want  string
	}{
		"no tag": {
			tag:   "",
			input: "ヤハッ！",
			want:  "ヤハッ！",
		},
		"undefined tag": {
			tag:   "usagi!!",
			input: "ヤハッ！",
			want:  "ヤハッ！",
		},
		"filled": {
			tag:   MaskTypeFilled,
			input: "ヤハッ！",
			want:  "****",
		},
		"zero": {
			tag:   MaskTypeZero,
			input: "ヤハッ！",
			want:  "",
		},
	}

	for name, tt := range tests {
		t.Run(defaultTestCase(name), func(t *testing.T) {
			defer cleanup(t)
			got, err := String(tt.tag, tt.input)
			assert.Nil(t, err)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Error(diff)
			}
		})
		t.Run(newMaskerTestCase(name), func(t *testing.T) {
			m := newMasker()
			got, err := m.String(tt.tag, tt.input)
			assert.Nil(t, err)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestInt(t *testing.T) {
	tests := map[string]struct {
		tag     string
		input   int
		want    int
		wantErr bool
	}{
		"no tag": {
			tag:   "",
			input: 20190122,
			want:  20190122,
		},
		"undefined tag": {
			tag:   "usagi!!",
			input: 20190122,
			want:  20190122,
		},
		"randomXX": {
			tag:     MaskTypeRandom + "XX",
			input:   20190122,
			wantErr: true,
		},
		"random30": {
			tag:   MaskTypeRandom + "30",
			input: 20190122,
			want:  9,
		},
		"random1000": {
			tag:   MaskTypeRandom + "1000",
			input: 20190122,
			want:  829,
		},
		"zero": {
			tag:   "zero",
			input: 0,
			want:  0,
		},
	}

	for name, tt := range tests {
		t.Run(defaultTestCase(name), func(t *testing.T) {
			rand.Seed(rand.NewSource(1).Int63())
			defer cleanup(t)
			got, err := Int(tt.tag, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Error(diff)
				}
			}
		})
		t.Run(newMaskerTestCase(name), func(t *testing.T) {
			rand.Seed(rand.NewSource(1).Int63())
			m := newMasker()
			got, err := m.Int(tt.tag, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Error(diff)
				}
			}
		})
	}
}

func TestUint(t *testing.T) {
	tests := map[string]struct {
		tag     string
		input   uint
		want    uint
		wantErr bool
	}{
		"no tag": {
			tag:   "",
			input: 20190122,
			want:  20190122,
		},
		"undefined tag": {
			tag:   "usagi!!",
			input: 20190122,
			want:  20190122,
		},
		"randomXX": {
			tag:     "urandomXX",
			input:   20190122,
			wantErr: true,
		},
		"random30": {
			tag:   "urandom30",
			input: 20190122,
			want:  9,
		},
		"random1000": {
			tag:   "urandom1000",
			input: 20190122,
			want:  829,
		},
		"zero": {
			tag:   "zero",
			input: 0,
			want:  0,
		},
	}

	for name, tt := range tests {
		t.Run(defaultTestCase(name), func(t *testing.T) {
			rand.Seed(rand.NewSource(1).Int63())
			defer cleanup(t)
			RegisterMaskUintFunc("urandom", func(arg string, value uint) (uint, error) {
				ret, err := defaultMasker.MaskRandomInt(arg, int(value))
				return uint(ret), err
			})
			got, err := Uint(tt.tag, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Error(diff)
				}
			}
		})
		t.Run(newMaskerTestCase(name), func(t *testing.T) {
			rand.Seed(rand.NewSource(1).Int63())
			m := newMasker()
			m.RegisterMaskUintFunc("urandom", func(arg string, value uint) (uint, error) {
				ret, err := defaultMasker.MaskRandomInt(arg, int(value))
				return uint(ret), err
			})
			got, err := m.Uint(tt.tag, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Error(diff)
				}
			}
		})
	}
}

func TestFloat64(t *testing.T) {
	tests := map[string]struct {
		tag     string
		input   float64
		want    float64
		wantErr bool
	}{
		"no tag": {
			tag:   "",
			input: 20190122,
			want:  20190122,
		},
		"undefined tag": {
			tag:   "usagi!!",
			input: 20190122,
			want:  20190122,
		},
		"randomXX.4": {
			tag:     MaskTypeRandom + "XX.4",
			input:   20190122,
			wantErr: true,
		},
		"random4.XX": {
			tag:     MaskTypeRandom + "10000.XX",
			input:   20190122,
			wantErr: true,
		},
		"random5.4": {
			tag:   MaskTypeRandom + "100000.4",
			input: 20190122,
			want:  96011.8989,
		},
		"random1.1": {
			tag:   MaskTypeRandom + "10.1",
			input: 20190122,
			want:  9.6,
		},
		"random1": {
			tag:   MaskTypeRandom + "10",
			input: 20190122,
			want:  9.0,
		},
		"zero": {
			tag:   "zero",
			input: 20190122,
			want:  0,
		},
	}

	for name, tt := range tests {
		t.Run(defaultTestCase(name), func(t *testing.T) {
			rand.Seed(rand.NewSource(1).Int63())
			defer cleanup(t)
			got, err := Float64(tt.tag, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Error(diff)
				}
			}
		})
		t.Run(newMaskerTestCase(name), func(t *testing.T) {
			rand.Seed(rand.NewSource(1).Int63())
			m := newMasker()
			got, err := m.Float64(tt.tag, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Error(diff)
				}
			}
		})
	}
}

func TestMaskFilled(t *testing.T) {
	type stringTest struct {
		Usagi string `mask:"filled"`
	}
	type stringMask5Test struct {
		Usagi string `mask:"filled5"`
	}
	type stringPtrTest struct {
		Usagi *string `mask:"filled"`
	}
	type stringPtrMask8Test struct {
		Usagi *string `mask:"filled8"`
	}
	type stringSliceTest struct {
		Usagi []string `mask:"filled"`
	}
	type stringSlicePtrTest struct {
		Usagi *[]string `mask:"filled"`
	}
	type stringToStringMapTest struct {
		Usagi map[string]string `mask:"filled"`
	}
	type intToStringMapTest struct {
		Usagi map[int]string `mask:"filled"`
	}
	type structToStringMapTest struct {
		Usagi map[stringTest]string `mask:"filled"`
	}

	tests := map[string]struct {
		input any
		want  any
	}{
		"string fields": {
			input: &stringTest{Usagi: "ヤハッ！"},
			want:  &stringTest{Usagi: "****"},
		},
		"zero string fields": {
			input: &stringTest{},
			want:  &stringTest{Usagi: ""},
		},
		"string ptr fields": {
			input: &stringPtrTest{Usagi: convertStringPtr("ヤハッ！")},
			want:  &stringPtrTest{Usagi: convertStringPtr("****")},
		},
		"nil string ptr fields": {
			input: &stringPtrTest{},
			want:  &stringPtrTest{Usagi: (*string)(nil)},
		},
		"string slice fields": {
			input: &stringSliceTest{Usagi: []string{"ハァ？", "ウラ", "フゥン"}},
			want:  &stringSliceTest{Usagi: []string{"***", "**", "***"}},
		},
		"nil string slice fields": {
			input: &stringSliceTest{},
			want:  &stringSliceTest{Usagi: ([]string)(nil)},
		},
		"string slice ptr fields": {
			input: &stringSlicePtrTest{Usagi: convertStringSlicePtr([]string{"ハァ？", "ウラ", "フゥン"})},
			want:  &stringSlicePtrTest{Usagi: convertStringSlicePtr([]string{"***", "**", "***"})},
		},
		"nil string slice ptr fields": {
			input: &stringSlicePtrTest{},
			want:  &stringSlicePtrTest{Usagi: (*[]string)(nil)},
		},
		"string to string map fields": {
			input: &stringToStringMapTest{Usagi: map[string]string{"うさぎ": "ハァ？", "うさぎ2": "ウラ", "うさぎ3": "フゥン"}},
			want:  &stringToStringMapTest{Usagi: map[string]string{"うさぎ": "***", "うさぎ2": "**", "うさぎ3": "***"}},
		},
		"int to string map fields": {
			input: &intToStringMapTest{Usagi: map[int]string{1: "ハァ？", 2: "ウラ", 3: "フゥン"}},
			want:  &intToStringMapTest{Usagi: map[int]string{1: "***", 2: "**", 3: "***"}},
		},
		"struct to string map fields": {
			input: &structToStringMapTest{Usagi: map[stringTest]string{{Usagi: "ヤハッ！"}: "ハァ？", {Usagi: "ヤハッ！！"}: "ウラ", {Usagi: "ヤハッ！！！"}: "フゥン"}},
			want:  &structToStringMapTest{Usagi: map[stringTest]string{{Usagi: "ヤハッ！"}: "***", {Usagi: "ヤハッ！！"}: "**", {Usagi: "ヤハッ！！！"}: "***"}},
		},
		"filled 5 chars": {
			input: stringMask5Test{Usagi: "ヤハッ！"},
			want:  stringMask5Test{Usagi: "*****"},
		},
		"filled 8 chars": {
			input: stringPtrMask8Test{Usagi: convertStringPtr("ヤハッ！")},
			want:  stringPtrMask8Test{Usagi: convertStringPtr("********")},
		},
	}

	for name, tt := range tests {
		t.Run(defaultTestCase(name), func(t *testing.T) {
			defer cleanup(t)
			got, err := Mask(tt.input)
			assert.Nil(t, err)
			if diff := cmp.Diff(tt.want, got, allowUnexported(tt.input)); diff != "" {
				t.Error(diff)
			}
		})
		t.Run(newMaskerTestCase(name), func(t *testing.T) {
			m := newMasker()
			got, err := m.Mask(tt.input)
			assert.Nil(t, err)
			if diff := cmp.Diff(tt.want, got, allowUnexported(tt.input)); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestMaskFixed(t *testing.T) {
	type stringTest struct {
		Usagi string `mask:"fixed"`
	}

	tests := map[string]struct {
		input any
		want  any
	}{
		"string fields": {
			input: &stringTest{Usagi: "ヤハッ！！！"},
			want:  &stringTest{Usagi: "********"},
		},
		"zero string fields": {
			input: &stringTest{},
			want:  &stringTest{Usagi: ""},
		},
	}

	for name, tt := range tests {
		t.Run(defaultTestCase(name), func(t *testing.T) {
			defer cleanup(t)
			got, err := Mask(tt.input)
			assert.Nil(t, err)
			if diff := cmp.Diff(tt.want, got, allowUnexported(tt.input)); diff != "" {
				t.Error(diff)
			}
		})
		t.Run(newMaskerTestCase(name), func(t *testing.T) {
			m := newMasker()
			got, err := m.Mask(tt.input)
			assert.Nil(t, err)
			if diff := cmp.Diff(tt.want, got, allowUnexported(tt.input)); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestMaskHashString(t *testing.T) {
	type stringTest struct {
		Usagi string `mask:"hash"`
	}
	type stringPtrTest struct {
		Usagi *string `mask:"hash"`
	}
	type stringSliceTest struct {
		Usagi []string `mask:"hash"`
	}
	type stringArrayTest struct {
		Usagi [3]string `mask:"hash"`
	}
	type stringSlicePtrTest struct {
		Usagi *[]string `mask:"hash"`
	}
	type stringToStringMapTest struct {
		Usagi map[string]string `mask:"hash"`
	}
	type intToStringMapTest struct {
		Usagi map[int]string `mask:"hash"`
	}
	type structToStringMapTest struct {
		Usagi map[stringTest]string `mask:"hash"`
	}

	tests := map[string]struct {
		input any
		want  any
	}{
		"string fields": {
			input: &stringTest{Usagi: "ヤハッ！"},
			want:  &stringTest{Usagi: "a6ab5728db57954641b2e155adc61f2cbdfc7063"},
		},
		"zero string fields": {
			input: &stringTest{},
			want:  &stringTest{Usagi: ""},
		},
		"string ptr fields": {
			input: &stringPtrTest{Usagi: convertStringPtr("ヤハッ！")},
			want:  &stringPtrTest{Usagi: convertStringPtr("a6ab5728db57954641b2e155adc61f2cbdfc7063")},
		},
		"nil string ptr fields": {
			input: &stringPtrTest{},
			want:  &stringPtrTest{Usagi: (*string)(nil)},
		},
		"string slice fields": {
			input: &stringSliceTest{Usagi: []string{"ハァ？", "ウラ", "フゥン"}},
			want: &stringSliceTest{Usagi: []string{
				"48a8b33f36a35631f584844686adaba89a6f156a",
				"ecef3e43f07f7150c089e99d5e1041259b1189d5",
				"17fa078ad3f2c34c17ee58b9119963548ddcf1ef",
			}},
		},
		"string array fields": {
			input: &stringArrayTest{Usagi: [3]string{"ハァ？", "ウラ", "フゥン"}},
			want: &stringArrayTest{Usagi: [3]string{
				"48a8b33f36a35631f584844686adaba89a6f156a",
				"ecef3e43f07f7150c089e99d5e1041259b1189d5",
				"17fa078ad3f2c34c17ee58b9119963548ddcf1ef",
			}},
		},
		"nil string slice fields": {
			input: &stringSliceTest{},
			want:  &stringSliceTest{Usagi: ([]string)(nil)},
		},
		"string slice ptr fields": {
			input: &stringSlicePtrTest{Usagi: convertStringSlicePtr([]string{"ハァ？", "ウラ", "フゥン"})},
			want: &stringSlicePtrTest{Usagi: convertStringSlicePtr([]string{
				"48a8b33f36a35631f584844686adaba89a6f156a",
				"ecef3e43f07f7150c089e99d5e1041259b1189d5",
				"17fa078ad3f2c34c17ee58b9119963548ddcf1ef",
			})},
		},
		"nil string slice ptr fields": {
			input: &stringSlicePtrTest{},
			want:  &stringSlicePtrTest{Usagi: (*[]string)(nil)},
		},
		"string to string map fields": {
			input: &stringToStringMapTest{Usagi: map[string]string{"うさぎ": "ハァ？", "うさぎ2": "ウラ", "うさぎ3": "フゥン"}},
			want: &stringToStringMapTest{Usagi: map[string]string{
				"うさぎ":  "48a8b33f36a35631f584844686adaba89a6f156a",
				"うさぎ2": "ecef3e43f07f7150c089e99d5e1041259b1189d5",
				"うさぎ3": "17fa078ad3f2c34c17ee58b9119963548ddcf1ef",
			}},
		},
		"int to string map fields": {
			input: &intToStringMapTest{Usagi: map[int]string{1: "ハァ？", 2: "ウラ", 3: "フゥン"}},
			want: &intToStringMapTest{Usagi: map[int]string{
				1: "48a8b33f36a35631f584844686adaba89a6f156a",
				2: "ecef3e43f07f7150c089e99d5e1041259b1189d5",
				3: "17fa078ad3f2c34c17ee58b9119963548ddcf1ef",
			}},
		},
		"struct to string map fields": {
			input: &structToStringMapTest{Usagi: map[stringTest]string{{Usagi: "ヤハッ！"}: "ハァ？", {Usagi: "ヤハッ！！"}: "ウラ", {Usagi: "ヤハッ！！！"}: "フゥン"}},
			want: &structToStringMapTest{Usagi: map[stringTest]string{{
				Usagi: "ヤハッ！"}: "48a8b33f36a35631f584844686adaba89a6f156a",
				{Usagi: "ヤハッ！！"}:  "ecef3e43f07f7150c089e99d5e1041259b1189d5",
				{Usagi: "ヤハッ！！！"}: "17fa078ad3f2c34c17ee58b9119963548ddcf1ef",
			}},
		},
	}

	for name, tt := range tests {
		t.Run(defaultTestCase(name), func(t *testing.T) {
			defer cleanup(t)
			got, err := Mask(tt.input)
			assert.Nil(t, err)

			if diff := cmp.Diff(tt.want, got, allowUnexported(tt.input)); diff != "" {
				t.Error(diff)
			}
		})
		t.Run(newMaskerTestCase(name), func(t *testing.T) {
			m := newMasker()
			got, err := m.Mask(tt.input)
			assert.Nil(t, err)

			if diff := cmp.Diff(tt.want, got, allowUnexported(tt.input)); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestMaskRandom(t *testing.T) {
	type intTest struct {
		Usagi int `mask:"random1000"`
	}
	type int16Test struct {
		Usagi int32 `mask:"random1000"`
	}
	type int32Test struct {
		Usagi int32 `mask:"random1000"`
	}
	type int64Test struct {
		Usagi int64 `mask:"random1000"`
	}
	type intPtrTest struct {
		Usagi *int `mask:"random1000"`
	}
	type intSliceTest struct {
		Usagi []int `mask:"random1000"`
	}
	type int32SliceTest struct {
		Usagi []int32 `mask:"random1000"`
	}
	type int64SliceTest struct {
		Usagi []int64 `mask:"random1000"`
	}
	type intArrayTest struct {
		Usagi [2]int `mask:"random1000"`
	}
	type int32ArrayTest struct {
		Usagi [2]int32 `mask:"random1000"`
	}
	type int64ArrayTest struct {
		Usagi [2]int64 `mask:"random1000"`
	}
	type intSlicePtrTest struct {
		Usagi *[]int `mask:"random1000"`
	}
	type float32Test struct {
		Usagi float32 `mask:"random100000.4"`
	}
	type float64Test struct {
		Usagi float64 `mask:"random100000.4"`
	}
	type float64PtrTest struct {
		Usagi *float64 `mask:"random100000.4"`
	}
	type float32SliceTest struct {
		Usagi []float32 `mask:"random100000.4"`
	}
	type float64SliceTest struct {
		Usagi []float64 `mask:"random100000.4"`
	}
	type float32ArrayTest struct {
		Usagi [3]float32 `mask:"random100000.4"`
	}
	type float64ArrayTest struct {
		Usagi [3]float64 `mask:"random100000.4"`
	}
	type float64SlicePtrTest struct {
		Usagi *[]float64 `mask:"random100000.4"`
	}
	type stringToIntTest struct {
		Usagi map[string]int `mask:"random1000"`
	}
	type stringToInt32Test struct {
		Usagi map[string]int32 `mask:"random1000"`
	}
	type stringToInt64Test struct {
		Usagi map[string]int64 `mask:"random1000"`
	}

	tests := map[string]struct {
		input any
		want  any
	}{
		"int fields": {
			input: &intTest{Usagi: 20190122},
			want:  &intTest{Usagi: 829},
		},
		"int16 fields": {
			input: &int16Test{Usagi: 2019},
			want:  &int16Test{Usagi: 829},
		},
		"int32 fields": {
			input: &int32Test{Usagi: 20190122},
			want:  &int32Test{Usagi: 829},
		},
		"int64 fields": {
			input: &int64Test{Usagi: 20190122},
			want:  &int64Test{Usagi: 829},
		},
		"zero int fields": {
			input: &intTest{},
			want:  &intTest{Usagi: 0},
		},
		"int ptr fields": {
			input: &intPtrTest{Usagi: convertIntPtr(20190122)},
			want:  &intPtrTest{Usagi: convertIntPtr(829)},
		},
		"nil int ptr fields": {
			input: &intPtrTest{},
			want:  &intPtrTest{Usagi: nil},
		},
		"int slice fields": {
			input: &intSliceTest{Usagi: []int{20190122, 20200501, 20200501}},
			want:  &intSliceTest{Usagi: []int{829, 830, 400}},
		},
		"int32 slice fields": {
			input: &int32SliceTest{Usagi: []int32{20190122, 20200501, 20200501}},
			want:  &int32SliceTest{Usagi: []int32{829, 830, 400}},
		},
		"int64 slice fields": {
			input: &int64SliceTest{Usagi: []int64{20190122, 20200501, 20200501}},
			want:  &int64SliceTest{Usagi: []int64{829, 830, 400}},
		},
		"int array fields": {
			input: &intArrayTest{Usagi: [2]int{20190122, 20200501}},
			want:  &intArrayTest{Usagi: [2]int{829, 830}},
		},
		"int32 array fields": {
			input: &int32ArrayTest{Usagi: [2]int32{20190122, 20200501}},
			want:  &int32ArrayTest{Usagi: [2]int32{829, 830}},
		},
		"int64 array fields": {
			input: &int64ArrayTest{Usagi: [2]int64{20190122, 20200501}},
			want:  &int64ArrayTest{Usagi: [2]int64{829, 830}},
		},
		"nil int slice fields": {
			input: &intSliceTest{},
			want:  &intSliceTest{Usagi: ([]int)(nil)},
		},
		"int slice ptr fields": {
			input: &intSlicePtrTest{Usagi: convertIntSlicePtr([]int{20190122, 20200501, 20200501})},
			want:  &intSlicePtrTest{Usagi: convertIntSlicePtr([]int{829, 830, 400})},
		},
		"nil int slice ptr fields": {
			input: &intSlicePtrTest{},
			want:  &intSlicePtrTest{Usagi: (*[]int)(nil)},
		},
		"float32 fields": {
			input: &float32Test{Usagi: 20190122},
			want:  &float32Test{Usagi: 96011.8989},
		},
		"float64 fields": {
			input: &float64Test{Usagi: 20190122},
			want:  &float64Test{Usagi: 96011.8989},
		},
		"zero float64 fields": {
			input: &float64Test{},
			want:  &float64Test{Usagi: 0},
		},
		"float64 ptr fields": {
			input: &float64PtrTest{Usagi: convertFloat64Ptr(20190122)},
			want:  &float64PtrTest{Usagi: convertFloat64Ptr(96011.8989)},
		},
		"nil float64 ptr fields": {
			input: &float64PtrTest{},
			want:  &float64PtrTest{Usagi: nil},
		},
		"float32 slice fields": {
			input: &float32SliceTest{Usagi: []float32{20190122, 20200501, 20200501}},
			want:  &float32SliceTest{Usagi: []float32{96011.8989, 90863.3149, 32310.0201}},
		},
		"float64 slice fields": {
			input: &float64SliceTest{Usagi: []float64{20190122, 20200501, 20200501}},
			want:  &float64SliceTest{Usagi: []float64{96011.8989, 90863.3149, 32310.0201}},
		},
		"float32 array fields": {
			input: &float32ArrayTest{Usagi: [3]float32{20190122, 20200501, 20200501}},
			want:  &float32ArrayTest{Usagi: [3]float32{96011.8989, 90863.3149, 32310.0201}},
		},
		"float64 array fields": {
			input: &float64ArrayTest{Usagi: [3]float64{20190122, 20200501, 20200501}},
			want:  &float64ArrayTest{Usagi: [3]float64{96011.8989, 90863.3149, 32310.0201}},
		},
		"nil float64 slice fields": {
			input: &float64SliceTest{},
			want:  &float64SliceTest{Usagi: ([]float64)(nil)},
		},
		"float64 slice ptr fields": {
			input: &float64SlicePtrTest{Usagi: convertFloat64SlicePtr([]float64{20190122, 20200501, 20200501})},
			want:  &float64SlicePtrTest{Usagi: convertFloat64SlicePtr([]float64{96011.8989, 90863.3149, 32310.0201})},
		},
		"nil float64 slice ptr fields": {
			input: &float64SlicePtrTest{},
			want:  &float64SlicePtrTest{Usagi: (*[]float64)(nil)},
		},
		"string to int map fields": {
			input: &stringToIntTest{Usagi: map[string]int{"うさぎ": 20190122}},
			want:  &stringToIntTest{Usagi: map[string]int{"うさぎ": 829}},
		},
		"string to int32 map fields": {
			input: &stringToInt32Test{Usagi: map[string]int32{"うさぎ": 20190122}},
			want:  &stringToInt32Test{Usagi: map[string]int32{"うさぎ": 829}},
		},
		"string to int64 map fields": {
			input: &stringToInt64Test{Usagi: map[string]int64{"うさぎ": 20190122}},
			want:  &stringToInt64Test{Usagi: map[string]int64{"うさぎ": 829}},
		},
	}

	for name, tt := range tests {
		t.Run(defaultTestCase(name), func(t *testing.T) {
			defer cleanup(t)
			rand.Seed(rand.NewSource(1).Int63())
			got, err := Mask(tt.input)
			if assert.NoError(t, err) {
				if diff := cmp.Diff(tt.want, got, cmpopts.SortMaps(func(i, j string) bool { return i < j })); diff != "" {
					t.Error(diff)
				}
			}
		})

		t.Run(newMaskerTestCase(name), func(t *testing.T) {
			rand.Seed(rand.NewSource(1).Int63())
			m := newMasker()
			got, err := m.Mask(tt.input)
			assert.Nil(t, err)
			if diff := cmp.Diff(tt.want, got, allowUnexported(tt.input)); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestMaskZero(t *testing.T) {
	type stringTest struct {
		Usagi string `mask:"zero"`
	}
	type stringPtrTest struct {
		Usagi *string `mask:"zero"`
	}
	type stringSliceTest struct {
		Usagi []string `mask:"zero"`
	}
	type stringArrayTest struct {
		Usagi [3]string `mask:"zero"`
	}
	type stringSlicePtrTest struct {
		Usagi *[]string `mask:"zero"`
	}
	type intTest struct {
		Usagi int `mask:"zero"`
	}
	type uintTest struct {
		Usagi uint `mask:"zero"`
	}
	type float64Test struct {
		Usagi float64 `mask:"zero"`
	}
	type boolTest struct {
		Usagi bool `mask:"zero"`
	}
	type mapStringToStringTest struct {
		Usagi map[string]string `mask:"zero"`
	}
	type structTest struct {
		StringTest stringTest `mask:"zero"`
	}

	tests := map[string]struct {
		input any
		want  any
	}{
		"string fields": {
			input: &stringTest{Usagi: "ヤハッ！"},
			want:  &stringTest{Usagi: ""},
		},
		"string empty fields": {
			input: &stringTest{},
			want:  &stringTest{Usagi: ""},
		},
		"string ptr fields": {
			input: &stringPtrTest{Usagi: convertStringPtr("ヤハッ！")},
			want:  &stringPtrTest{},
		},
		"nil string ptr fields": {
			input: &stringPtrTest{},
			want:  &stringPtrTest{Usagi: nil},
		},
		"string slice fields": {
			input: &stringSliceTest{Usagi: []string{"ハァ？", "ウラ", "フゥン"}},
			want:  &stringSliceTest{},
		},
		"string array fields": {
			input: &stringArrayTest{Usagi: [3]string{"ハァ？", "ウラ", "フゥン"}},
			want:  &stringArrayTest{Usagi: [3]string{}},
		},
		"nil string slice fields": {
			input: &stringSliceTest{},
			want:  &stringSliceTest{Usagi: ([]string)(nil)},
		},
		"string slice ptr fields": {
			input: &stringSlicePtrTest{Usagi: convertStringSlicePtr([]string{"ハァ？", "ウラ", "フゥン"})},
			want:  &stringSlicePtrTest{},
		},
		"nil string slice ptr fields": {
			input: &stringSlicePtrTest{},
			want:  &stringSlicePtrTest{Usagi: (*[]string)(nil)},
		},
		"int fields": {
			input: &intTest{Usagi: 20190122},
			want:  &intTest{Usagi: 0},
		},
		"zero int fields": {
			input: &intTest{},
			want:  &intTest{Usagi: 0},
		},
		"uint fields": {
			input: &uintTest{Usagi: 20190122},
			want:  &uintTest{Usagi: 0},
		},
		"zero uint fields": {
			input: &uintTest{},
			want:  &uintTest{Usagi: 0},
		},
		"float64 fields": {
			input: &float64Test{Usagi: 20190122},
			want:  &float64Test{Usagi: 0},
		},
		"zero float64 fields": {
			input: &float64Test{},
			want:  &float64Test{Usagi: 0},
		},
		"bool fields": {
			input: &boolTest{Usagi: true},
			want:  &boolTest{Usagi: false},
		},
		"zero bool fields": {
			input: &boolTest{},
			want:  &boolTest{},
		},
		"map string to string fields": {
			input: &mapStringToStringTest{Usagi: map[string]string{"うさぎ": "ハァ？", "うさぎ2": "ウラ", "うさぎ3": "フゥン"}},
			want:  &mapStringToStringTest{},
		},
		"nil map string to string fields": {
			input: &mapStringToStringTest{},
			want:  &mapStringToStringTest{},
		},
		"struct fields": {
			input: &structTest{
				StringTest: stringTest{Usagi: "ヤハッ！"},
			},
			want: &structTest{
				StringTest: stringTest{},
			},
		},
	}

	for name, tt := range tests {
		t.Run(defaultTestCase(name), func(t *testing.T) {
			defer cleanup(t)
			rand.Seed(rand.NewSource(1).Int63())
			got, err := Mask(tt.input)
			assert.Nil(t, err)
			if diff := cmp.Diff(tt.want, got, allowUnexported(tt.input)); diff != "" {
				t.Error(diff)
			}
		})

		t.Run(newMaskerTestCase(name), func(t *testing.T) {
			rand.Seed(rand.NewSource(1).Int63())
			m := newMasker()
			got, err := m.Mask(tt.input)
			assert.Nil(t, err)
			if diff := cmp.Diff(tt.want, got, allowUnexported(tt.input)); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestMaskFieldName(t *testing.T) {
	tests := map[string]struct {
		before func(m *Masker)
		input  any
		want   any
	}{
		"matches field of string type": {
			before: func(m *Masker) {
				m.RegisterMaskField("S", MaskTypeFilled+"4")
			},
			input: struct {
				S string
			}{
				S: "Hello World",
			},
			want: struct {
				S string
			}{
				S: "****",
			},
		},
		"matches field of int type": {
			before: func(m *Masker) {
				rand.Seed(rand.NewSource(1).Int63())
				m.RegisterMaskField("I", MaskTypeRandom+"100")
			},
			input: struct {
				I int
			}{
				I: 10,
			},
			want: struct {
				I int
			}{
				I: 29,
			},
		},
		"matches field of float64 type": {
			before: func(m *Masker) {
				rand.Seed(rand.NewSource(1).Int63())
				m.RegisterMaskField("F", MaskTypeRandom+"100.3")
			},
			input: struct {
				F float64
			}{
				F: 10,
			},
			want: struct {
				F float64
			}{
				F: 96.011,
			},
		},
		"matches field of any type": {
			before: func(m *Masker) {
				rand.Seed(rand.NewSource(1).Int63())
				m.RegisterMaskField("S", MaskTypeZero)
				m.RegisterMaskField("I", MaskTypeZero)
				m.RegisterMaskField("F", MaskTypeZero)
				m.RegisterMaskField("A", MaskTypeZero)
			},
			input: struct {
				S string
				I int
				F float64
				A any
			}{
				S: "String",
				I: 10,
				F: 10.2,
				A: struct{ SS string }{"Child"},
			},
			want: struct {
				S string
				I int
				F float64
				A any
			}{
				S: "",
				I: 0,
				F: 0,
				A: struct{ SS string }{},
			},
		},
		"matches field of map type": {
			before: func(m *Masker) {
				rand.Seed(rand.NewSource(1).Int63())
				m.RegisterMaskField("S", MaskTypeFilled+"4")
				m.RegisterMaskField("I", MaskTypeRandom+"100")
			},
			input: map[string]any{
				"M": map[string]any{
					"S": "Hello world",
					"I": 10,
				},
			},
			want: map[string]any{
				"M": map[string]any{
					"S": "****",
					"I": 29,
				},
			},
		},
		"matches field of json": {
			before: func(m *Masker) {
				rand.Seed(rand.NewSource(1).Int63())
				m.RegisterMaskField("S", MaskTypeFilled+"4")
				m.RegisterMaskField("F", MaskTypeRandom+"100")
			},
			input: unmarshalJson(t, `{"M":{"S": "Hello World","F":10}}`),
			want: map[string]any{
				"M": map[string]any{
					"S": "****",
					"F": 96.0,
				},
			},
		},
		"no type match with registered tags": {
			before: func(m *Masker) {
				m.RegisterMaskField("S", MaskTypeRandom+"100") // random tag is only int or float
			},
			input: struct {
				S string
			}{
				S: "Hello World",
			},
			want: struct {
				S string
			}{
				S: "Hello World",
			},
		},
	}
	for name, tt := range tests {
		t.Run(defaultTestCase(name), func(t *testing.T) {
			defer cleanup(t)
			tt.before(defaultMasker)
			got, err := Mask(tt.input)
			assert.Nil(t, err)
			if diff := cmp.Diff(tt.want, got, allowUnexported(tt.input)); diff != "" {
				t.Error(diff)
			}
		})

		t.Run(newMaskerTestCase(name), func(t *testing.T) {
			m := newMasker()
			tt.before(m)
			got, err := m.Mask(tt.input)
			assert.Nil(t, err)
			if diff := cmp.Diff(tt.want, got, allowUnexported(tt.input)); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestSetTagName(t *testing.T) {
	t.Run("change a tag name", func(t *testing.T) {
		m := newMasker()
		m.SetTagName("fake")

		input := struct {
			SM string `mask:"filled4"`
			SF string `fake:"filled4"`
		}{
			SM: "Hello World",
			SF: "Hello World",
		}
		want := struct {
			SM string `mask:"filled4"`
			SF string `fake:"filled4"`
		}{
			SM: "Hello World",
			SF: "****",
		}
		got, _ := m.Mask(input)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Error(diff)
		}
	})
	t.Run("change a empty tag name", func(t *testing.T) {
		m := newMasker()
		m.SetTagName("")

		input := struct {
			SM string `mask:"filled4"`
			SF string `fake:"filled4"`
		}{
			SM: "Hello World",
			SF: "Hello World",
		}
		want := struct {
			SM string `mask:"filled4"`
			SF string `fake:"filled4"`
		}{
			SM: "****",
			SF: "Hello World",
		}
		got, _ := m.Mask(input)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Error(diff)
		}
	})
}

func TestSetMaskChar(t *testing.T) {
	t.Run("change a mask character", func(t *testing.T) {
		defer cleanup(t)
		SetMaskChar("-")

		input := struct {
			S string `mask:"filled4"`
		}{
			S: "Hello World",
		}
		want := struct {
			S string `mask:"filled4"`
		}{
			S: "----",
		}
		got, _ := Mask(input)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Error(diff)
		}
	})
	t.Run("change a empty mask character", func(t *testing.T) {
		defer cleanup(t)
		SetMaskChar("")

		input := struct {
			S string `mask:"filled4"`
		}{
			S: "Hello World",
		}
		want := struct {
			S string `mask:"filled4"`
		}{
			S: "",
		}
		got, _ := Mask(input)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Error(diff)
		}
	})
}

func allowUnexported(v any) cmp.Options {
	var options cmp.Options
	if !reflect.ValueOf(v).IsValid() {
		return options
	}
	rt, ok := getStructType(reflect.TypeOf(v))
	if !ok {
		return options
	}

	rv := reflect.New(rt).Elem()
	options = append(options, cmp.AllowUnexported(rv.Interface()))
	for i := 0; i < rv.NumField(); i++ {
		if rt2, ok := getStructType(rv.Field(i).Type()); ok {
			rv2 := reflect.New(rt2).Elem()
			options = append(options, allowUnexported(rv2.Interface())...)
		}
	}

	return options
}

func getStructType(rt reflect.Type) (reflect.Type, bool) {
	switch rt.Kind() {
	case reflect.Ptr, reflect.Slice:
		return getStructType(rt.Elem())
	case reflect.Struct:
		return rt, true
	default:
		return rt, false
	}
}

func convertStringPtr(s string) *string {
	return &s
}
func convertStringSlicePtr(s []string) *[]string {
	return &s
}
func convertInt16Ptr(i int16) *int16 {
	return &i
}
func convertInt32Ptr(i int32) *int32 {
	return &i
}
func convertInt64Ptr(i int64) *int64 {
	return &i
}
func convertIntPtr(i int) *int {
	return &i
}
func convertIntSlicePtr(i []int) *[]int {
	return &i
}
func convertFloat32Ptr(f float32) *float32 {
	return &f
}
func convertFloat64Ptr(f float64) *float64 {
	return &f
}
func convertFloat64SlicePtr(f []float64) *[]float64 {
	return &f
}
func convertBoolPtr(v bool) *bool {
	return &v
}
func convertAnyPtr(v any) *any {
	return &v
}

func unmarshalJson(t *testing.T, s string) any {
	t.Helper()
	var ret any
	if err := json.Unmarshal([]byte(s), &ret); err != nil {
		t.Fatal(err)
	}

	return ret
}

func defaultTestCase(name string) string {
	return "default Masker:" + name
}
func newMaskerTestCase(name string) string {
	return "newMasker:" + name
}

func cleanup(t *testing.T) {
	t.Helper()
	defaultMasker.typeToStructCache = make(map[string]structType)
	SetMaskChar(maskChar)
}

func newMasker() *Masker {
	m := NewMasker()
	m.RegisterMaskStringFunc(MaskTypeFilled, m.MaskFilledString)
	m.RegisterMaskStringFunc(MaskTypeFixed, m.MaskFixedString)
	m.RegisterMaskStringFunc(MaskTypeHash, m.MaskHashString)
	m.RegisterMaskIntFunc(MaskTypeRandom, m.MaskRandomInt)
	m.RegisterMaskFloat64Func(MaskTypeRandom, m.MaskRandomFloat64)
	m.RegisterMaskAnyFunc(MaskTypeZero, m.MaskZero)
	return m
}
