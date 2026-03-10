package main

import (
	"context"
	"flag"
	"fmt"
	"htpatcher/internal/domain"
	"htpatcher/internal/repository"
	"htpatcher/internal/service"
	"os"
)

// CLILogger prints log messages to stdout/stderr
type CLILogger struct{}

func (l *CLILogger) Info(message string)    { fmt.Println("[INFO]", message) }
func (l *CLILogger) Success(message string) { fmt.Println("[OK]  ", message) }
func (l *CLILogger) Error(message string)   { fmt.Fprintln(os.Stderr, "[ERR] ", message) }
func (l *CLILogger) Warn(message string)    { fmt.Println("[WARN]", message) }

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	logger := &CLILogger{}
	patchRepo := repository.NewPatchRepository()
	gameService := service.NewGameService(logger)
	patchService := service.NewPatchService(patchRepo, logger)
	backupService := service.NewBackupService(logger)
	downloadService := service.NewDownloadService(patchRepo, logger)

	switch os.Args[1] {
	case "patch":
		cmdApply(os.Args[2:], gameService, patchService, backupService, downloadService)
	case "restore":
		cmdRestore(os.Args[2:], gameService, backupService)
	case "list":
		cmdList(patchService)
	case "version":
		fmt.Printf("htpatcher v%d\n", service.Version)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Usage: htpatcher-cli <command> [options]

Commands:
  patch      Apply a patch to a game (local file or downloaded by store code)
  restore    Restore a game to its pre-patch state
  list       List all available patches
  version    Print the current version

Run "htpatcher-cli <command> -help" for command-specific options.`)
}

func cmdApply(args []string, gameSvc *service.GameService, patchSvc *service.PatchService, backupSvc *service.BackupService, dlSvc *service.DownloadService) {
	fs := flag.NewFlagSet("patch", flag.ExitOnError)
	gamePath := fs.String("game", "", "Path to the game executable (required)")
	patchPath := fs.String("patch", "", "Path to a local .htpatch file")
	storeCode := fs.String("store-code", "", "Store code of the game (required)")
	backup := fs.Bool("backup", true, "Create a backup before patching")
	launch := fs.Bool("launch", false, "Launch the game after patching (Windows only)")
	msgHide := fs.Bool("message-hide", false, "Inject MessageWindowHidden plugin")
	noPatchChecker := fs.Bool("no-patch-checker", false, "Disable HT_PatchChecker plugin injection")
	fs.Parse(args)

	if *gamePath == "" || *storeCode == "" {
		fmt.Fprintln(os.Stderr, "error: --game and --store-code are required")
		fs.Usage()
		os.Exit(1)
	}


	gameInfo, err := gameSvc.GetGameInfoFromExePath(*gamePath)
	if err != nil {
		die("error reading game: %v", err)
	}

	var patchInfo *domain.PatchInfo
	patchVersion := "0"

	if *patchPath == "" {
		gamePatchInfo, err := patchSvc.FetchGamePatchInfo(*storeCode)
		if err != nil {
			die("error fetching patch info: %v", err)
		}
		if gamePatchInfo == nil || len(gamePatchInfo.Patches) == 0 {
			die("no patch found for store code: %s", *storeCode)
		}
		latest := gamePatchInfo.Patches[0]
		if latest.Download == nil {
			die("patch has no download URL")
		}
		patchInfo, err = dlSvc.DownloadPatch(latest.Download.Url, latest.Download.FileName, func(filePath string) (*domain.PatchInfo, error) {
			return patchSvc.LoadPatchInfo(filePath)
		})
		if err != nil {
			die("download failed: %v", err)
		}
		patchVersion = latest.Version
	} else {
		patchInfo, err = patchSvc.LoadPatchInfo(*patchPath)
		if err != nil {
			die("error loading patch: %v", err)
		}
	}

	if *backup {
		if err := backupSvc.BackupGameData(gameInfo, patchInfo, *msgHide); err != nil {
			die("backup failed: %v", err)
		}
	}

	if err := patchSvc.ApplyPatch(context.Background(), gameInfo, patchInfo, *msgHide, !*noPatchChecker, *storeCode, patchVersion); err != nil {
		die("patch failed: %v", err)
	}

	if *launch {
		if err := gameSvc.LaunchGame(gameInfo.ExePath); err != nil {
			die("launch failed: %v", err)
		}
	}
}

func cmdRestore(args []string, gameSvc *service.GameService, backupSvc *service.BackupService) {
	fs := flag.NewFlagSet("restore", flag.ExitOnError)
	gamePath := fs.String("game", "", "Path to the game executable (required)")
	fs.Parse(args)

	if *gamePath == "" {
		fmt.Fprintln(os.Stderr, "error: --game is required")
		fs.Usage()
		os.Exit(1)
	}

	gameInfo, err := gameSvc.GetGameInfoFromExePath(*gamePath)
	if err != nil {
		die("error reading game: %v", err)
	}

	if err := backupSvc.RestoreBackup(gameInfo); err != nil {
		die("restore failed: %v", err)
	}
}

func cmdList(patchSvc *service.PatchService) {
	patches, err := patchSvc.FetchAllPatches()
	if err != nil {
		die("error fetching patches: %v", err)
	}

	if len(patches) == 0 {
		fmt.Println("No patches available.")
		return
	}

	fmt.Printf("%-22s  %-8s  %s\n", "STORE CODE", "STORE", "TITLE")
	fmt.Println("--------------------------------------------------")
	for _, p := range patches {
		fmt.Printf("%-22s  %-8s  %s\n", p.StoreCode, p.Store, p.Title)
	}
}

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}
