package resize

func (c *CropGroup) Width(width int) *CropGroup {
	c.width = width

	return c
}

func (c *CropGroup) WidthPercent(width float32) *CropGroup {
	c.width = width

	return c
}

func (c *CropGroup) WidthExpr(width string) *CropGroup {
	c.width = width

	return c
}

func (c *CropGroup) Height(height int) *CropGroup {
	c.height = height

	return c
}

func (c *CropGroup) HeightPercent(height float32) *CropGroup {
	c.height = height

	return c
}

func (c *CropGroup) HeightExpr(height string) *CropGroup {
	c.height = height

	return c
}

func (c *CropGroup) AspectRatio(aspectRatio int) *CropGroup {
	c.aspectRatio = aspectRatio

	return c
}

func (c *CropGroup) AspectRatioPercent(aspectRatio float32) *CropGroup {
	c.aspectRatio = aspectRatio

	return c
}

func (c *CropGroup) AspectRatioExpr(aspectRatio string) *CropGroup {
	c.aspectRatio = aspectRatio

	return c
}

func (c *CropGroup) TestField(testField int) *CropGroup {
	c.testField = testField

	return c
}

func (c *CropGroup) TestFieldExpr(testField string) *CropGroup {
	c.testField = testField

	return c
}

func (c *CropGroup) TestFieldFloat(testField float32) *CropGroup {
	c.testField = testField

	return c
}

func (c *CropGroup) Test3Field(test3Field string) *CropGroup {
	c.test3Field = test3Field

	return c
}

func (c *CropGroup) Test3FieldPercent(test3Field int) *CropGroup {
	c.test3Field = test3Field

	return c
}

func (c *CropGroup) Test4FieldTest(Test4Field int) *CropGroup {
	c.Test4Field = Test4Field

	return c
}

func (c *CropGroup) Test4FieldExpr(Test4Field string) *CropGroup {
	c.Test4Field = Test4Field

	return c
}
