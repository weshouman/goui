# GoUI Framework

A reusable Terminal User Interface framework for building command-driven applications in Go.

## Features

- **Onion Architecture**: Clean separation of concerns with domain, service, UI, and app layers
- **Command System**: Extensible command registration with autocomplete and state transitions
- **Multiple Layout Types**: Support for tables, lists, text views, trees, and combinations
- **Mode Management**: Normal, search, and command input modes
- **State Management**: Configurable application states with template rendering
- **Dependency Injection**: Clean service layer with registry pattern

## Architecture

```
pkg/
  app/            # Application scaffolding and coordination
  domain/         # Core types: State, Command, Mode, Config
  service/        # Business logic: Registry, StateManager, SearchService
  ui/             # UI components: ViewContainer, renderer abstractions
  util/           # Utility functions
```

## Usage

Import the framework and create your app:

```go
import (
    gotui "github.com/ourorg/gotui-08-framework/pkg/app"
    "github.com/ourorg/gotui-08-framework/pkg/service"
)

reg := service.NewRegistry()
// Register your states and commands
app := gotui.New(reg, gotui.Options{
    Logo: "Your App Logo",
    DebugUI: false,
})
app.Run()
```

## Key Components

- **Registry**: Manages states, commands, and their relationships
- **State**: Defines application screens and their display layout
- **Command**: Executable actions with handlers and state transitions
- **App**: Main application coordinator with input handling
- **ViewContainer**: Manages UI components for different layout types

## Build

```bash
go build ./pkg/...
```

