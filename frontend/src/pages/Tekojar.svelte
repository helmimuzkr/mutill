<script lang="ts">
  import TekojarList from "../components/TekojarList.svelte";
  import TekojarLogs from "../components/TekojarLogs.svelte";
  import TekojarControls from "../components/TekojarControls.svelte";

  let selected = $state("user-service");
</script>

<div class="flex h-full">
  <TekojarList {selected} />

  <div class="flex flex-col flex-1 p-4 gap-3">
    <TekojarControls {selected} />
    <TekojarLogs />
  </div>
</div>
