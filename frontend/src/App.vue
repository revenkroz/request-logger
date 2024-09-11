<script setup lang="ts">
import type { Ref } from 'vue'
import { computed, onMounted, ref } from 'vue'

import LogEntry from '~/components/LogEntry.vue'
import type { Log } from '~/types'
import { listenLogs } from '~/lib/sse'

const items: Ref<Log[]> = ref([])
const currentItem: Ref<Log|null> = ref(null)
const selectedItems = computed(() => items.value.filter(item => item.selected))

onMounted(() => {
  listenLogs((log: Log) => {
    items.value = [log, ...items.value]
  })
})

function clearLogs() {
  items.value = []
  currentItem.value = null
}

function openLog(log: Log) {
  currentItem.value = log
}

function copyToClipboard(text: string) {
  navigator.clipboard.writeText(text)
}
</script>

<template>
  <div class="w-full overflow-hidden h-screen">
    <div class="w-full h-full overflow-y-auto">
      <div class="grid grid-cols-[120px_1fr] md:grid-cols-[300px_1fr] gap-4 h-full">
        <div
          class="h-full w-full overflow-y-auto p-4"
        >
          <div class="top-bar"></div>
          <div class="mb-2">
            <button
              class="btn-1"
              @click="() => clearLogs()"
              v-text="`Clear Logs`"
            />
          </div>
          <LogEntry
            v-for="log in items"
            :key="log.req"
            :log="log"
            :selectable="true"
            class="cursor-pointer my-1 w-full"
            :class="{
              'bg-slate-300': currentItem === log,
            }"
            @click="openLog(log)"
            @select="log.selected = !log.selected"
          />
          <div class="mt-2 flex flex-col">
            <span
              class="text-secondary"
              v-text="'Total: ' + items.length"
            />
            <span
              class="text-secondary"
              v-text="'Selected: ' + selectedItems.length"
            />
          </div>
        </div>

        <div class="p-4">
          <div class="w-full top-bar bg-gray-100 rounded-md p-4">
            <LogEntry
              v-if="currentItem"
              :key="currentItem.req"
              :log="currentItem"
              :full="true"
              class="w-full"
            />
          </div>

          <div class="grid grid-cols-[50%_50%]">
            <div
              class="h-full overflow-hidden min-w-80vw md:min-w-300px mr-2"
              v-if="currentItem"
            >
              <div>
                <button
                  class="btn-1"
                  @click="copyToClipboard(currentItem?.req || '')"
                  v-text="`Copy`"
                />
              </div>
              <div class="mt-2">
                <span class="text-secondary" v-text="'Request'" />
              </div>
              <pre
                class="code-block"
                v-text="currentItem.req"
              />
            </div>

            <div
              class="h-full overflow-hidden min-w-80vw md:min-w-300px ml-2"
              v-if="currentItem"
            >
              <div>
                <button
                  class="btn-1"
                  @click="copyToClipboard(currentItem?.res || '')"
                  v-text="`Copy`"
                />
              </div>
              <div class="mt-2">
                <span class="text-secondary" v-text="'Response'" />
              </div>
              <pre
                class="code-block"
                v-text="currentItem.res"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
