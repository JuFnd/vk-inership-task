package communication

type (
	SigninRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	SignupRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	AddActorRequest struct {
		Name      string `json:"name"`
		Gender    string `json:"gender"`
		BirthDate string `json:"birth_date"`
	}

	EditActorRequest struct {
		Id        int64   `json:"id"`
		Name      string  `json:"name"`
		Gender    string  `json:"gender"`
		BirthDate string  `json:"birth_date"`
		Films     []int64 `json:"films"`
	}

	DeleteActorRequest struct {
		Id int64 `json:"id"`
	}

	AddFilmRequest struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Rating      float64 `json:"rating"`
		ReleaseDate string  `json:"release_date"`
		Crew        []int64 `json:"crew"`
	}

	EditFilmRequest struct {
		Id          int64   `json:"id"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Rating      float64 `json:"rating"`
		ReleaseDate string  `json:"release_date"`
		Crew        []int64 `json:"crew"`
	}

	DeleteFilmRequest struct {
		Id int64 `json:"id"`
	}
)
