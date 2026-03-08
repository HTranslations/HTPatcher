package main

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"htpatcher/internal/domain"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func jsonResult(v any) (*mcp.CallToolResult, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("json marshal error: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func registerTools(s *server.MCPServer, svc *Services) {
	// ── Game Detection ─────────────────────────────────────────────

	s.AddTool(mcp.NewTool("detect_game",
		mcp.WithDescription("Detect game version, title, and paths from a game executable path"),
		mcp.WithString("exe_path", mcp.Required(), mcp.Description("Path to the game executable (e.g. Game.exe)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path, err := req.RequireString("exe_path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		info, err := svc.Game.GetGameInfoFromExePath(path)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(info)
	})

	// ── Collection Management ──────────────────────────────────────

	s.AddTool(mcp.NewTool("list_games",
		mcp.WithDescription("List all games in the collection"),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		games, err := svc.Collection.GetGamesCollection()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(games)
	})

	s.AddTool(mcp.NewTool("add_game",
		mcp.WithDescription("Add a game to the collection"),
		mcp.WithString("exe_path", mcp.Required(), mcp.Description("Path to the game executable")),
		mcp.WithString("rj_code", mcp.Required(), mcp.Description("Store code for the game (e.g. RJ123456 for DLsite, d_123456 for DMM)")),
		mcp.WithString("friendly_name", mcp.Required(), mcp.Description("Display name for the game")),
		mcp.WithString("tags_json", mcp.Description("JSON array of tags, e.g. [\"rpg\",\"fantasy\"]. Defaults to []")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		exePath, err := req.RequireString("exe_path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		storeCode, err := req.RequireString("rj_code")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		friendlyName, err := req.RequireString("friendly_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		var tags []string
		tagsJSON := req.GetString("tags_json", "[]")
		if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("invalid tags_json: %v", err)), nil
		}

		gameDir := filepath.Dir(exePath)
		locatedGame := &domain.LocatedGame{
			GameDir: gameDir,
			ExePath: exePath,
		}

		if err := svc.Collection.AddGameToCollection(locatedGame, storeCode, friendlyName, tags); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(locatedGame)
	})

	s.AddTool(mcp.NewTool("remove_game",
		mcp.WithDescription("Remove a game from the collection by ID"),
		mcp.WithString("game_id", mcp.Required(), mcp.Description("UUID of the game to remove")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := req.RequireString("game_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := svc.Collection.RemoveGameFromCollection(id); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("game removed"), nil
	})

	s.AddTool(mcp.NewTool("update_game",
		mcp.WithDescription("Update game metadata (friendly name, tags)"),
		mcp.WithString("game_id", mcp.Required(), mcp.Description("UUID of the game")),
		mcp.WithString("friendly_name", mcp.Required(), mcp.Description("New display name")),
		mcp.WithString("tags_json", mcp.Description("JSON array of tags. Defaults to []")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := req.RequireString("game_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		name, err := req.RequireString("friendly_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var tags []string
		tagsJSON := req.GetString("tags_json", "[]")
		if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("invalid tags_json: %v", err)), nil
		}
		if err := svc.Collection.UpdateGameMetadata(id, name, tags); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("game updated"), nil
	})

	s.AddTool(mcp.NewTool("set_game_play_status",
		mcp.WithDescription("Set the play status of a game (unplayed, playing, on-hold, finished, given-up)"),
		mcp.WithString("game_id", mcp.Required(), mcp.Description("UUID of the game")),
		mcp.WithString("status", mcp.Required(), mcp.Description("Play status: unplayed, playing, on-hold, finished, or given-up")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := req.RequireString("game_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		status, err := req.RequireString("status")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := svc.Collection.SetGamePlayStatus(id, status); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("play status updated"), nil
	})

	// ── Patch Operations ───────────────────────────────────────────

	s.AddTool(mcp.NewTool("load_patch",
		mcp.WithDescription("Load and inspect a .htpatch file, returning its config, dictionary size, and overrides"),
		mcp.WithString("patch_path", mcp.Required(), mcp.Description("Path to the .htpatch file")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path, err := req.RequireString("patch_path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		patchInfo, err := svc.Patch.LoadPatchInfo(path)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		// Return a summary instead of the full dictionary
		summary := map[string]any{
			"patchPath":      patchInfo.PatchPath,
			"dictionarySize": len(patchInfo.Dictionary),
			"overrides":      patchInfo.Overrides,
			"config":         patchInfo.Config,
		}
		return jsonResult(summary)
	})

	s.AddTool(mcp.NewTool("apply_patch",
		mcp.WithDescription("Apply a .htpatch translation patch to a game. Optionally backs up game data first."),
		mcp.WithString("exe_path", mcp.Required(), mcp.Description("Path to the game executable")),
		mcp.WithString("patch_path", mcp.Required(), mcp.Description("Path to the .htpatch file")),
		mcp.WithBoolean("backup", mcp.Description("Create a backup before patching (default: true)")),
		mcp.WithBoolean("inject_message_hide", mcp.Description("Inject MessageWindowHidden plugin for MV/MZ games (default: false)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		exePath, err := req.RequireString("exe_path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		patchPath, err := req.RequireString("patch_path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		backup := req.GetBool("backup", true)
		injectMessageHide := req.GetBool("inject_message_hide", false)

		// Detect game
		gameInfo, err := svc.Game.GetGameInfoFromExePath(exePath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to detect game: %v", err)), nil
		}

		// Load patch
		patchInfo, err := svc.Patch.LoadPatchInfo(patchPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to load patch: %v", err)), nil
		}

		// Backup if requested
		if backup {
			if err := svc.Backup.BackupGameData(gameInfo, patchInfo, injectMessageHide); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("backup failed: %v", err)), nil
			}
		}

		// Apply patch
		if err := svc.Patch.ApplyPatch(ctx, gameInfo, patchInfo, injectMessageHide, false, "", ""); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("patch failed: %v", err)), nil
		}

		svc.Collection.SetGameTranslatedByExePath(exePath, true)

		return mcp.NewToolResultText(fmt.Sprintf("patch applied successfully to %q (%s)", gameInfo.GameTitle, gameInfo.GameVersion)), nil
	})

	s.AddTool(mcp.NewTool("fetch_available_patches",
		mcp.WithDescription("Fetch all available translation patches from HTranslations"),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		patches, err := svc.Patch.FetchAllPatches()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(patches)
	})

	s.AddTool(mcp.NewTool("download_and_apply_patch",
		mcp.WithDescription("Download a patch from a URL, then apply it to a game"),
		mcp.WithString("exe_path", mcp.Required(), mcp.Description("Path to the game executable")),
		mcp.WithString("url", mcp.Required(), mcp.Description("Direct download URL for the .htpatch file")),
		mcp.WithString("file_name", mcp.Required(), mcp.Description("Filename for the downloaded patch")),
		mcp.WithBoolean("backup", mcp.Description("Create a backup before patching (default: true)")),
		mcp.WithBoolean("inject_message_hide", mcp.Description("Inject MessageWindowHidden plugin for MV/MZ games (default: false)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		exePath, err := req.RequireString("exe_path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		url, err := req.RequireString("url")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		fileName, err := req.RequireString("file_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		backup := req.GetBool("backup", true)
		injectMessageHide := req.GetBool("inject_message_hide", false)

		// Detect game
		gameInfo, err := svc.Game.GetGameInfoFromExePath(exePath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to detect game: %v", err)), nil
		}

		// Download and load patch
		patchInfo, err := svc.Download.DownloadPatch(url, fileName, func(filePath string) (*domain.PatchInfo, error) {
			return svc.Patch.LoadPatchInfo(filePath)
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to download patch: %v", err)), nil
		}

		// Backup if requested
		if backup {
			if err := svc.Backup.BackupGameData(gameInfo, patchInfo, injectMessageHide); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("backup failed: %v", err)), nil
			}
		}

		// Apply patch
		if err := svc.Patch.ApplyPatch(ctx, gameInfo, patchInfo, injectMessageHide, false, "", ""); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("patch failed: %v", err)), nil
		}

		svc.Collection.SetGameTranslatedByExePath(exePath, true)

		return mcp.NewToolResultText(fmt.Sprintf("patch downloaded and applied successfully to %q (%s)", gameInfo.GameTitle, gameInfo.GameVersion)), nil
	})

	// ── Backup & Restore ───────────────────────────────────────────

	s.AddTool(mcp.NewTool("restore_backup",
		mcp.WithDescription("Restore a game from its backup, reverting all patch changes"),
		mcp.WithString("exe_path", mcp.Required(), mcp.Description("Path to the game executable")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		exePath, err := req.RequireString("exe_path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		gameInfo, err := svc.Game.GetGameInfoFromExePath(exePath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to detect game: %v", err)), nil
		}
		if err := svc.Backup.RestoreBackup(gameInfo); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		svc.Collection.SetGameTranslatedByExePath(exePath, false)

		return mcp.NewToolResultText("backup restored successfully"), nil
	})

	s.AddTool(mcp.NewTool("has_backup",
		mcp.WithDescription("Check if a game has a backup that can be restored"),
		mcp.WithString("exe_path", mcp.Required(), mcp.Description("Path to the game executable")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		exePath, err := req.RequireString("exe_path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		gameDir := filepath.Dir(exePath)
		backupPath := filepath.Join(gameDir, ".backup")
		_, statErr := os.Stat(backupPath)
		return jsonResult(map[string]bool{"hasBackup": statErr == nil})
	})

	// ── Export ──────────────────────────────────────────────────────

	s.AddTool(mcp.NewTool("export_patched_files",
		mcp.WithDescription("Export all patched files from a game to a ZIP archive"),
		mcp.WithString("exe_path", mcp.Required(), mcp.Description("Path to the game executable")),
		mcp.WithString("output_path", mcp.Required(), mcp.Description("Path for the output ZIP file")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		exePath, err := req.RequireString("exe_path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		outputPath, err := req.RequireString("output_path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		gameDir := filepath.Dir(exePath)

		// Read patch summary
		summaryPath := filepath.Join(gameDir, "patch-summary.json")
		summaryData, err := os.ReadFile(summaryPath)
		if err != nil {
			return mcp.NewToolResultError("patch-summary.json not found - game may not be patched"), nil
		}

		var patchSummary domain.PatchSummary
		if err := json.Unmarshal(summaryData, &patchSummary); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to parse patch-summary.json: %v", err)), nil
		}

		if len(patchSummary.PatchedFiles) == 0 {
			return mcp.NewToolResultError("no patched files found in patch summary"), nil
		}

		// Create ZIP file
		zipFile, err := os.Create(outputPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create ZIP file: %v", err)), nil
		}
		defer zipFile.Close()

		zipWriter := zip.NewWriter(zipFile)
		defer zipWriter.Close()

		filesAdded := 0
		for _, relPath := range patchSummary.PatchedFiles {
			srcPath := filepath.Join(gameDir, relPath)
			if _, err := os.Stat(srcPath); os.IsNotExist(err) {
				continue
			}

			srcFile, err := os.Open(srcPath)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to open %s: %v", relPath, err)), nil
			}

			zipPath := filepath.ToSlash(relPath)
			writer, err := zipWriter.Create(zipPath)
			if err != nil {
				srcFile.Close()
				return mcp.NewToolResultError(fmt.Sprintf("failed to create ZIP entry %s: %v", relPath, err)), nil
			}

			_, err = io.Copy(writer, srcFile)
			srcFile.Close()
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to write %s to ZIP: %v", relPath, err)), nil
			}
			filesAdded++
		}

		return mcp.NewToolResultText(fmt.Sprintf("exported %d patched files to %s", filesAdded, outputPath)), nil
	})
}
