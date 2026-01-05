<script lang="ts">
  import { main } from "../../wailsjs/go/models.js";
  import { getDlsiteImageUrl } from "../lib/utils.js";
  
  export let game: main.LocatedGame;
  export let onOpenGame: (game: main.LocatedGame) => void;
  export let onOpenFolder: (game: main.LocatedGame) => void;
  export let onLaunchGame: (game: main.LocatedGame) => void;
  export let onLocateGame: (game: main.LocatedGame) => void;
  export let onDeleteGame: (game: main.LocatedGame) => void;

  let loadingFolder = false;
  let loadingPlay = false;

  async function handleOpenFolder() {
    loadingFolder = true;
    await new Promise(resolve => setTimeout(resolve, 1000));
    onOpenFolder(game);
    loadingFolder = false;
  }

  async function handleLaunchGame() {
    loadingPlay = true;
    await new Promise(resolve => setTimeout(resolve, 1000));
    onLaunchGame(game);
    loadingPlay = false;
  }
</script>

<div class="bg-zinc-900 border border-zinc-800">
  <div class="relative">
    <img
      src={getDlsiteImageUrl(game.rjCode)}
      alt={game.rjCode || "Game"}
      class="w-full aspect-4/3 object-cover bg-zinc-800"
    />
    {#if game.exePath}
      {#if game.translated}
        <div class="absolute top-2 right-2 bg-emerald-500 text-white text-xs font-bold px-2 py-1 shadow-lg">
          PATCHED
        </div>
      {:else}
        <div class="absolute top-2 right-2 bg-blue-500 text-white text-xs font-bold px-2 py-1 shadow-lg">
          NOT PATCHED
        </div>
      {/if}
    {/if}
  </div>
  <div class="p-4 text-left">
    <h3 class="font-medium text-sm mb-1 line-clamp-2">{game.rjCode || "Unknown Game"}</h3>
    <p class="text-xs text-zinc-500 mb-3 font-mono truncate">{game.exePath || game.gameDir}</p>
    {#if game.exePath}
      <div class="flex items-center gap-2">
        <button
          onclick={() => onDeleteGame(game)}
          class="bg-zinc-800 border border-zinc-700 text-zinc-300 text-sm px-3 py-2 font-medium cursor-pointer"
          title="Delete game"
        >
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
          </svg>
        </button>
        <button
          onclick={() => onOpenGame(game)}
          class="flex-1 bg-zinc-800 border border-zinc-700 text-zinc-300 text-sm px-3 py-2 font-medium cursor-pointer"
        >
          {game.translated ? "Remove Patch" : "Patch Game"}
        </button>
        <button
          onclick={handleOpenFolder}
          disabled={loadingFolder}
          class="bg-yellow-600 border border-yellow-600 text-white px-3 py-2 cursor-pointer disabled:opacity-50 disabled:cursor-wait"
          title="Open folder"
        >
          {#if loadingFolder}
            <svg class="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          {:else}
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 12.75V12A2.25 2.25 0 014.5 9.75h15A2.25 2.25 0 0121.75 12v.75m-8.69-6.44l-2.12-2.12a1.5 1.5 0 00-1.061-.44H4.5A2.25 2.25 0 002.25 6v12a2.25 2.25 0 002.25 2.25h15A2.25 2.25 0 0021.75 18V9a2.25 2.25 0 00-2.25-2.25h-5.379a1.5 1.5 0 01-1.06-.44z" />
            </svg>
          {/if}
        </button>
        <button
          onclick={handleLaunchGame}
          disabled={loadingPlay}
          class="bg-blue-600 border border-blue-600 text-white px-3 py-2 cursor-pointer disabled:opacity-50 disabled:cursor-wait"
          title="Launch game"
        >
          {#if loadingPlay}
            <svg class="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          {:else}
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M5.25 5.653c0-.856.917-1.398 1.667-.986l11.54 6.348a1.125 1.125 0 010 1.971l-11.54 6.347a1.125 1.125 0 01-1.667-.986V5.653z" />
            </svg>
          {/if}
        </button>
      </div>
    {:else}
      <button
        onclick={() => onLocateGame(game)}
        class="w-full bg-zinc-800 border border-zinc-700 text-zinc-300 text-sm px-3 py-2"
      >
        Locate Game
      </button>
    {/if}
  </div>
</div>

