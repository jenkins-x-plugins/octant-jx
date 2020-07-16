package pluginctx

// Context the current context the plugin is running inside
type Context struct {
	// Namespace is the current kubernetes namespace
	Namespace string

	// Composite indicates if this page is inside an Overview or composite view and so we should exclude the breadcrumb/heading
	Composite bool
}
