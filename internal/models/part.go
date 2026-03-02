package models

type Part struct {
	ID           string   `quire:"ID" json:"-"`
	Name         string   `quire:"Nombre" json:"Name"`
	Brand        string   `quire:"Marca" json:"Brand"`
	Category     string   `quire:"Categoria" json:"Category"`
	Subcategoria string   `quire:"Subcategoria" json:"Subcategoria"`
	Model        string   `quire:"Modelo" json:"Model"`
	Year         string   `quire:"Año" json:"Year"`
	Price        *float64 `quire:"Precio" json:"Price,omitempty"`
	Description  string   `quire:"Descripcion" json:"Description"`
	Estado       string   `quire:"Estado" json:"Estado"`
	Imagenes     string   `quire:"Imagenes" json:"-"`
	ImagenesArr  []string `quire:"-" json:"ImagenesArr"`
}

type SearchResult struct {
	Part  Part
	Score int
}

type FilterOptions struct {
	Brand        string
	Category     string
	Subcategoria string
}

type User struct {
	Username string `quire:"Usuario"`
	Password string `quire:"Password"`
}

type Session struct {
	Username string
}
