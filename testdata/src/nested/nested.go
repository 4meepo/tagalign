package nested

type S struct {
	A struct {
		B string `env:"" required:""`
		C struct {
			D int `default:"" env:""`
		} `embed:"" prefix:"" envprefix:""` // want `tag is not aligned , should be: embed:"" envprefix:"" prefix:""`
	} `embed:"" envprefix:"" prefix:""`
}
