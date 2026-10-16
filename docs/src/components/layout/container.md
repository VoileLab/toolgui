# Container

Container is the most basic layout component.
Its definition is in tgframe.
The Main "container" and Sidebar "container" are Containers.

## Usage

A page does not create a bare container. Every layout component hands one out:
[Box](box.md) for a plain one, and [Column](column.md), [Tab](tab.md),
[Expand](expand.md) and the rest for theirs.

```go
box := tgcomp.Box(p.Main)
tgcomp.Text(box, "inside the box")
```

A [custom component](../../architecture/custom-components.md) that owns
containers creates them with `Container.AddContainerTo`, so they sit under the
component in the node tree.
