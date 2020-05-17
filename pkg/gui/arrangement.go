package gui

func (gui *Gui) getViewDimensions() map[string]dimensions {
	width, height := gui.g.Size()

	main := "main"
	secondary := "secondary"
	if gui.State.Panels.LineByLine != nil && gui.State.Panels.LineByLine.SecondaryFocused {
		main, secondary = secondary, main
	}

	mainSectionChildren := []*box{
		{
			viewName: main,
			weight:   1,
		},
	}

	if gui.State.SplitMainPanel {
		mainSectionChildren = append(mainSectionChildren, &box{
			viewName: secondary,
			weight:   1,
		})
	}

	// we originally specified this as a ratio i.e. .20 would correspond to a weight of 1 against 4
	sidePanelWidthRatio := gui.Config.GetUserConfig().GetFloat64("gui.sidePanelWidth")
	// we could make this better by creating ratios like 2:3 rather than always 1:something
	mainSectionWeight := int(1/sidePanelWidthRatio) - 1
	sideSectionWeight := 1

	if gui.State.SplitMainPanel {
		mainSectionWeight = 5 // need to shrink side panel to make way for main panels if side-by-side
	}
	currentViewName := gui.currentViewName()
	if currentViewName == "main" {
		if gui.State.ScreenMode == SCREEN_HALF || gui.State.ScreenMode == SCREEN_FULL {
			sideSectionWeight = 0
		}
	} else {
		if gui.State.ScreenMode == SCREEN_HALF {
			mainSectionWeight = 1
		} else if gui.State.ScreenMode == SCREEN_FULL {
			mainSectionWeight = 0
		}
	}

	sidePanelsDirection := COLUMN
	portraitMode := width <= 84 && height > 50
	if portraitMode {
		sidePanelsDirection = ROW
	}

	root := &box{
		direction: ROW,
		children: []*box{
			{
				direction: sidePanelsDirection,
				weight:    1,
				children: []*box{
					{
						direction:           ROW,
						weight:              sideSectionWeight,
						conditionalChildren: gui.sidePanelChildren,
					},
					{
						conditionalDirection: func(width int, height int) int {
							if width < 160 && height > 30 { // 2 80 character width panels
								return ROW
							} else {
								return COLUMN
							}
						},
						direction: COLUMN,
						weight:    mainSectionWeight,
						children:  mainSectionChildren,
					},
				},
			},
			// TODO: actually handle options here. Currently we're just hard-coding it to be set on the bottom row in our layout function given that we need some custom logic to have it share space with other views on that row.
			{
				viewName: "options",
				size:     1,
			},
		},
	}

	return gui.arrangeViews(root, 0, 0, width, height)
}

func (gui *Gui) sidePanelChildren(width int, height int) []*box {
	currentCyclableViewName := gui.currentCyclableViewName()

	if gui.State.ScreenMode == SCREEN_FULL || gui.State.ScreenMode == SCREEN_HALF {
		fullHeightBox := func(viewName string) *box {
			if viewName == currentCyclableViewName {
				return &box{
					viewName: viewName,
					weight:   1,
				}
			} else {
				return &box{
					viewName: viewName,
					size:     0,
				}
			}
		}

		return []*box{
			fullHeightBox("status"),
			fullHeightBox("files"),
			fullHeightBox("branches"),
			fullHeightBox("commits"),
			fullHeightBox("stash"),
		}
	} else if height >= 28 {
		return []*box{
			{
				viewName: "status",
				size:     3,
			},
			{
				viewName: "files",
				weight:   1,
			},
			{
				viewName: "branches",
				weight:   1,
			},
			{
				viewName: "commits",
				weight:   1,
			},
			{
				viewName: "stash",
				size:     3,
			},
		}
	} else {
		squashedHeight := 1
		if height >= 21 {
			squashedHeight = 3
		}

		squashedSidePanelBox := func(viewName string) *box {
			if viewName == currentCyclableViewName {
				return &box{
					viewName: viewName,
					weight:   1,
				}
			} else {
				return &box{
					viewName: viewName,
					size:     squashedHeight,
				}
			}
		}

		return []*box{
			squashedSidePanelBox("status"),
			squashedSidePanelBox("files"),
			squashedSidePanelBox("branches"),
			squashedSidePanelBox("commits"),
			squashedSidePanelBox("stash"),
		}
	}
}

func (gui *Gui) currentCyclableViewName() string {
	currView := gui.g.CurrentView()
	currentCyclebleView := gui.State.PreviousView
	if currView != nil {
		viewName := currView.Name()
		usePreviousView := true
		for _, view := range cyclableViews {
			if view == viewName {
				currentCyclebleView = viewName
				usePreviousView = false
				break
			}
		}
		if usePreviousView {
			currentCyclebleView = gui.State.PreviousView
		}
	}

	// unfortunate result of the fact that these are separate views, have to map explicitly
	if currentCyclebleView == "commitFiles" {
		return "commits"
	}

	return currentCyclebleView
}