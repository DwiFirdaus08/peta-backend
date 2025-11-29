package model

// Location dipindahkan ke sini agar bisa di-import
// Field harus diawali Huruf Besar agar bisa diakses (Exported)
type Location struct {
	ID   string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name string  `json:"name" bson:"name"`
	Lat  float64 `json:"lat" bson:"lat"`
	Lng  float64 `json:"lng" bson:"lng"`
	Desc string  `json:"desc" bson:"desc"`
}