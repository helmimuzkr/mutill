<script lang="ts">
  import type { TekojarSetting } from './type';
  import * as Field from '$lib/components/ui/field/index.js';
  import { Input } from '$lib/components/ui/input/index.js';
  import Button from '$lib/components/ui/button/button.svelte';
  import { Checkbox } from '$lib/components/ui/checkbox/index.js';
  import * as Table from '$lib/components/ui/table';
  import { onMount } from 'svelte';
  import { tekojar } from '../../wailsjs/go/models';
  import { GetSetting, SaveSetting } from '../../wailsjs/go/backend/TekojarApp';
  import Skeleton from '$lib/components/ui/skeleton/skeleton.svelte';

  let error = $state<string | null>(null);
  let tekojarSetting = $state<TekojarSetting>();

  onMount(async () => {
    try {
      const result = await GetSetting();
      tekojarSetting = result;

      $inspect(tekojarSetting);
    } catch (err) {
      console.log(err instanceof Error ? err.message : 'Failed to fetch services');
    }
  });

  function addService() {
    error = null;
    tekojarSetting.service_settings = [
      ...tekojarSetting.service_settings,
      { name: '', path: '', skip_flag: false, delay: 0 },
    ];
  }

  function removeService(name: string) {
    tekojarSetting.service_settings = tekojarSetting.service_settings.filter((s) => s.name !== name);
  }

  async function save() {
    // call Go backend later
    tekojarSetting.service_settings = tekojarSetting.service_settings.filter((s) => s.name);
    await SaveSetting(tekojarSetting);
    console.log(tekojarSetting);
  }

  function handleServiceNameChange(service: tekojar.ServiceSetting, name: string) {
    service.name = name;
    error = null;

    validateDuplicateService(name);
  }

  function validateDuplicateService(name: string) {
    const isDuplicate = tekojarSetting.service_settings.filter((s) => s.name === name).length > 1;
    error = isDuplicate ? `"${name}" already exists` : null;
  }
</script>

<div class="w-full max-w-2xl">
  {#if tekojarSetting}
    <form>
      <Field.Group>
        <Field.Set>
          <Field.Legend>General</Field.Legend>
          <Field.Group>
            <Field.Field orientation="horizontal">
              <Field.Label class="text-sm w-24 font-normal">Command</Field.Label>
              <Input bind:value={tekojarSetting.command} />
            </Field.Field>

            <Field.Field orientation="horizontal">
              <Checkbox bind:checked={tekojarSetting.auto_shutdown} />
              <Field.Label class="text-sm w-24 font-normal">Auto Shutdown</Field.Label>
            </Field.Field>
          </Field.Group>
        </Field.Set>

        <Field.Separator />

        <Field.Set>
          <Field.Legend>Service</Field.Legend>
          <div class="flex flex-col gap-3">
            <Table.Root>
              <Table.Header>
                <Table.Row>
                  <Table.Head>Name</Table.Head>
                  <Table.Head>Path</Table.Head>
                  <Table.Head>Skip</Table.Head>
                  <Table.Head>Delay (s)</Table.Head>
                  <Table.Head></Table.Head>
                </Table.Row>
              </Table.Header>
              <Table.Body>
                {#each tekojarSetting?.service_settings as service, i (i)}
                  <Table.Row>
                    <Table.Cell>
                      <Input
                        value={service.name}
                        placeholder="service.jar"
                        oninput={(e) => handleServiceNameChange(service, e.currentTarget.value)}
                      />
                    </Table.Cell>
                    <Table.Cell>
                      <Input bind:value={service.path} placeholder="/home/user/service" />
                    </Table.Cell>
                    <Table.Cell>
                      <Checkbox bind:checked={service.skip_flag} />
                    </Table.Cell>
                    <Table.Cell>
                      <Input type="number" bind:value={service.delay} class="w-20" />
                    </Table.Cell>
                    <Table.Cell>
                      <Button variant="destructive" size="sm" onclick={() => removeService(service.name)}>Remove</Button
                      >
                    </Table.Cell>
                  </Table.Row>
                {/each}
              </Table.Body>
            </Table.Root>

            {#if error}
              <p class="text-destructive text-sm">{error}</p>
            {/if}

            <Button variant="outline" size="sm" class="w-fit" onclick={addService}>Add Service</Button>
          </div>
        </Field.Set>

        <Field.Field orientation="horizontal" class="justify-end">
          <Button onclick={save}>Save</Button>
        </Field.Field>
      </Field.Group>
    </form>
  {:else}
    <Skeleton />
  {/if}
</div>
