<script lang="ts">
  import { domain } from "../../wailsjs/go/models.js";
  import { getDlsiteImageUrl } from "../lib/utils.js";

  export let game: domain.LocatedGame;
  export let onOpenGame: (game: domain.LocatedGame) => void;
  export let onOpenFolder: (game: domain.LocatedGame) => void;
  export let onLaunchGame: (game: domain.LocatedGame) => void;
  export let onLocateGame: (game: domain.LocatedGame) => void;
  export let onEditGame: (game: domain.LocatedGame) => void;
  export let onPinGame: (game: domain.LocatedGame) => void;

  let loadingFolder = false;
  let loadingPlay = false;
  let isHovered = false;

  function getStatusColor(status: string | undefined): string {
    switch (status) {
      case "playing":
        return "bg-blue-500";
      case "on-hold":
        return "bg-amber-500";
      case "finished":
        return "bg-emerald-500";
      case "given-up":
        return "bg-red-500";
      case "unplayed":
      default:
        return "bg-zinc-600";
    }
  }

  async function handleOpenFolder() {
    loadingFolder = true;
    await new Promise((resolve) => setTimeout(resolve, 1000));
    onOpenFolder(game);
    loadingFolder = false;
  }

  async function handleLaunchGame() {
    loadingPlay = true;
    await new Promise((resolve) => setTimeout(resolve, 1000));
    onLaunchGame(game);
    loadingPlay = false;
  }
</script>

<div
  class="bg-zinc-900 border border-zinc-800"
  on:mouseenter={() => (isHovered = true)}
  on:mouseleave={() => (isHovered = false)}
>
  <div class="relative">
    <img
      src={getDlsiteImageUrl(game.rjCode)}
      alt={game.rjCode || "Game"}
      class="w-full aspect-4/3 object-cover bg-zinc-800"
    />
    <!-- Bookmark Button -->
    {#if game.pinned || isHovered}
      <button
        on:click={() => onPinGame(game)}
        class="absolute top-2 left-2 p-1.5 bg-zinc-900/70 cursor-pointer transition-opacity {game.pinned ? 'opacity-100' : 'opacity-0 group-hover:opacity-100'}"
        class:opacity-100={isHovered || game.pinned}
        title={game.pinned ? "Remove bookmark" : "Bookmark game"}
      >
        {#if game.pinned}
          <!-- Filled bookmark icon (blue) -->
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="currentColor"
            class="w-5 h-5 text-blue-500"
          >
            <path fill-rule="evenodd" d="M6.32 2.577a49.255 49.255 0 0 1 11.36 0c1.497.174 2.57 1.46 2.57 2.93V21a.75.75 0 0 1-1.085.67L12 18.089l-7.165 3.583A.75.75 0 0 1 3.75 21V5.507c0-1.47 1.073-2.756 2.57-2.93Z" clip-rule="evenodd" />
          </svg>
        {:else}
          <!-- Outline bookmark icon (gray) -->
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            class="w-5 h-5 text-zinc-400 hover:text-zinc-200"
          >
            <path stroke-linecap="round" stroke-linejoin="round" d="M17.593 3.322c1.1.128 1.907 1.077 1.907 2.185V21L12 17.25 4.5 21V5.507c0-1.108.806-2.057 1.907-2.185a48.507 48.507 0 0 1 11.186 0Z" />
          </svg>
        {/if}
      </button>
    {/if}
    {#if game.exePath}
      {#if game.translated}
        <div
          class="absolute top-2 right-2 bg-emerald-500 text-white text-xs font-bold px-2 py-1 shadow-lg"
        >
          PATCHED
        </div>
      {:else}
        <div
          class="absolute top-2 right-2 bg-blue-500 text-white text-xs font-bold px-2 py-1 shadow-lg"
        >
          NOT PATCHED
        </div>
      {/if}
    {/if}
  </div>
  <!-- Play Status Bar -->
  <div class="h-1.5 w-full {getStatusColor(game.playStatus)}"></div>
  <div class="p-4 text-left">
    <h3 class="font-medium text-sm mb-1 line-clamp-2">
      {game.friendlyName || game.rjCode || "Unknown Game"}
    </h3>
    <p class="text-xs text-zinc-500 mb-1 font-mono truncate">
      {game.rjCode || "No RJ Code"}
    </p>
    {#if game.tags && game.tags.length > 0}
      <div class="flex flex-wrap gap-1 mb-2">
        {#each game.tags as tag}
          <span
            class="text-xs bg-zinc-800 border border-zinc-700 text-zinc-400 px-2 py-0.5"
          >
            {tag}
          </span>
        {/each}
      </div>
    {/if}
    <p class="text-xs text-zinc-600 mb-3 font-mono truncate">
      {game.exePath || game.gameDir}
    </p>
    {#if game.exePath}
      <div class="flex items-center gap-2">
        <button
          on:click={() => onEditGame(game)}
          class="bg-zinc-800 border border-zinc-700 text-zinc-300 text-sm px-3 py-2 font-medium cursor-pointer"
          title="Edit game details"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            class="w-5 h-5"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M9.594 3.94c.09-.542.56-.94 1.11-.94h2.593c.55 0 1.02.398 1.11.94l.213 1.281c.063.374.313.686.645.87.074.04.147.083.22.127.324.196.72.257 1.075.124l1.217-.456a1.125 1.125 0 011.37.49l1.296 2.247a1.125 1.125 0 01-.26 1.431l-1.003.827c-.293.24-.438.613-.431.992a6.759 6.759 0 010 .255c-.007.378.138.75.43.99l1.005.828c.424.35.534.954.26 1.43l-1.298 2.247a1.125 1.125 0 01-1.369.491l-1.217-.456c-.355-.133-.75-.072-1.076.124a6.57 6.57 0 01-.22.128c-.331.183-.581.495-.644.869l-.213 1.28c-.09.543-.56.941-1.11.941h-2.594c-.55 0-1.02-.398-1.11-.94l-.213-1.281c-.062-.374-.312-.686-.644-.87a6.52 6.52 0 01-.22-.127c-.325-.196-.72-.257-1.076-.124l-1.217.456a1.125 1.125 0 01-1.369-.49l-1.297-2.247a1.125 1.125 0 01.26-1.431l1.004-.827c.292-.24.437-.613.43-.992a6.932 6.932 0 010-.255c.007-.378-.138-.75-.43-.99l-1.004-.828a1.125 1.125 0 01-.26-1.43l1.297-2.247a1.125 1.125 0 011.37-.491l1.216.456c.356.133.751.072 1.076-.124.072-.044.146-.087.22-.128.332-.183.582-.495.644-.869l.214-1.281z"
            />
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
            />
          </svg>
        </button>
        <button
          on:click={() => onOpenGame(game)}
          class="flex-1 bg-zinc-800 border border-zinc-700 text-zinc-300 text-sm px-3 py-2 font-medium cursor-pointer"
        >
          {game.translated ? "Remove Patch" : "Patch Game"}
        </button>
        <button
          on:click={handleOpenFolder}
          disabled={loadingFolder}
          class="bg-yellow-600 border border-yellow-600 text-white px-3 py-2 cursor-pointer disabled:opacity-50 disabled:cursor-wait"
          title="Open folder"
        >
          {#if loadingFolder}
            <svg
              class="animate-spin h-5 w-5"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"
              ></circle>
              <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              ></path>
            </svg>
          {:else}
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="w-5 h-5"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M2.25 12.75V12A2.25 2.25 0 014.5 9.75h15A2.25 2.25 0 0121.75 12v.75m-8.69-6.44l-2.12-2.12a1.5 1.5 0 00-1.061-.44H4.5A2.25 2.25 0 002.25 6v12a2.25 2.25 0 002.25 2.25h15A2.25 2.25 0 0121.75 18V9a2.25 2.25 0 00-2.25-2.25h-5.379a1.5 1.5 0 01-1.06-.44z"
              />
            </svg>
          {/if}
        </button>
        <button
          on:click={handleLaunchGame}
          disabled={loadingPlay}
          class="bg-blue-600 border border-blue-600 text-white px-3 py-2 cursor-pointer disabled:opacity-50 disabled:cursor-wait"
          title="Launch game"
        >
          {#if loadingPlay}
            <svg
              class="animate-spin h-5 w-5"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"
              ></circle>
              <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              ></path>
            </svg>
          {:else}
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="w-5 h-5"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M5.25 5.653c0-.856.917-1.398 1.667-.986l11.54 6.348a1.125 1.125 0 010 1.971l-11.54 6.347a1.125 1.125 0 01-1.667-.986V5.653z"
              />
            </svg>
          {/if}
        </button>
      </div>
    {:else}
      <button
        on:click={() => onLocateGame(game)}
        class="w-full bg-zinc-800 border border-zinc-700 text-zinc-300 text-sm px-3 py-2"
      >
        Locate Game
      </button>
    {/if}
  </div>
</div>
