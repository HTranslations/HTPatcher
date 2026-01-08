<script lang="ts">
  import { BrowserOpenURL } from "../../wailsjs/runtime/runtime.js";

  function openWebsite() {
    BrowserOpenURL("https://htranslations.com");
  }

  function openDiscord() {
    BrowserOpenURL("https://discord.gg/sKXZDn72cr");
  }

  function openKofi() {
    BrowserOpenURL("https://ko-fi.com/htranslations");
  }
</script>

<div class="bg-zinc-900 border-b border-zinc-800 px-6 py-4">
  <div class="flex items-center justify-between">
    <div class="flex items-center gap-4">
      <h1 class="text-xl font-semibold">HTranslations Patcher</h1>
    </div>
    <div class="flex items-center gap-4">
      <button onclick={openWebsite} class="text-sm text-zinc-400 hover:text-zinc-200 transition-colors">
        Website
      </button>
      <button onclick={openDiscord} class="text-sm text-zinc-400 hover:text-zinc-200 transition-colors">
        Discord
      </button>
      <button onclick={openKofi} class="text-sm text-emerald-400 font-medium hover:text-emerald-300 transition-colors">
        Support Us
      </button>
    </div>
  </div>
</div>

