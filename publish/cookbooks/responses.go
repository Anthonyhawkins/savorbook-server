package cookbooks

type CookbookResponse struct {
	ID       uint              `json:"id"`
	Title    string            `json:"title"`
	SubTitle string            `json:"subTitle"`
	Blurb    string            `json:"blurb"`
	Image    string            `json:"image"`
	Sections []SectionResponse `json:"sections"`
}

type SectionResponse struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Overview string  `json:"overview"`
	Recipes  []int64 `json:"recipes"`
}

func (r *CookbookResponse) SerializeCookbook(model *CookbookModel) {
	r.ID = model.ID
	r.Title = model.Title
	r.SubTitle = model.SubTitle
	r.Blurb = model.Blurb
	r.Image = model.Image
	r.SerializeSections(model.Sections)
}

func (r *CookbookResponse) SerializeSections(sectionModels []SectionModel) {
	sections := make([]SectionResponse, 0)
	for _, sectionModel := range sectionModels {
		var section SectionResponse
		section.ID = sectionModel.ID
		section.Name = sectionModel.Name
		section.Overview = sectionModel.Overview
		if sectionModel.Recipes == nil {
			section.Recipes = make([]int64, 0)
		} else {
			section.Recipes = sectionModel.Recipes
		}
		sections = append(sections, section)
	}
	r.Sections = sections
}
