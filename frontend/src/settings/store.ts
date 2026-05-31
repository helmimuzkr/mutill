import { derived, writable } from "svelte/store";

export const selectedSettingName = writable<string | null>(null);
export const selectedChildSetting = writable<string | null>(null)



