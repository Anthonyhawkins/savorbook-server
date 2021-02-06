package recipes

type RecipeResponse struct {
	ID               uint                      `json:"id"`
	Name             string                    `json:"name"`
	Image            string                    `json:"image"`
	Description      string                    `json:"description"`
	PrepTime         string                    `json:"prepTime"`
	Servings         string                    `json:"servings"`
	Tags             []string                  `json:"tags"`
	ParentRecipes    []ParentRecipeResponse    `json:"parentRecipes"`
	DependentRecipes []DependentRecipeResponse `json:"dependentRecipes"`
	IngredientGroups []IngredientGroupResponse `json:"ingredientGroups"`
	Steps            []StepResponse            `json:"steps"`
}

type DependentRecipeResponse struct {
	ID              uint   `json:"id"`
	DependentRecipe uint   `json:"dependentRecipe"`
	RecipeName      string `json:"name,omitempty"`
	Qty             string `json:"qty"`
}

type ParentRecipeResponse struct {
	ID         uint   `json:"id"`
	RecipeName string `json:"name"`
}

type IngredientGroupResponse struct {
	ID          uint                 `json:"id"`
	GroupName   string               `json:"groupName"`
	Ingredients []IngredientResponse `json:"ingredients"`
}

type IngredientResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Qty  string `json:"qty"`
	Unit string `json:"unit"`
}

type StepResponse struct {
	ID         uint                `json:"id"`
	Type       string              `json:"type"`
	Text       string              `json:"text"`
	StepImages []StepImageResponse `json:"images"`
}

type StepImageResponse struct {
	ID    uint   `json:"id"`
	Image string `json:"src"`
	Text  string `json:"text"`
}

func (r *RecipeResponse) SerializeRecipe(model *RecipeModel) {
	r.ID = model.ID
	r.Name = model.Name
	r.Description = model.Description
	r.PrepTime = model.PrepTime
	r.Servings = model.Servings
	r.Tags = SerializeTags(model.Tags)
	r.serializeDependentRecipes(model.DependentRecipes)
	r.Image = model.Image
	r.serializeSteps(model.Steps)
	r.serializeIngredientGroups(model.IngredientGroups)
	r.ParentRecipes = SerializeParentRecipes(model.ParentRecipes)
}

func SerializeTags(tagModels []TagModel) []string {
	tags := make([]string, 0)
	for _, tagModel := range tagModels {
		tags = append(tags, tagModel.Tag)
	}
	return tags
}

func (r *RecipeResponse) serializeDependentRecipes(recipeDependencies []RecipeDependencyModel) {
	dependencies := make([]DependentRecipeResponse, 0)
	for _, recipeDependency := range recipeDependencies {
		var dependency DependentRecipeResponse
		dependency.ID = recipeDependency.ID
		dependency.DependentRecipe = recipeDependency.DependentRecipe
		dependency.RecipeName = recipeDependency.RecipeName
		dependency.Qty = recipeDependency.Qty
		dependencies = append(dependencies, dependency)
	}
	r.DependentRecipes = dependencies
}

func (r *RecipeResponse) serializeIngredientGroups(ingredientGroupModels []IngredientGroupModel) {
	groups := make([]IngredientGroupResponse, 0)
	for _, groupModel := range ingredientGroupModels {
		var group IngredientGroupResponse
		group.ID = groupModel.ID
		group.GroupName = groupModel.GroupName
		group.serializeIngredients(groupModel.Ingredients)
		groups = append(groups, group)
	}
	r.IngredientGroups = groups
}

func (r *IngredientGroupResponse) serializeIngredients(ingredientModels []IngredientModel) {
	ingredients := make([]IngredientResponse, 0)
	for _, ingredientModel := range ingredientModels {
		var ingredient IngredientResponse
		ingredient.ID = ingredientModel.ID
		ingredient.Name = ingredientModel.Name
		ingredient.Qty = ingredientModel.Qty
		ingredient.Unit = ingredientModel.Unit
		ingredients = append(ingredients, ingredient)
	}
	r.Ingredients = ingredients
}

func (r *RecipeResponse) serializeSteps(stepModels []StepModel) {
	steps := make([]StepResponse, 0)
	for _, stepModel := range stepModels {
		var step StepResponse
		step.ID = stepModel.ID
		step.Type = stepModel.Type
		step.Text = stepModel.Text
		step.serializeStepImages(stepModel.StepImages)
		steps = append(steps, step)
	}
	r.Steps = steps
}

func (r *StepResponse) serializeStepImages(stepImageModels []StepImageModel) {
	stepImages := make([]StepImageResponse, 0)
	for _, stepImageModel := range stepImageModels {
		var stepImage StepImageResponse
		stepImage.ID = stepImageModel.ID
		stepImage.Image = stepImageModel.Image
		stepImage.Text = stepImageModel.Text
		stepImages = append(stepImages, stepImage)
	}
	r.StepImages = stepImages
}

func SerializeParentRecipes(recipeParents []RecipeDependencyModel) []ParentRecipeResponse {
	parents := make([]ParentRecipeResponse, 0)
	for _, recipeParent := range recipeParents {
		var parent ParentRecipeResponse
		parent.ID = recipeParent.RecipeID
		parent.RecipeName = recipeParent.RecipeName
		parents = append(parents, parent)
	}
	return parents
}
