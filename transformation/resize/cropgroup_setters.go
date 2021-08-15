package resize

func (c *CropGroup) X(x int) *CropGroup {
	c.x = x

	return c
}

func (c *CropGroup) XPercent(x float32) *CropGroup {
	c.x = x

	return c
}

func (c *CropGroup) XExpr(x string) *CropGroup {
	c.x = x

	return c
}

func (c *CropGroup) Y(y int) *CropGroup {
	c.y = y

	return c
}

func (c *CropGroup) YPercent(y float32) *CropGroup {
	c.y = y

	return c
}

func (c *CropGroup) YExpr(y string) *CropGroup {
	c.y = y

	return c
}

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
