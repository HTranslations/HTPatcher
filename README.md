# HTPatcher

HTPatcher is an application for applying translations to RPG Maker games. It works with translation patches downloaded from [HTranslations](https://htranslations.com).

## Usage

1. Download a translation patch from HTranslations
2. Launch HTPatcher
3. Select your game executable
4. Select the translation patch file
5. Click "Apply Patch" to apply the translation to your game

## Development

This application is built using [Wails](https://wails.io) with a Svelte frontend and Go backend.

### Prerequisites

- Go
- Bun
- Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### Building from Source

1. Clone the repository:
```bash
git clone https://github.com/HTranslations/HTPatcher
cd HTPatcher
```

2. Install frontend dependencies:
```bash
cd frontend
bun install
cd ..
```

3. Build the application:
```bash
wails build
```

The compiled executable will be in the `build/bin` directory.

### Development Mode

To run in live development mode with hot reload:

```bash
wails dev
```
