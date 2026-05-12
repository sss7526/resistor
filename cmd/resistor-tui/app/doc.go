// Package app implements the Bubble Tea application layer for the
// interactive resistor TUI.
//
// Architecture Overview
//
// The TUI is structured around a dynamic view-router pattern.
// Each screen implements the View interface and is fully autonomous.
//
// The AppModel acts as a shell that:
//
//   - Holds the currently active View
//   - Delegates Update and View calls
//   - Handles global quit keybindings
//   - Propagates terminal resize events
//   - Applies global layout styling (header/footer)
//
// No business logic is implemented in this layer. All computation
// is delegated to the core resistor library.
//
// Design Goals
//
//   - Strict separation from core logic
//   - No CLI coupling
//   - No routing switches or enums
//   - Fully dynamic view transitions
//   - Clean resize propagation
//   - Scalable for future views
//
// Each view:
//
//   - Manages its own state
//   - Decides when to transition to another view
//   - Handles its own key events (except global quit)
//   - Implements Resizable if it depends on terminal dimensions
//
// Resize Strategy
//
// The AppModel stores the latest terminal width and height.
// When a WindowSizeMsg is received, the current view is resized
// if it implements Resizable.
//
// Whenever a view transition occurs, the new view is immediately
// resized using the last known dimensions. This guarantees layout
// consistency across view transitions.
//
// This architecture avoids tight coupling, routing switches,
// and repetitive boilerplate in individual views.
//
// Keybinding Policy:
//
//   - q and ctrl+c always quit the application.
//   - ESC returns to the main menu when in a subview.
//   - ESC does nothing in the main menu unless filtering.
//   - Each view is responsible for handling its own keys.
//
// This separation prevents tight coupling between
// the router and individual view behavior.
package app