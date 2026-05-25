
import type { Component } from "svelte";

export type Page = {
  id: string;
  title: string;
  icon: Component;
  component: Component | null;
};
