<script lang="ts">
  import { onMount } from "svelte";
  import { domain } from "../../wailsjs/go/models.js";
  import PageHeader from "./PageHeader.svelte";
  import {
    CheckForUpdate,
    GetCurrentVersion,
    GetLatestReleaseInfo,
  } from "../../wailsjs/go/main/App.js";
  import { BrowserOpenURL } from "../../wailsjs/runtime/runtime.js";

  export let dataPath: string;
  export let onDeleteData: () => void;

  let currentVersion = 0;
  let updateReleaseInfo: domain.ReleaseInfo | null = null;
  let latestReleaseInfo: domain.ReleaseInfo | null = null;
  let checkingUpdate = false;

  async function loadVersion() {
    try {
      currentVersion = await GetCurrentVersion();
    } catch (error) {
      console.error("Failed to get current version:", error);
    }
  }

  async function loadLatestReleaseInfo() {
    try {
      latestReleaseInfo = await GetLatestReleaseInfo();
    } catch (error) {
      console.error("Failed to get latest release info:", error);
    }
  }

  async function checkForUpdate() {
    checkingUpdate = true;
    try {
      const releaseInfo = await CheckForUpdate();
      updateReleaseInfo = releaseInfo;
      // Always load latest release info to show the latest version
      await loadLatestReleaseInfo();
    } catch (error) {
      console.error("Failed to check for update:", error);
    } finally {
      checkingUpdate = false;
    }
  }

  function handleUpdate() {
    if (updateReleaseInfo) {
      const exeAsset = updateReleaseInfo.assets?.find((asset) =>
        asset.name.endsWith(".exe")
      );
      if (exeAsset?.browser_download_url) {
        BrowserOpenURL(exeAsset.browser_download_url);
      } else if (updateReleaseInfo.html_url) {
        BrowserOpenURL(updateReleaseInfo.html_url);
      }
    }
  }

  function formatFileSize(bytes: number): string {
    if (bytes === 0) return "0 Bytes";
    const k = 1024;
    const sizes = ["Bytes", "KB", "MB", "GB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + " " + sizes[i];
  }

  onMount(async () => {
    await loadVersion();
    await loadLatestReleaseInfo();
    await checkForUpdate();
  });
</script>

<PageHeader
  title="Settings"
  description="Manage application settings and preferences"
/>

<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
  <!-- Version Section -->
  <div class="bg-zinc-900 border border-zinc-800 p-6 flex flex-col">
    <h3 class="text-lg font-semibold text-zinc-100 mb-6">Version Information</h3>
    <div class="flex-1 space-y-4">
      <div class="flex items-center justify-between">
        <span class="text-sm text-zinc-400">Current Version</span>
        <span class="text-sm font-medium text-zinc-300">v{currentVersion}</span>
      </div>
      {#if latestReleaseInfo}
        <div class="flex items-center justify-between pt-4 border-t border-zinc-800">
          <span class="text-sm text-zinc-400">Latest Version</span>
          <span class="text-sm font-medium text-zinc-300">
            {latestReleaseInfo.tag_name}
          </span>
        </div>
      {/if}
    </div>
    {#if updateReleaseInfo}
      <div class="pt-6 mt-4 border-t border-zinc-800 space-y-4">
        <div>
          <span class="text-sm text-zinc-400">Update Available</span>
          <p class="text-xs text-zinc-500 mt-1">
            {updateReleaseInfo.tag_name} is available
          </p>
        </div>
        {#if updateReleaseInfo.assets}
          {@const exeAsset = updateReleaseInfo.assets.find((asset) =>
            asset.name.endsWith(".exe")
          )}
          {#if exeAsset}
            <div class="text-xs text-zinc-500">
              File size: {formatFileSize(exeAsset.size)}
            </div>
          {/if}
        {/if}
        <button
          onclick={handleUpdate}
          class="w-full px-4 py-2 text-sm bg-blue-600 border border-blue-600 text-white hover:bg-blue-700 transition-colors"
        >
          Update Now
        </button>
      </div>
    {:else}
      <div class="pt-6 mt-6 border-t border-zinc-800 space-y-4">
        <div class="flex items-center justify-between">
          <span class="text-sm text-zinc-400">Update Status</span>
          {#if checkingUpdate}
            <span class="text-xs text-zinc-500">Checking...</span>
          {:else}
            <span class="text-xs text-emerald-500">Up to date</span>
          {/if}
        </div>
        <button
          onclick={checkForUpdate}
          disabled={checkingUpdate}
          class="w-full px-4 py-2 text-sm bg-zinc-800 border border-zinc-700 text-zinc-300 hover:bg-zinc-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Check for Updates
        </button>
      </div>
    {/if}
  </div>

  <!-- Data Location Section -->
  <div class="bg-zinc-900 border border-zinc-800 p-6 flex flex-col">
    <h3 class="text-lg font-semibold text-zinc-100 mb-4">Data Location</h3>
    <div class="space-y-3 flex-1">
      <div>
        <p class="text-sm text-zinc-400 mb-2">Persistent data file location:</p>
        <div class="bg-zinc-800/50 border border-zinc-700 p-3 rounded">
          <p class="text-xs text-zinc-300 font-mono break-all">{dataPath}</p>
        </div>
      </div>
    </div>
    <div class="pt-4 border-t border-zinc-800 mt-auto">
      <p class="text-sm text-zinc-400 mb-3">
        Delete all persistent data. This will remove all games from your
        collection and cannot be undone.
      </p>
      <button
        onclick={onDeleteData}
        class="w-full px-4 py-2 text-sm bg-red-600 border border-red-600 text-white hover:bg-red-700 transition-colors"
      >
        Delete All Data
      </button>
    </div>
  </div>
</div>
