package jsondoc

import (
	"fmt"
	"time"
)

func ExampleMarshalIndent_comments() {
	var strct = struct {
		A string `c:"注释A"`
		B string `c:"注释B"`
		C string
		D string `c:"<注释D>"`
	}{
		A: "1",
		B: "2",
		C: "3",
		D: "4",
	}
	b, err := MarshalIndent(strct, false, ``, `  `)
	fmt.Println(string(b), err)

	// Output:
	// {
	//   "A": "1",	 # 注释A
	//   "B": "2",	 # 注释B
	//   "C": "3",
	//   "D": "4"	 # <注释D>
	// } <nil>
}

func ExampleMarshalIndent_empty_slice() {
	type node struct {
		Name     string `c:"名称"`
		Children []node `c:"孩子"`
	}
	b, err := MarshalIndent(node{}, false, ``, `  `)
	fmt.Println(string(b), err)

	// Output:
	// {
	//   "Name": "",	 # 名称
	//   "Children": [	 # 孩子
	//     {
	//       "Name": "",	 # 名称
	//       "Children": null	 # 孩子
	//     }
	//   ]
	// } <nil>
}
func ExampleMarshalIndent_empty_map() {
	type node struct {
		Name     string       `c:"名称"`
		Children map[int]node `c:"孩子"`
	}
	b, err := MarshalIndent(node{}, false, ``, `  `)
	fmt.Println(string(b), err)

	// Output:
	// {
	//   "Name": "",	 # 名称
	//   "Children": {	 # 孩子
	//     "0": {
	//       "Name": "",	 # 名称
	//       "Children": null	 # 孩子
	//     }
	//   }
	// } <nil>
}

func ExampleMarshalIndent_nil_pointer() {
	type node struct {
		Name string `c:"名称"`
		Next *node  `c:"后一个"`
		Time *time.Time
	}
	b, err := MarshalIndent(node{}, false, ``, `  `)
	fmt.Println(string(b), err)

	// Output:
	// {
	//   "Name": "",	 # 名称
	//   "Next": {	 # 后一个
	//     "Name": "",	 # 名称
	//     "Next": null,	 # 后一个
	//     "Time": "0001-01-01T00:00:00Z"
	//   },
	//   "Time": "0001-01-01T00:00:00Z"
	// } <nil>
}

func ExampleMarshalIndent_nil_pointer_to_anonymous_field() {
	type node struct {
		Name string `c:"名称"`
		Next *node  `c:"后一个"`
		Time *time.Time
	}
	b, err := MarshalIndent(struct{ *node }{}, false, ``, `  `)
	fmt.Println(string(b), err)

	// Output:
	// {
	//   "Name": "",	 # 名称
	//   "Next": null,	 # 后一个
	//   "Time": "0001-01-01T00:00:00Z"
	// } <nil>
}
