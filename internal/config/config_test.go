package config

// func TestGetRecursiveTags(t *testing.T) {
// 	tt := []struct {
// 		name     string
// 		config   any
// 		expected []*pathConfig
// 		prefix   string
// 	}{
// 		{
// 			name: "Simple Struct",
// 			config: struct {
// 				a string `ssm:"/a"`
// 				b string `ssm:"/b"`
// 				c string `ssm:"/c"`
// 			}{},
// 			expected: []*pathConfig{
// 				{
// 					name:     "/a",
// 					required: false,
// 				},
// 				{
// 					name:     "/b",
// 					required: false,
// 				},
// 				{
// 					name:     "/c",
// 					required: false,
// 				},
// 			},
// 		},
// 		{
// 			name:   "Simple Struct With Prefix",
// 			prefix: "/p",
// 			config: struct {
// 				a string `ssm:"/a"`
// 				b string `ssm:"/b"`
// 				c string `ssm:"/c"`
// 			}{},
// 			expected: []*pathConfig{
// 				{
// 					name:     "/p/a",
// 					required: false,
// 				},
// 				{
// 					name:     "/p/b",
// 					required: false,
// 				},
// 				{
// 					name:     "/p/c",
// 					required: false,
// 				},
// 			},
// 		},
// 		{
// 			name:   "Simple Struct With Ignored",
// 			prefix: "/p",
// 			config: struct {
// 				a string `ssm:"/a"`
// 				b string `ssm:"/b"`
// 				c string
// 			}{},
// 			expected: []*pathConfig{
// 				{
// 					name:     "/p/a",
// 					required: false,
// 				},
// 				{
// 					name:     "/p/b",
// 					required: false,
// 				},
// 			},
// 		},

// 		{
// 			name:   "Simple Struct With Required And Ignored",
// 			prefix: "/p",
// 			config: struct {
// 				a string `ssm:"/a"`
// 				b string `ssm:"/b,required"`
// 				c string
// 			}{},
// 			expected: []*pathConfig{
// 				{
// 					name:     "/p/a",
// 					required: false,
// 				},
// 				{
// 					name:     "/p/b",
// 					required: true,
// 				},
// 			},
// 		},
// 		{
// 			name: "Simple Struct With Required Attribute",
// 			config: struct {
// 				a string `ssm:"/a,required"`
// 			}{},
// 			expected: []*pathConfig{
// 				{
// 					name:     "/a",
// 					required: true,
// 				},
// 			},
// 		},
// 		{
// 			name: "Struct With Nested Struct",
// 			config: struct {
// 				A string `ssm:"/a"`
// 				B struct {
// 					C string `ssm:"/c"`
// 				} `ssm:"/b"`
// 			}{},
// 			expected: []*pathConfig{
// 				{
// 					name:     "/a",
// 					required: false,
// 				},
// 				{
// 					name:     "/b/c",
// 					required: false,
// 				},
// 			},
// 		},
// 		{
// 			name: "Struct With Nested Struct With Required In Parent",
// 			config: struct {
// 				A string `ssm:"/a,required"`
// 				B struct {
// 					C string `ssm:"/c"`
// 				} `ssm:"/b"`
// 			}{},
// 			expected: []*pathConfig{
// 				{
// 					name:     "/a",
// 					required: true,
// 				},
// 				{
// 					name:     "/b/c",
// 					required: false,
// 				},
// 			},
// 		},
// 		{
// 			name: "Struct With Nested Struct With Required In Child",
// 			config: struct {
// 				A string `ssm:"/a"`
// 				B struct {
// 					C string `ssm:"/c,required"`
// 				} `ssm:"/b"`
// 			}{},
// 			expected: []*pathConfig{
// 				{
// 					name:     "/a",
// 					required: false,
// 				},
// 				{
// 					name:     "/b/c",
// 					required: true,
// 				},
// 			},
// 		},
// 		{
// 			name: "Struct With Nested Struct With Required In Parent and Child",
// 			config: struct {
// 				A string `ssm:"/a,required"`
// 				B struct {
// 					C string `ssm:"/c,required"`
// 				} `ssm:"/b"`
// 			}{},
// 			expected: []*pathConfig{
// 				{
// 					name:     "/a",
// 					required: true,
// 				},
// 				{
// 					name:     "/b/c",
// 					required: true,
// 				},
// 			},
// 		},
// 	}

// 	for _, tc := range tt {
// 		t.Run(tc.name, func(t *testing.T) {
// 			c := getRecursiveTags(reflect.ValueOf(tc.config), tc.prefix)
// 			if len(c) != len(tc.expected) {
// 				t.Fail()
// 			}

// 			for i, e := range c {
// 				fmt.Println(e.name, tc.expected[i].name)
// 				if !cmp.Equal(e, tc.expected[i], cmp.AllowUnexported(pathConfig{})) {
// 					t.Fail()
// 				}
// 			}
// 		})
// 	}
// }
