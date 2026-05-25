<script lang="ts">
  import TekojarList from "../tekojar/components/TekojarList.svelte";
  import TekojarLogs from "../tekojar/components/TekojarLogs.svelte";
  import TekojarControls from "../tekojar/components/TekojarControls.svelte";

  let selected = $state("user-service");
</script>

<div class="flex h-full">
  <TekojarList {selected} />

  <div class="flex flex-col flex-1 p-4 gap-3">
    <TekojarControls {selected} />
    <TekojarLogs />
  </div>
</div>
