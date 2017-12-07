package main

import (
	"fmt"
	"strings"
)

const testMarkdownString = `
## An exmple

Let's assume that we want to re-create the following HTML layout:
___html
<div id="percentage_multiple_nested_with_padding_margin_and_percentage_values" style="width: 200px; height: 200px; flex-direction: column;">
  <div style="flex-grow: 1; flex-basis: 10%; min-width: 60%; margin: 5px; padding: 3px;">
    <div style="width: 50%; margin: 5px; padding: 3%;">
      <div style="width: 45%; margin: 5%; padding: 3px;"></div>
    </div>
  </div>
  <div style="flex-grow: 4; flex-basis: 15%; min-width: 20%;"></div>
</div>
___

The equivalent Go code is:
___go
	config := flex.NewConfig()

	root := flex.NewNodeWithConfig(config)
	root.StyleSetWidth(200)
	root.StyleSetHeight(200)

	rootChild1 := flex.NewNodeWithConfig(config)
	rootChild1.StyleSetFlexGrow(4)
	rootChild1.StyleSetFlexBasisPercent(15)
	rootChild1.StyleSetMinWidthPercent(20)
	root.InsertChild(rootChild1, 1)
	flex.CalculateLayout(root, flex.Undefined, flex.Undefined, DirectionLTR)
___

After CalculateLayout we can see the position of each node e.g.:

___
	fmt.Printf("root left: %f\n", root.LayoutGetLeft()) // 0
	fmt.Printf"root top: %f\n", root.LayoutGetTop()) // 0
___

To see example for every flexbox property, look into [github.com/facebook/yoga/gentest/fixtures](https://github.com/facebook/yoga/tree/master/gentest/fixtures). Their names hint at which properties are being used.
`

func adhocTestChroma() {
	s := testMarkdownString
	s = strings.Replace(s, "___", "```", -1)
	in := []byte(s)
	debugMarkdownCodeHighligh = true
	markdownCodeHighligh(in)

	fmt.Print("\n\n")
	//htmlFormatter.WriteCSS(os.Stdout, highlightStyle)
}
