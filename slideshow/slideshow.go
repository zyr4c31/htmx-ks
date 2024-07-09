package slideshow

type slide struct {
	number int
	name   string
}

type Slideshow struct {
	Slides []slide
}

func (s *Slideshow) Add(name string) {
	slide := slide{
		number: len(s.Slides),
		name:   name,
	}
	s.Slides = append(s.Slides, slide)
}
