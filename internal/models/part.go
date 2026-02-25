package models

type Part struct {
	ID          string   `quire:"ID" json:"-"`
	Name        string   `quire:"Nombre" json:"Name"`
	Brand       string   `quire:"Marca" json:"Brand"`
	Type        string   `quire:"Tipo" json:"Type"`
	Model       string   `quire:"Modelo" json:"Model"`
	Year        string   `quire:"Año" json:"Year"`
	Price       *float64 `quire:"Precio" json:"Price,omitempty"`
	Description string   `quire:"Descripcion" json:"Description"`
	Estado      string   `quire:"Estado" json:"Estado"`
	Imagenes    string   `quire:"Imagenes" json:"-"`
	ImagenesArr []string `quire:"-" json:"ImagenesArr"`
}

type SearchResult struct {
	Part  Part
	Score int
}

type FilterOptions struct {
	Brand string
	Type  string
}

type User struct {
	Username string `quire:"Usuario"`
	Password string `quire:"Password"`
}

type Session struct {
	Username string
}
