<script lang="ts">
  import "./app.css";
  import { onMount } from "svelte";
  import { domain } from "../wailsjs/go/models.js";
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
    SetGamePinned,
    SetGamePlayStatus,
    RestoreGameBackup,
    GetPersistentDataPath,
    DeletePersistentData,
    UpdateGameMetadata,
    GetGamesPerRow,
    SetGamesPerRow,
    CheckForUpdate,
    GetCurrentVersion,
    ExportPatchedFiles,
  } from "../wailsjs/go/main/App.js";
  import { EventsOn } from "../wailsjs/runtime/runtime.js";

  import Header from "./components/Header.svelte";
  import Sidebar from "./components/Sidebar.svelte";
  import GamesPage from "./components/GamesPage.svelte";
  import PatchesPage from "./components/PatchesPage.svelte";
  import SettingsPage from "./components/SettingsPage.svelte";
  import RequestTranslationPage from "./components/RequestTranslationPage.svelte";
  import AddGameDialog from "./components/AddGameDialog.svelte";
  import ConfirmDialog from "./components/ConfirmDialog.svelte";
  import TranslateGameDrawer from "./components/TranslateGameDrawer.svelte";
  import RestoreBackupDrawer from "./components/RestoreBackupDrawer.svelte";
  import EditGameDrawer from "./components/EditGameDrawer.svelte";
  import UpdateDialog from "./components/UpdateDialog.svelte";
  import { getDlsiteImageUrl } from "./lib/utils.js";

  let games: domain.LocatedGame[] = [];
  let patches: domain.PatchEntry[] = [];
  let currentPage: "games" | "patches" | "settings" | "request-translation" =
    "games";
  let selectedCategory = "all";
  let searchQuery = "";
  let showAddDialog = false;
  let locatedGame: domain.LocatedGame | null = null;
  let rjCode = "";
  let friendlyName = "";
  let tags = "";
  let previewImageUrl = "";
  let previewImageLoaded = false;
  let showDeleteDialog = false;
  let gameToDelete: domain.LocatedGame | null = null;
  let showDeleteDataDialog = false;
  let dataPath = "";

  // Translate game drawer state
  let showTranslateDrawer = false;
  let translateGameInfo: domain.GameInfo | null = null;
  let translatePatchInfo: domain.PatchInfo | null = null;
  let translateLogs: Array<{
    message: string;
    type: "info" | "success" | "error";
  }> = [];
  let isPatching = false;
  let launchAfterPatch = true;
  let selectedPatch: domain.PatchEntry | null = null;
  let patchSearchQuery = "";
  let currentTranslatingGame: domain.LocatedGame | null = null;

  // Restore backup drawer state
  let showRestoreDrawer = false;
  let restoreGameInfo: domain.GameInfo | null = null;
  let restoreLogs: Array<{
    message: string;
    type: "info" | "success" | "error";
  }> = [];
  let isRestoring = false;
  let currentRestoringGame: domain.LocatedGame | null = null;

  // Edit game drawer state
  let showEditDrawer = false;
  let gameToEdit: domain.LocatedGame | null = null;

  // Update dialog state
  let showUpdateDialog = false;
  let updateReleaseInfo: domain.ReleaseInfo | null = null;

  // Update success toast state
  let showUpdateSuccessToast = false;
  let updatedToVersion = 0;

  // Settings
  let gamesPerRow = 3;

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

  onMount(async () => {
    loadGames();
    loadPatches();

    // Load data path
    try {
      dataPath = await GetPersistentDataPath();
    } catch (error) {
      console.error("Failed to get data path:", error);
    }

    // Load games per row setting
    try {
      gamesPerRow = await GetGamesPerRow();
    } catch (error) {
      console.error("Failed to get games per row:", error);
      gamesPerRow = 3;
    }

    // Check for updates
    try {
      const releaseInfo = await CheckForUpdate();
      if (releaseInfo) {
        updateReleaseInfo = releaseInfo;
        showUpdateDialog = true;
      }
    } catch (error) {
      console.error("Failed to check for updates:", error);
    }

    // Listen for log events
    EventsOn(
      "log",
      (logMessage: { message: string; type: "info" | "success" | "error" }) => {
        if (showTranslateDrawer) {
          translateLogs = [...translateLogs, logMessage];
        } else if (showRestoreDrawer) {
          restoreLogs = [...restoreLogs, logMessage];
        }
      },
    );

    // Listen for update success event
    EventsOn("app:updated", (version: number) => {
      updatedToVersion = version;
      showUpdateSuccessToast = true;
    });
  });

  $: categories = [
    { id: "all", name: "All Games", count: games.length },
    {
      id: "patched",
      name: "Patched",
      count: games.filter((g) => g.translated).length,
    },
    {
      id: "not-patched",
      name: "Not Patched",
      count: games.filter((g) => !g.translated).length,
    },
  ];

  $: playStatusCategories = [
    {
      id: "unplayed",
      name: "Unplayed",
      count: games.filter((g) => !g.playStatus || g.playStatus === "unplayed").length,
    },
    {
      id: "playing",
      name: "Playing",
      count: games.filter((g) => g.playStatus === "playing").length,
    },
    {
      id: "on-hold",
      name: "On Hold",
      count: games.filter((g) => g.playStatus === "on-hold").length,
    },
    {
      id: "finished",
      name: "Finished",
      count: games.filter((g) => g.playStatus === "finished").length,
    },
    {
      id: "given-up",
      name: "Given Up",
      count: games.filter((g) => g.playStatus === "given-up").length,
    },
  ];

  $: filteredGames = (() => {
    let filtered = games;

    // Apply category filter
    if (selectedCategory === "patched") {
      filtered = filtered.filter((g) => g.translated);
    } else if (selectedCategory === "not-patched") {
      filtered = filtered.filter((g) => !g.translated);
    } else if (selectedCategory === "unplayed") {
      filtered = filtered.filter((g) => !g.playStatus || g.playStatus === "unplayed");
    } else if (selectedCategory === "playing") {
      filtered = filtered.filter((g) => g.playStatus === "playing");
    } else if (selectedCategory === "on-hold") {
      filtered = filtered.filter((g) => g.playStatus === "on-hold");
    } else if (selectedCategory === "finished") {
      filtered = filtered.filter((g) => g.playStatus === "finished");
    } else if (selectedCategory === "given-up") {
      filtered = filtered.filter((g) => g.playStatus === "given-up");
    }

    // Apply search filter
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter((g) => {
        const friendlyNameMatch = g.friendlyName?.toLowerCase().includes(query);
        const rjCodeMatch = g.rjCode?.toLowerCase().includes(query);
        const tagsMatch = g.tags?.some((tag) =>
          tag.toLowerCase().includes(query),
        );
        return friendlyNameMatch || rjCodeMatch || tagsMatch;
      });
    }

    // Sort pinned games first (stable sort preserves original order)
    return [...filtered].sort((a, b) => {
      if (a.pinned && !b.pinned) return -1;
      if (!a.pinned && b.pinned) return 1;
      return 0;
    });
  })();

  $: {
    const newUrl = rjCode ? getDlsiteImageUrl(rjCode) : "";
    if (newUrl !== previewImageUrl) {
      previewImageLoaded = false;
    }
    previewImageUrl = newUrl;
  }

  function locateGame(game: domain.LocatedGame) {
    // Placeholder for file selection
    console.log("Locate game:", game);
  }

  async function openGame(game: domain.LocatedGame) {
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
        selectedPatch =
          patches.find(
            (patch) => patch.systemGameTitle === gameInfo.gameTitle,
          ) || null;
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

  function togglePatch(patch: domain.PatchEntry) {
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
        await ApplyPatch(
          translateGameInfo,
          translatePatchInfo,
          launchAfterPatch,
          true,
        );
        await SetGameTranslated(currentTranslatingGame.id, true);
        await loadGames();
      }
    } catch (error) {
      translateLogs = [
        ...translateLogs,
        { message: `Error: ${error}`, type: "error" },
      ];
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
      restoreLogs = [
        ...restoreLogs,
        { message: `Error: ${error}`, type: "error" },
      ];
    } finally {
      isRestoring = false;
    }
  }

  async function launchGame(game: domain.LocatedGame) {
    if (game.exePath) {
      try {
        await LaunchGameFromPath(game.exePath);
      } catch (error) {
        console.error("Failed to launch game:", error);
      }
    }
  }

  async function openGameFolder(game: domain.LocatedGame) {
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
      friendlyName = "";
      tags = "";
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
    friendlyName = "";
    tags = "";
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

  function handleFriendlyNameInput(event: Event) {
    const target = event.target as HTMLInputElement;
    friendlyName = target.value;
  }

  function handleTagsInput(event: Event) {
    const target = event.target as HTMLInputElement;
    tags = target.value;
  }

  async function addGameToCollection() {
    if (!locatedGame || !rjCode || !friendlyName) return;

    try {
      // Parse tags from comma-separated string
      const tagArray = tags
        .split(",")
        .map((tag) => tag.trim())
        .filter((tag) => tag.length > 0);

      await AddGameToCollection(locatedGame, rjCode, friendlyName, tagArray);
      await loadGames();
      closeAddDialog();
    } catch (error) {
      console.error("Failed to add game to collection:", error);
    }
  }

  function requestDeleteGame(game: domain.LocatedGame) {
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

  function requestDeleteData() {
    showDeleteDataDialog = true;
  }

  function cancelDeleteData() {
    showDeleteDataDialog = false;
  }

  async function confirmDeleteData() {
    try {
      await DeletePersistentData();
      await loadGames();
      cancelDeleteData();
    } catch (error) {
      console.error("Failed to delete data:", error);
    }
  }

  function openEditDrawer(game: domain.LocatedGame) {
    gameToEdit = game;
    showEditDrawer = true;
  }

  function closeEditDrawer() {
    showEditDrawer = false;
    gameToEdit = null;
  }

  async function saveGameEdit(friendlyName: string, tags: string[], playStatus: string) {
    if (!gameToEdit || !gameToEdit.id) return;

    try {
      await UpdateGameMetadata(gameToEdit.id, friendlyName, tags);
      await SetGamePlayStatus(gameToEdit.id, playStatus);
      await loadGames();
      closeEditDrawer();
    } catch (error) {
      console.error("Failed to update game metadata:", error);
    }
  }

  async function exportPatchedFiles(game: domain.LocatedGame) {
    if (!game.gameDir) return;

    try {
      await ExportPatchedFiles(game.gameDir, game.friendlyName || game.rjCode || "game");
    } catch (error) {
      console.error("Failed to export patched files:", error);
    }
  }

  async function handleGamesPerRowChange(count: number) {
    try {
      await SetGamesPerRow(count);
      gamesPerRow = count;
    } catch (error) {
      console.error("Failed to set games per row:", error);
    }
  }

  async function togglePinGame(game: domain.LocatedGame) {
    try {
      await SetGamePinned(game.id, !game.pinned);
      await loadGames();
    } catch (error) {
      console.error("Failed to toggle pin:", error);
    }
  }

  function closeUpdateDialog() {
    showUpdateDialog = false;
    updateReleaseInfo = null;
  }

  function closeUpdateSuccessToast() {
    showUpdateSuccessToast = false;
  }
</script>

<div class="flex flex-col h-screen bg-zinc-950 text-zinc-100">
  <Header />

  <!-- Main Content -->
  <div class="flex flex-1 overflow-hidden">
    <Sidebar bind:currentPage bind:selectedCategory {categories} {playStatusCategories} />

    <!-- Content Area -->
    <div class="flex-1 overflow-y-auto">
      <div class="p-6">
        {#if currentPage === "games"}
          <GamesPage
            games={filteredGames}
            {searchQuery}
            {gamesPerRow}
            onAddGame={addGame}
            onOpenFolder={openGameFolder}
            onLaunchGame={launchGame}
            onLocateGame={locateGame}
            onTranslateGame={openGame}
            onEditGame={openEditDrawer}
            onPinGame={togglePinGame}
            onSearchQueryChange={(query) => (searchQuery = query)}
            onGamesPerRowChange={handleGamesPerRowChange}
          />
        {:else if currentPage === "patches"}
          <PatchesPage {patches} />
        {:else if currentPage === "settings"}
          <SettingsPage {dataPath} onDeleteData={requestDeleteData} />
        {:else if currentPage === "request-translation"}
          <RequestTranslationPage />
        {/if}
      </div>
    </div>
  </div>

  <AddGameDialog
    show={showAddDialog}
    {locatedGame}
    bind:rjCode
    bind:friendlyName
    bind:tags
    {previewImageUrl}
    {previewImageLoaded}
    onClose={closeAddDialog}
    onRjCodeInput={handleRjCodeInput}
    onFriendlyNameInput={handleFriendlyNameInput}
    onTagsInput={handleTagsInput}
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

  <ConfirmDialog
    show={showDeleteDataDialog}
    title="Delete All Data"
    message="Are you sure you want to delete all your persistent data? This will remove all games from your collection and cannot be undone."
    confirmText="Delete All Data"
    cancelText="Cancel"
    onConfirm={confirmDeleteData}
    onCancel={cancelDeleteData}
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
    onLaunchAfterPatchChange={(value) => (launchAfterPatch = value)}
    onPatchSearchQueryChange={(value) => (patchSearchQuery = value)}
  />

  <RestoreBackupDrawer
    show={showRestoreDrawer}
    gameInfo={restoreGameInfo}
    logs={restoreLogs}
    {isRestoring}
    onClose={closeRestoreDrawer}
    onRestoreBackup={restoreBackup}
  />

  <EditGameDrawer
    show={showEditDrawer}
    game={gameToEdit}
    onClose={closeEditDrawer}
    onSave={saveGameEdit}
    onDelete={requestDeleteGame}
    onExport={exportPatchedFiles}
  />

  <UpdateDialog
    show={showUpdateDialog}
    releaseInfo={updateReleaseInfo}
    onClose={closeUpdateDialog}
  />

  <!-- Update Success Toast -->
  {#if showUpdateSuccessToast}
    <div class="fixed bottom-4 right-4 z-50 bg-green-600 border border-green-500 text-white px-4 py-3 shadow-lg flex items-center gap-3">
      <span>Successfully updated to v{updatedToVersion}!</span>
      <button
        onclick={closeUpdateSuccessToast}
        class="text-white/80 hover:text-white cursor-pointer"
        aria-label="Close notification"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
        </svg>
      </button>
    </div>
  {/if}
</div>
