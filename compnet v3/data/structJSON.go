package data

type Person struct{
	//digunakan untuk decode dan encoding
	Name string `json:"name"`

	Age int `json:"age"`
}