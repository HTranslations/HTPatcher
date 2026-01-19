<script lang="ts">
  import { domain } from "../../wailsjs/go/models.js";
  import PageHeader from "./PageHeader.svelte";
  import GameCard from "./GameCard.svelte";

  export let games: domain.LocatedGame[];
  export let searchQuery: string;
  export let gamesPerRow: number;
  export let onAddGame: () => void;
  export let onOpenFolder: (game: domain.LocatedGame) => void;
  export let onLaunchGame: (game: domain.LocatedGame) => void;
  export let onLocateGame: (game: domain.LocatedGame) => void;
  export let onTranslateGame: (game: domain.LocatedGame) => void;
  export let onEditGame: (game: domain.LocatedGame) => void;
  export let onPinGame: (game: domain.LocatedGame) => void;
  export let onSearchQueryChange: (query: string) => void;
  export let onGamesPerRowChange: (count: number) => void;
</script>

<PageHeader
  title="Your Games"
  description="Select a game to view available translations and apply patches"
/>

<!-- Search Bar -->
<div class="mb-4">
  <div class="relative">
    <input
      type="text"
      value={searchQuery}
      oninput={(e) => onSearchQueryChange((e.target as HTMLInputElement).value)}
      placeholder="Search by game name, RJ code, or tags..."
      class="w-full bg-zinc-900 border border-zinc-800 px-4 py-3 pl-10 text-sm text-zinc-300 focus:outline-none focus:border-zinc-700 transition-colors"
    />
    <svg
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
      stroke-width="1.5"
      stroke="currentColor"
      class="w-5 h-5 text-zinc-500 absolute left-3 top-1/2 -translate-y-1/2"
    >
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z"
      />
    </svg>
  </div>
</div>

<!-- Games Per Row Toggle -->
<div class="flex items-center justify-end gap-3 mb-4">
  <span class="text-sm text-zinc-400">Games per row:</span>
  <div class="flex gap-2">
    <button
      onclick={() => onGamesPerRowChange(3)}
      class="px-4 py-2 text-sm font-medium transition-colors {gamesPerRow === 3
        ? 'bg-emerald-600 text-white'
        : 'bg-zinc-800 text-zinc-300 hover:bg-zinc-700'}"
    >
      3
    </button>
    <button
      onclick={() => onGamesPerRowChange(4)}
      class="px-4 py-2 text-sm font-medium transition-colors {gamesPerRow === 4
        ? 'bg-emerald-600 text-white'
        : 'bg-zinc-800 text-zinc-300 hover:bg-zinc-700'}"
    >
      4
    </button>
  </div>
</div>

<div
  class="grid grid-cols-1 md:grid-cols-2 {gamesPerRow === 3
    ? 'lg:grid-cols-3'
    : 'lg:grid-cols-4'} gap-4"
>
  <!-- Add Game Card -->
  <button
    onclick={onAddGame}
    class="bg-zinc-900 border border-zinc-800 text-left flex flex-col"
  >
    <div class="aspect-4/3 bg-zinc-800 flex items-center justify-center">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        class="w-12 h-12 text-zinc-500"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M12 4.5v15m7.5-7.5h-15"
        />
      </svg>
    </div>
    <div class="p-4 flex flex-col flex-1">
      <h3 class="font-medium text-sm mb-1">Add Game</h3>
      <p class="text-xs text-zinc-500 mb-3">Locate game from disk</p>
      <div
        class="mt-auto w-full bg-zinc-800 border border-zinc-700 text-zinc-300 text-sm px-3 py-2 text-center"
      >
        Select Game to Add
      </div>
    </div>
  </button>

  {#each games as game}
    <GameCard
      {game}
      onOpenGame={onTranslateGame}
      {onOpenFolder}
      {onLaunchGame}
      {onLocateGame}
      {onEditGame}
      {onPinGame}
    />
  {/each}
</div>
