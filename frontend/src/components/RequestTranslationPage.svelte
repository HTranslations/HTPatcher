<script lang="ts">
  import PageHeader from "./PageHeader.svelte";
  import { BrowserOpenURL } from "../../wailsjs/runtime/runtime.js";

  function openDiscord() {
    BrowserOpenURL("https://discord.gg/sKXZDn72cr");
  }

  function openKofi() {
    BrowserOpenURL("https://ko-fi.com/htranslations");
  }
</script>

<PageHeader
  title="Request Translation"
  description="Learn how to request translations for your favorite games"
/>

<div class="max-w-3xl space-y-6 text-left">
  <!-- Discord Section -->
  <div class="bg-zinc-900 border border-zinc-800 p-6">
    <div class="flex items-start gap-4 text-left">
      <div class="flex-shrink-0 w-12 h-12 bg-indigo-600 flex items-center justify-center">
        <svg xmlns="http://www.w3.org/2000/svg" fill="currentColor" viewBox="0 0 24 24" class="w-7 h-7 text-white">
          <path d="M20.317 4.37a19.791 19.791 0 0 0-4.885-1.515.074.074 0 0 0-.079.037c-.21.375-.444.864-.608 1.25a18.27 18.27 0 0 0-5.487 0 12.64 12.64 0 0 0-.617-1.25.077.077 0 0 0-.079-.037A19.736 19.736 0 0 0 3.677 4.37a.07.07 0 0 0-.032.027C.533 9.046-.32 13.58.099 18.057a.082.082 0 0 0 .031.057 19.9 19.9 0 0 0 5.993 3.03.078.078 0 0 0 .084-.028c.462-.63.874-1.295 1.226-1.994a.076.076 0 0 0-.041-.106 13.107 13.107 0 0 1-1.872-.892.077.077 0 0 1-.008-.128 10.2 10.2 0 0 0 .372-.292.074.074 0 0 1 .077-.01c3.928 1.793 8.18 1.793 12.062 0a.074.074 0 0 1 .078.01c.12.098.246.198.373.292a.077.077 0 0 1-.006.127 12.299 12.299 0 0 1-1.873.892.077.077 0 0 0-.041.107c.36.698.772 1.362 1.225 1.993a.076.076 0 0 0 .084.028 19.839 19.839 0 0 0 6.002-3.03.077.077 0 0 0 .032-.054c.5-5.177-.838-9.674-3.549-13.66a.061.061 0 0 0-.031-.03zM8.02 15.33c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.956-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.956 2.418-2.157 2.418zm7.975 0c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.955-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.946 2.418-2.157 2.418z"/>
        </svg>
      </div>
      <div class="flex-1">
        <h3 class="text-lg font-semibold mb-2">Vote on Discord</h3>
        <p class="text-sm text-zinc-400 mb-4">
          Join our Discord community to vote on translation requests. The games with the most votes will be prioritized for translation by our team.
        </p>
         <button
           onclick={openDiscord}
           class="inline-flex items-center justify-center gap-2 bg-indigo-600 hover:bg-indigo-500 text-white px-4 py-2 text-sm font-semibold transition-colors w-48"
         >
           <span>Join Discord</span>
           <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-4 h-4">
             <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 6H5.25A2.25 2.25 0 003 8.25v10.5A2.25 2.25 0 005.25 21h10.5A2.25 2.25 0 0018 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" />
           </svg>
         </button>
      </div>
    </div>
  </div>

  <!-- Ko-fi Section -->
  <div class="bg-zinc-900 border border-zinc-800 p-6">
    <div class="flex items-start gap-4 text-left">
      <div class="flex-shrink-0 w-12 h-12 bg-cyan-600 flex items-center justify-center">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-7 h-7 text-white">
          <path stroke-linecap="round" stroke-linejoin="round" d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12z" />
        </svg>
      </div>
      <div class="flex-1">
        <h3 class="text-lg font-semibold mb-2">Support on Ko-fi</h3>
        <p class="text-sm text-zinc-400 mb-4">
          Support us on Ko-fi to request a specific game translation. Your contribution helps us create high-quality translations and maintain this tool.
        </p>
         <button
           onclick={openKofi}
           class="inline-flex items-center justify-center gap-2 bg-cyan-600 hover:bg-cyan-500 text-white px-4 py-2 text-sm font-semibold transition-colors w-48"
         >
           <span>Support on Ko-fi</span>
           <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-4 h-4">
             <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 6H5.25A2.25 2.25 0 003 8.25v10.5A2.25 2.25 0 005.25 21h10.5A2.25 2.25 0 0018 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" />
           </svg>
         </button>
      </div>
    </div>
  </div>
</div>

