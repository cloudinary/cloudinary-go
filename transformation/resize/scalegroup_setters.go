package resize

func (s *ScaleGroup) Width(width int) *ScaleGroup {
	s.width = width

	return s
}

func (s *ScaleGroup) WidthPercent(width float32) *ScaleGroup {
	s.width = width

	return s
}

func (s *ScaleGroup) WidthExpr(width string) *ScaleGroup {
	s.width = width

	return s
}

func (s *ScaleGroup) Height(height int) *ScaleGroup {
	s.height = height

	return s
}

func (s *ScaleGroup) HeightPercent(height float32) *ScaleGroup {
	s.height = height

	return s
}

func (s *ScaleGroup) HeightExpr(height string) *ScaleGroup {
	s.height = height

	return s
}

func (s *ScaleGroup) AspectRatio(aspectRatio int) *ScaleGroup {
	s.aspectRatio = aspectRatio

	return s
}

func (s *ScaleGroup) AspectRatioPercent(aspectRatio float32) *ScaleGroup {
	s.aspectRatio = aspectRatio

	return s
}

func (s *ScaleGroup) AspectRatioExpr(aspectRatio string) *ScaleGroup {
	s.aspectRatio = aspectRatio

	return s
}
