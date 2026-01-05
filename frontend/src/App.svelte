<script lang="ts">
  import "./app.css";
  import { onMount } from "svelte";
  import { main } from "../wailsjs/go/models.js";
  import { 
    PrepareGameToAddToCollection,
    AddGameToCollection,
    GetGamesCollection,
    LaunchGameFromPath,
    OpenFolder,
    FetchAllPatches,
    RemoveGameFromCollection,
    SelectGameExeFile,
    GetGameInfoFromExePath,
    SelectPatchFile,
    DownloadPatch,
    ApplyPatch,
    SetGameTranslated,
    RestoreGameBackup
  } from "../wailsjs/go/main/App.js";
  import { EventsOn } from "../wailsjs/runtime/runtime.js";
  
  import Header from "./components/Header.svelte";
  import Sidebar from "./components/Sidebar.svelte";
  import GamesPage from "./components/GamesPage.svelte";
  import PatchesPage from "./components/PatchesPage.svelte";
  import AddGameDialog from "./components/AddGameDialog.svelte";
  import ConfirmDialog from "./components/ConfirmDialog.svelte";
  import TranslateGameDrawer from "./components/TranslateGameDrawer.svelte";
  import RestoreBackupDrawer from "./components/RestoreBackupDrawer.svelte";
  import { getDlsiteImageUrl } from "./lib/utils.js";

  let games: main.LocatedGame[] = [];
  let patches: main.PatchEntry[] = [];
  let currentPage: "games" | "patches" = "games";
  let selectedCategory = "all";
  let showAddDialog = false;
  let locatedGame: main.LocatedGame | null = null;
  let rjCode = "";
  let previewImageUrl = "";
  let previewImageLoaded = false;
  let showDeleteDialog = false;
  let gameToDelete: main.LocatedGame | null = null;
  
  // Translate game drawer state
  let showTranslateDrawer = false;
  let translateGameInfo: main.GameInfo | null = null;
  let translatePatchInfo: main.PatchInfo | null = null;
  let translateLogs: Array<{ message: string; type: "info" | "success" | "error" }> = [];
  let isPatching = false;
  let launchAfterPatch = true;
  let selectedPatch: main.PatchEntry | null = null;
  let patchSearchQuery = "";
  let currentTranslatingGame: main.LocatedGame | null = null;
  
  // Restore backup drawer state
  let showRestoreDrawer = false;
  let restoreGameInfo: main.GameInfo | null = null;
  let restoreLogs: Array<{ message: string; type: "info" | "success" | "error" }> = [];
  let isRestoring = false;
  let currentRestoringGame: main.LocatedGame | null = null;

  async function loadGames() {
    try {
      games = await GetGamesCollection();
    } catch (error) {
      console.error("Failed to load games:", error);
      games = [];
    }
  }

  async function loadPatches() {
    try {
      patches = await FetchAllPatches();
    } catch (error) {
      console.error("Failed to load patches:", error);
      patches = [];
    }
  }

  onMount(() => {
    loadGames();
    loadPatches();
    
    // Listen for log events
    EventsOn("log", (logMessage: { message: string; type: "info" | "success" | "error" }) => {
      if (showTranslateDrawer) {
        translateLogs = [...translateLogs, logMessage];
      } else if (showRestoreDrawer) {
        restoreLogs = [...restoreLogs, logMessage];
      }
    });
  });

  $: categories = [
    { id: "all", name: "All Games", count: games.length },
    { id: "patched", name: "Patched", count: games.filter(g => g.translated).length },
    { id: "not-patched", name: "Not Patched", count: games.filter(g => !g.translated).length },
  ];

  $: filteredGames = selectedCategory === "all" 
    ? games 
    : selectedCategory === "patched"
    ? games.filter(g => g.translated)
    : selectedCategory === "not-patched"
    ? games.filter(g => !g.translated)
    : games;

  $: {
    const newUrl = rjCode ? getDlsiteImageUrl(rjCode) : "";
    if (newUrl !== previewImageUrl) {
      previewImageLoaded = false;
    }
    previewImageUrl = newUrl;
  }

  function locateGame(game: main.LocatedGame) {
    // Placeholder for file selection
    console.log("Locate game:", game);
  }

  async function openGame(game: main.LocatedGame) {
    if (!game.exePath) return;
    
    try {
      const gameInfo = await GetGameInfoFromExePath(game.exePath);
      if (!gameInfo) return;
      
      if (game.translated) {
        // Show restore backup drawer
        restoreGameInfo = gameInfo;
        currentRestoringGame = game;
        restoreLogs = [];
        showRestoreDrawer = true;
      } else {
        // Show translate drawer
        translateGameInfo = gameInfo;
        currentTranslatingGame = game;
        // Try to find a matching patch
        selectedPatch = patches.find(patch => patch.systemGameTitle === gameInfo.gameTitle) || null;
        translateLogs = [];
        translatePatchInfo = null;
        patchSearchQuery = "";
        showTranslateDrawer = true;
      }
    } catch (error) {
      console.error("Failed to get game info:", error);
    }
  }
  
  function closeTranslateDrawer() {
    showTranslateDrawer = false;
    translateGameInfo = null;
    translatePatchInfo = null;
    selectedPatch = null;
    translateLogs = [];
    patchSearchQuery = "";
    currentTranslatingGame = null;
  }
  
  function closeRestoreDrawer() {
    showRestoreDrawer = false;
    restoreGameInfo = null;
    restoreLogs = [];
    currentRestoringGame = null;
  }
  
  async function selectPatchFile() {
    try {
      translatePatchInfo = await SelectPatchFile();
      selectedPatch = null;
    } catch (error) {
      console.error("Failed to select patch file:", error);
    }
  }
  
  function togglePatch(patch: main.PatchEntry) {
    if (selectedPatch === patch) {
      selectedPatch = null;
    } else {
      selectedPatch = patch;
      translatePatchInfo = null;
    }
  }
  
  function clearCustomPatch() {
    translatePatchInfo = null;
  }
  
  async function applyPatch() {
    if (!translateGameInfo || !(translatePatchInfo || selectedPatch)) return;
    
    translateLogs = [];
    isPatching = true;
    
    try {
      if (selectedPatch) {
        translatePatchInfo = await DownloadPatch(selectedPatch.patchDownloadId);
      }
      
      if (translatePatchInfo && currentTranslatingGame) {
        await ApplyPatch(translateGameInfo, translatePatchInfo, launchAfterPatch, true);
        await SetGameTranslated(currentTranslatingGame.id, true);
        await loadGames();
      }
    } catch (error) {
      translateLogs = [...translateLogs, { message: `Error: ${error}`, type: "error" }];
    } finally {
      isPatching = false;
    }
  }
  
  async function restoreBackup() {
    if (!restoreGameInfo || !currentRestoringGame) return;
    
    restoreLogs = [];
    isRestoring = true;
    
    try {
      await RestoreGameBackup(restoreGameInfo);
      await SetGameTranslated(currentRestoringGame.id, false);
      await loadGames();
    } catch (error) {
      restoreLogs = [...restoreLogs, { message: `Error: ${error}`, type: "error" }];
    } finally {
      isRestoring = false;
    }
  }

  async function launchGame(game: main.LocatedGame) {
    if (game.exePath) {
      try {
        await LaunchGameFromPath(game.exePath);
      } catch (error) {
        console.error("Failed to launch game:", error);
      }
    }
  }

  async function openGameFolder(game: main.LocatedGame) {
    const folderPath = game.exePath ? game.gameDir : game.gameDir;
    if (folderPath) {
      try {
        await OpenFolder(folderPath);
      } catch (error) {
        console.error("Failed to open folder:", error);
      }
    }
  }

  async function addGame() {
    try {
      const game = await PrepareGameToAddToCollection();
      locatedGame = game;
      rjCode = "";
      previewImageUrl = "";
      previewImageLoaded = false;
      showAddDialog = true;
    } catch (error) {
      console.error("Failed to prepare game:", error);
    }
  }

  function closeAddDialog() {
    showAddDialog = false;
    locatedGame = null;
    rjCode = "";
    previewImageUrl = "";
    previewImageLoaded = false;
  }

  function handlePreviewImageLoad() {
    previewImageLoaded = true;
  }

  function handlePreviewImageError() {
    previewImageLoaded = false;
  }

  function handleRjCodeInput(event: Event) {
    const target = event.target as HTMLInputElement;
    let value = target.value.toUpperCase();
    // Remove any non-alphanumeric characters except RJ prefix
    value = value.replace(/[^RJ0-9]/g, "");
    // Ensure it starts with RJ
    if (value && !value.startsWith("RJ")) {
      value = "RJ" + value.replace(/RJ/g, "");
    }
    // Limit to reasonable length (RJ + 8 digits max)
    if (value.length > 10) {
      value = value.substring(0, 10);
    }
    rjCode = value;
  }

  async function addGameToCollection() {
    if (!locatedGame || !rjCode) return;
    
    try {
      await AddGameToCollection(locatedGame, rjCode);
      await loadGames();
      closeAddDialog();
    } catch (error) {
      console.error("Failed to add game to collection:", error);
    }
  }

  function requestDeleteGame(game: main.LocatedGame) {
    gameToDelete = game;
    showDeleteDialog = true;
  }

  function cancelDelete() {
    showDeleteDialog = false;
    gameToDelete = null;
  }

  async function confirmDelete() {
    if (!gameToDelete || !gameToDelete.id) return;
    
    try {
      await RemoveGameFromCollection(gameToDelete.id);
      await loadGames();
      cancelDelete();
    } catch (error) {
      console.error("Failed to delete game:", error);
    }
  }
</script>

<div class="flex flex-col h-screen bg-zinc-950 text-zinc-100">
  <Header />

  <!-- Main Content -->
  <div class="flex flex-1 overflow-hidden">
    <Sidebar bind:currentPage bind:selectedCategory {categories} />

    <!-- Content Area -->
    <div class="flex-1 overflow-y-auto">
      <div class="p-6">
        {#if currentPage === "games"}
          <GamesPage
            games={filteredGames}
            onAddGame={addGame}
            onOpenFolder={openGameFolder}
            onLaunchGame={launchGame}
            onLocateGame={locateGame}
            onDeleteGame={requestDeleteGame}
            onTranslateGame={openGame}
          />
        {:else if currentPage === "patches"}
          <PatchesPage {patches} />
        {/if}
      </div>
    </div>
  </div>

  <AddGameDialog
    show={showAddDialog}
    {locatedGame}
    bind:rjCode
    {previewImageUrl}
    {previewImageLoaded}
    onClose={closeAddDialog}
    onRjCodeInput={handleRjCodeInput}
    onPreviewImageLoad={handlePreviewImageLoad}
    onPreviewImageError={handlePreviewImageError}
    onAddGame={addGameToCollection}
  />

  <ConfirmDialog
    show={showDeleteDialog}
    title="Delete Game"
    message="Are you sure you want to remove this game from your collection? This won't delete the game from your disk."
    confirmText="Delete"
    cancelText="Cancel"
    onConfirm={confirmDelete}
    onCancel={cancelDelete}
  />

  <TranslateGameDrawer
    show={showTranslateDrawer}
    gameInfo={translateGameInfo}
    {patches}
    logs={translateLogs}
    {isPatching}
    bind:launchAfterPatch
    {selectedPatch}
    patchInfo={translatePatchInfo}
    bind:patchSearchQuery
    onClose={closeTranslateDrawer}
    onSelectPatchFile={selectPatchFile}
    onTogglePatch={togglePatch}
    onClearCustomPatch={clearCustomPatch}
    onApplyPatch={applyPatch}
    onLaunchAfterPatchChange={(value) => launchAfterPatch = value}
    onPatchSearchQueryChange={(value) => patchSearchQuery = value}
  />

  <RestoreBackupDrawer
    show={showRestoreDrawer}
    gameInfo={restoreGameInfo}
    logs={restoreLogs}
    isRestoring={isRestoring}
    onClose={closeRestoreDrawer}
    onRestoreBackup={restoreBackup}
  />
</div>
