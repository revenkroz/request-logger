<script setup lang="ts">
import type { Log } from '~/types'

const emit = defineEmits<{
  (e: 'select', value: Log): void
}>()

const props = defineProps<{
  log: Log
  full: boolean
  selectable: boolean
}>()

const time = !props.full ?
  `${props.log.elapsed}ms` :
  `${props.log.elapsed}ms / ${props.log.done_at}`

const status = parseInt(props.log.status.split(' ')[0])
</script>

<template>
  <div
    class="flex flex-row items-center p-2 bg-gray-100 color-black rounded-md"
  >
    <div
      v-if="selectable"
      class="mr-4"
    >
      <div
        class="rounded-full w-4 h-4 cursor-pointer"
        :class="{
          'bg-blue-500': log.selected,
          'bg-gray-400': !log.selected,
        }"
        @click="() => emit('select', log)"
      />
    </div>
    <div class="flex flex-col w-full">
      <div class="grid grid-cols-[auto_1fr] gap-2">
        <div
          :class="`method-${log.method.toLowerCase()}`"
          v-text="log.method"
        />
        <div
          class="whitespace-nowrap"
          v-text="full ? log.full_url : log.url.path"
        />
      </div>
      <div class="flex justify-between items-center">
        <span
          class="text-sm color-gray-500"
          v-text="time"
        />
        <span
          class="text-sm"
          :class="{
            'color-green-500': status >= 200 && status < 300,
            'color-yellow-500': status >= 300 && status < 400,
            'color-red-500': status >= 400,
          }"
          v-text="full ? log.status : status"
        />
      </div>
    </div>
  </div>
</template>
