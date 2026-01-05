<script lang="ts">
  import { main } from "../../wailsjs/go/models.js";
  import { BrowserOpenURL } from "../../wailsjs/runtime/runtime.js";
  import { getDlsiteImageUrl } from "../lib/utils.js";
  
  export let patch: main.PatchEntry;
  export let index: number;
</script>

<div class="bg-zinc-900 border border-zinc-800 hover:border-zinc-700 transition-colors">
  <div class="flex gap-3 py-3 px-3">
    <!-- Thumbnail -->
    <div class="shrink-0 w-40 self-stretch">
      {#if patch.rjCode && getDlsiteImageUrl(patch.rjCode)}
        <img
          src={getDlsiteImageUrl(patch.rjCode)}
          alt={patch.title}
          class="w-full h-full object-cover bg-zinc-800 border border-zinc-700"
          onerror={(e) => {
            const img = e.currentTarget as HTMLImageElement;
            img.style.display = 'none';
          }}
        />
      {:else}
        <div class="w-full h-full bg-zinc-800 border border-zinc-700 flex items-center justify-center">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-8 h-8 text-zinc-600">
            <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 15.75l5.159-5.159a2.25 2.25 0 013.182 0l5.159 5.159m-1.5-1.5l1.409-1.409a2.25 2.25 0 013.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 001.5-1.5V6a1.5 1.5 0 00-1.5-1.5H3.75A1.5 1.5 0 002.25 6v12a1.5 1.5 0 001.5 1.5zm10.5-11.25h.008v.008h-.008V8.25zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z" />
          </svg>
        </div>
      {/if}
    </div>

    <!-- Content -->
    <div class="flex-1 min-w-0 text-left flex flex-col justify-between">
      <div>
        <div class="flex items-start justify-between gap-3 mb-2">
          <div class="flex-1 min-w-0">
            <h3 class="text-lg font-semibold text-zinc-100 mb-1 line-clamp-2 text-left">
              {patch.title}
            </h3>
            {#if patch.systemGameTitle}
              <p class="text-sm text-zinc-400 mb-2 truncate text-left">
                {patch.systemGameTitle}
              </p>
            {/if}
          </div>
          {#if index < 5}
            <span class="shrink-0 bg-emerald-500/20 text-emerald-400 text-xs font-semibold px-2.5 py-1 border border-emerald-500/30">
              NEW
            </span>
          {/if}
        </div>
      </div>
      <div class="flex items-center justify-between gap-3 mt-auto">
        <div class="flex items-center gap-3">
          {#if patch.rjCode}
            <div class="flex items-center gap-2 text-xs text-zinc-500 font-mono">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-3.5 h-3.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9.568 3H5.25A2.25 2.25 0 003 5.25v4.318c0 .597.237 1.17.659 1.591l9.581 9.581c.699.699 1.78.872 2.607.33a18.095 18.095 0 005.223-5.223c.542-.827.369-1.908-.33-2.607L11.16 3.66A2.25 2.25 0 009.568 3z" />
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 6h.008v.008H6V6z" />
              </svg>
              <span>{patch.rjCode}</span>
            </div>
          {/if}
          {#if patch.storeLink}
            <button
              onclick={() => BrowserOpenURL(patch.storeLink)}
              class="bg-blue-600 hover:bg-blue-700 text-white text-xs font-medium px-3 py-1.5 flex items-center gap-1.5 transition-colors"
            >
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-3.5 h-3.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 6H5.25A2.25 2.25 0 003 8.25v10.5A2.25 2.25 0 005.25 21h10.5A2.25 2.25 0 0018 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" />
              </svg>
              <span>View on DLsite</span>
            </button>
          {/if}
        </div>
        {#if patch.releaseDate}
          <p class="text-xs text-zinc-500 text-right">
            {patch.releaseDate}
          </p>
        {/if}
      </div>
    </div>
  </div>
</div>

