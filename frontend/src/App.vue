<script setup lang="ts">
import { computed, Ref } from 'vue'
import { onMounted, ref } from 'vue'

import type { Log } from '~/types'
import LogEntry from '~/components/LogEntry.vue'

const items: Ref<Log[]> = ref([])
const currentItem: Ref<Log|null> = ref(null)
const selectedItems = computed(() => items.value.filter(item => item.selected))

onMounted(() => {
  startSse()
})

function startSse() {
  let es = new EventSource('/logs')

  es.addEventListener('log', event => {
    items.value = [JSON.parse(event.data), ...items.value]
  })
}

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
          class="h-full overflow-y-auto w-full p-4"
        >
          <div class="h-80px mb-4"></div>
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
            :full="false"
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
              class="text-sm color-gray-500"
              v-text="'Total: ' + items.length"
            />
            <span
              class="text-sm color-gray-500"
              v-text="'Selected: ' + selectedItems.length"
            />
          </div>
        </div>

        <div class="p-4">
          <div class="w-full h-80px bg-gray-100 rounded-md p-4 mb-4">
            <LogEntry
              v-if="currentItem"
              :key="currentItem.req"
              :log="currentItem"
              :full="true"
              :selectable="false"
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
                <span
                  class="text-sm color-gray-500"
                  v-text="'Request'"
                />
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
                <span
                  class="text-sm color-gray-500"
                  v-text="'Response'"
                />
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
