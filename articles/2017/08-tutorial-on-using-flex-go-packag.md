---
Id: 9
Title: Tutorial for github.com/kjk/flex Go package (implementation of CSS flexbox algorithm)
Date: 2017-08-04T18:32:51-07:00
Format: Markdown
Tags: go
---

Package [github.com/kjk/flex](https://github.com/kjk/flex) implements [CSS flexbox](https://www.w3.org/TR/css-flexbox-1/) layout algorithm in Go.

It's a pure Go [port](/article/wN9R/experience-porting-4.5k-loc-of-c-to-go-facebooks-css-flexbox-implementation-yoga.html) of Facebook's [Yoga](https://github.com/facebook/yoga) C library.

## High-level API overview

Despite implementing CSS flexbox spec, it isn't tied to CSS/HTML in any way. Yoga, for example, can be integrated with iOS app and used to layout UIView hierarchy.

The library works on abstract tree of nodes. In HTML a node would correspond to a block element like a `div`. When used in Cocoa app, a node could represent `UIView` or `NSView`.

When used on windows, it could represent a HWND-based control.

The high-level use is:
* create a tree of nodes that represents a layout you want to represent
* set desired flexbox properties on each node using `node.StyleSet*()` [functions](https://github.com/kjk/flex/blob/master/yoga_props.go#L58)
* call `flex.CalculateLayout(rootNode, parentWidth, parentHeight, direction)`
* each node is now measured and positioned so you can e.g. size and position widgets associated with each node. You can get the size on position of nodes with `node.LayoutGet*()` [functions](https://github.com/kjk/flex/blob/master/yoga_props.go#L531)
* when layout hierachy changes (e.g. the node represents a label and you've changed its text, which changes it's intrisic size), call `node.MarkDirty()` and `CalculateLayout()` to re-calculate new size/position of the nodes

## An exmple

Let's assume that we want to re-create the following HTML layout:
```html
<div id="percentage_multiple_nested_with_padding_margin_and_percentage_values" style="width: 200px; height: 200px; flex-direction: column;">
  <div style="flex-grow: 1; flex-basis: 10%; min-width: 60%; margin: 5px; padding: 3px;">
    <div style="width: 50%; margin: 5px; padding: 3%;">
      <div style="width: 45%; margin: 5%; padding: 3px;"></div>
    </div>
  </div>
  <div style="flex-grow: 4; flex-basis: 15%; min-width: 20%;"></div>
</div>
```

The equivalent Go code is:
```go
	config := flex.NewConfig()

	root := flex.NewNodeWithConfig(config)
	root.StyleSetWidth(200)
	root.StyleSetHeight(200)

	rootChild0 := flex.NewNodeWithConfig(config)
	rootChild0.StyleSetFlexGrow(1)
	rootChild0.StyleSetFlexBasisPercent(10)
	rootChild0.StyleSetMargin(EdgeLeft, 5)
	rootChild0.StyleSetMargin(EdgeTop, 5)
	rootChild0.StyleSetMargin(EdgeRight, 5)
	rootChild0.StyleSetMargin(EdgeBottom, 5)
	rootChild0.StyleSetPadding(EdgeLeft, 3)
	rootChild0.StyleSetPadding(EdgeTop, 3)
	rootChild0.StyleSetPadding(EdgeRight, 3)
	rootChild0.StyleSetPadding(EdgeBottom, 3)
	rootChild0.StyleSetMinWidthPercent(60)
	root.InsertChild(rootChild0, 0)

	rootChild0Child0 := flex.NewNodeWithConfig(config)
	rootChild0Child0.StyleSetMargin(EdgeLeft, 5)
	rootChild0Child0.StyleSetMargin(EdgeTop, 5)
	rootChild0Child0.StyleSetMargin(EdgeRight, 5)
	rootChild0Child0.StyleSetMargin(EdgeBottom, 5)
	rootChild0Child0.StyleSetPaddingPercent(EdgeLeft, 3)
	rootChild0Child0.StyleSetPaddingPercent(EdgeTop, 3)
	rootChild0Child0.StyleSetPaddingPercent(EdgeRight, 3)
	rootChild0Child0.StyleSetPaddingPercent(EdgeBottom, 3)
	rootChild0Child0.StyleSetWidthPercent(50)
	rootChild0.InsertChild(rootChild0Child0, 0)

	rootChild0Child0Child0 := flex.NewNodeWithConfig(config)
	rootChild0Child0Child0.StyleSetMarginPercent(EdgeLeft, 5)
	rootChild0Child0Child0.StyleSetMarginPercent(EdgeTop, 5)
	rootChild0Child0Child0.StyleSetMarginPercent(EdgeRight, 5)
	rootChild0Child0Child0.StyleSetMarginPercent(EdgeBottom, 5)
	rootChild0Child0Child0.StyleSetPadding(EdgeLeft, 3)
	rootChild0Child0Child0.StyleSetPadding(EdgeTop, 3)
	rootChild0Child0Child0.StyleSetPadding(EdgeRight, 3)
	rootChild0Child0Child0.StyleSetPadding(EdgeBottom, 3)
	rootChild0Child0Child0.StyleSetWidthPercent(45)
	rootChild0Child0.InsertChild(rootChild0Child0Child0, 0)

	rootChild1 := flex.NewNodeWithConfig(config)
	rootChild1.StyleSetFlexGrow(4)
	rootChild1.StyleSetFlexBasisPercent(15)
	rootChild1.StyleSetMinWidthPercent(20)
	root.InsertChild(rootChild1, 1)
	flex.CalculateLayout(root, flex.Undefined, flex.Undefined, DirectionLTR)
```

After `CalculateLayout` we can see the position of each node e.g.:

```go
	fmt.Printf("root left: %f\n", root.LayoutGetLeft()) // 0
	fmt.Printf(("root top: %f\n", root.LayoutGetTop()) // 0
	fmt.Printf("root width: %f\n", root.LayoutGetWidth()) // 200
	fmt.Printf("root height: %f\n", root.LayoutGetHeight()) // 200
```

To see example for every flexbox property, look into [github.com/facebook/yoga/gentest/fixtures](https://github.com/facebook/yoga/tree/master/gentest/fixtures). Their names hint at which properties are being used.

Each file there has corresponding `*_test.go` file in [github.com/kjk/flex](https://github.com/kjk/flex) directory which shows how to express it in Go.

## Size of root's parent

Notice that in this particular example we used `flex.Undefined` as both height and width of the parent container.

Imagine you're using `flex` to implment layout for a dekstop application where each `flex.Node` represents a control inside the window.

Window is the parent of root node.

In response to user resizing the window, you want to pass width/height of the window to `flex.CalculateLayout()`.

When you create the window initially, you might do the reverse: pass `flex.Undefined` as width/height of parent container and then use the size of root node as the size of the window, to size it to its content.

## Measure function

Imagine that a node represents an OS button. The button has some intrisic size dictated by its text.

To represent that size `flex` allows setting a measuring function with `node.SetMeasureFunc(measureFunc MeasureFunc)`. It's definition is:

```go
type MeasureFunc func(node *Node, width float32, widthMode MeasureMode, height float32, heightMode MeasureMode) Size
```

The functions takes a hint width/height which is the size of parent container and returns intrinsic size of node.

This is usefule e.g. when a node represents a paragraph of text. When you know width of the parent container, you can break it into multi-line text.

If measuring function needs some state, you can use `node.Context` to store it.
