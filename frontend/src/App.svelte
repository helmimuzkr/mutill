<script lang="ts">
  import { ModeWatcher } from "mode-watcher";
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import Braces from "@lucide/svelte/icons/braces";
  import Monitor from "@lucide/svelte/icons/monitor";
  import SettingsIcon from "@lucide/svelte/icons/settings";

  import Tekojar from "./pages/Tekojar.svelte";
  import JsonataQuery from "./pages/JsonataQuery.svelte";

  let page = $state("tekojar");

  const items = [
    {
      title: "Tekojar",
      onclick: () => (page = "tekojar"),
      icon: Monitor,
    },
    {
      title: "JSONata",
      onclick: () => (page = "jsonata"),

      icon: Braces,
    },
    {
      title: "Settings",
      url: "#",
      icon: SettingsIcon,
    },
  ];
</script>

<!-- Dark Mode Whatcher -->
<ModeWatcher />

<div class="flex h-screen">
  <Sidebar.Provider>
    <Sidebar.Sidebar>
      <Sidebar.Content>
        <Sidebar.Group>
          <Sidebar.GroupLabel>Menu</Sidebar.GroupLabel>

          <Sidebar.GroupContent>
            <Sidebar.Menu>
              {#each items as item (item.title)}
                <Sidebar.MenuItem>
                  <Sidebar.MenuButton>
                    {#snippet child({ props })}
                      <button onclick={item.onclick} {...props}>
                        <item.icon />
                        <span>{item.title}</span>
                      </button>
                    {/snippet}
                  </Sidebar.MenuButton>
                </Sidebar.MenuItem>
              {/each}
            </Sidebar.Menu>
          </Sidebar.GroupContent>
        </Sidebar.Group>
      </Sidebar.Content>
    </Sidebar.Sidebar>

    <Sidebar.Inset>
      <div class="flex-1">
        {#if page === "tekojar"}
          <Tekojar />
        {:else}
          <JsonataQuery />
        {/if}
      </div>
    </Sidebar.Inset>
  </Sidebar.Provider>
</div>
