// Package formats implements graph exporters and importers in different formats.
// Exporters can implement either GraphExporter or LayoutExporter interface or both.
// Also high-level helper functions are available for most common use-cases.
//
// Example:
//
//   g := graph.New()
//   ...
//   formats.ToD3JSON(g, "network.json)
//
//   j := formats.NewD3JSON(fd, true)
//   j.ExportGraph(g)
//
//   w := formats.NewGraphBinary("data/")
//   w.ExportGraph(g)
package formats
