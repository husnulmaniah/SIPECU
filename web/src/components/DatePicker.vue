<template>
  <div class="relative">
    <input
      :value="displayValue"
      @input="onInput"
      @focus="onFocus"
      @blur="onBlur"
      type="text"
      :placeholder="placeholder || 'DD/MM/YYYY'"
      maxlength="10"
      autocomplete="off"
      :class="inputClass"
    />
    <!-- Calendar Dropdown -->
    <div
      v-if="showPicker"
      class="absolute z-50 mt-1 rounded-xl shadow-2xl p-3 w-72"
      style="background: #FFFFFF; border: 1.5px solid #CBD5E0; box-shadow: 0 8px 32px rgba(26,54,93,0.15);"
      @mousedown.prevent
    >
      <!-- Header: Month/Year Navigation -->
      <div class="flex items-center justify-between mb-3">
        <button type="button" @click="prevMonth"
                class="p-1 rounded-lg transition-colors"
                style="color: #718096;"
                @mouseover="e => e.target.style.background='#EBF0F8'"
                @mouseout="e => e.target.style.background='transparent'">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/></svg>
        </button>
        <div class="flex gap-2">
          <select v-model="viewMonth" @change="onViewChange"
                  class="text-xs rounded-lg px-2 py-1 outline-none cursor-pointer font-semibold"
                  style="background: #EBF0F8; border: 1px solid #CBD5E0; color: #1A365D;">
            <option v-for="(m, i) in monthNames" :key="i" :value="i">{{ m }}</option>
          </select>
          <select v-model="viewYear" @change="onViewChange"
                  class="text-xs rounded-lg px-2 py-1 outline-none cursor-pointer font-semibold"
                  style="background: #EBF0F8; border: 1px solid #CBD5E0; color: #1A365D;">
            <option v-for="y in yearRange" :key="y" :value="y">{{ y }}</option>
          </select>
        </div>
        <button type="button" @click="nextMonth"
                class="p-1 rounded-lg transition-colors"
                style="color: #718096;"
                @mouseover="e => e.target.style.background='#EBF0F8'"
                @mouseout="e => e.target.style.background='transparent'">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/></svg>
        </button>
      </div>

      <!-- Day Names -->
      <div class="grid grid-cols-7 mb-1">
        <span v-for="d in ['Min','Sen','Sel','Rab','Kam','Jum','Sab']" :key="d"
              class="text-center text-[10px] font-bold py-1"
              style="color: #718096;">{{ d }}</span>
      </div>

      <!-- Calendar Days -->
      <div class="grid grid-cols-7 gap-0.5">
        <span v-for="blank in firstDayOffset" :key="'b'+blank"></span>
        <button
          v-for="day in daysInMonth"
          :key="day"
          type="button"
          @click="selectDay(day)"
          class="text-center text-xs py-1.5 rounded-lg transition-all font-medium"
          :style="isSelected(day)
            ? 'background: #1A365D; color: #FFFFFF; box-shadow: 0 2px 8px rgba(26,54,93,0.3);'
            : isToday(day)
              ? 'background: #FEF3C7; color: #92400E; border: 1.5px solid #D69E2E; font-weight: 700;'
              : 'color: #2D3748;'"
          @mouseover="e => { if (!isSelected(day) && !isToday(day)) e.target.style.background='#EBF0F8' }"
          @mouseout="e => { if (!isSelected(day) && !isToday(day)) e.target.style.background='' }"
        >
          {{ day }}
        </button>
      </div>

      <!-- Footer: Clear & Today -->
      <div class="flex justify-between mt-3 pt-2" style="border-top: 1px solid #E2E8F0;">
        <button type="button" @click="clearDate"
                class="text-[10px] font-semibold transition-colors"
                style="color: #A0AEC0;"
                @mouseover="e => e.target.style.color='#E53E3E'"
                @mouseout="e => e.target.style.color='#A0AEC0'">
          Hapus Tanggal
        </button>
        <button type="button" @click="goToday"
                class="text-[10px] font-semibold transition-colors"
                style="color: #D69E2E;"
                @mouseover="e => e.target.style.color='#B7851F'"
                @mouseout="e => e.target.style.color='#D69E2E'">
          Hari Ini
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  modelValue: { type: String, default: '' },      // Internal format: YYYY-MM-DD
  placeholder: { type: String, default: 'DD/MM/YYYY' },
  inputClass: { type: String, default: '' }
})
const emit = defineEmits(['update:modelValue'])

const showPicker = ref(false)
const today = new Date()

// View state: the month/year currently displayed
const viewMonth = ref(today.getMonth())
const viewYear = ref(today.getFullYear())

// Initialize view to modelValue's month/year if available
const syncView = (val) => {
  if (val && !val.startsWith('0001') && val.length >= 7) {
    const [y, m] = val.split('-')
    const yr = parseInt(y), mo = parseInt(m) - 1
    if (!isNaN(yr) && yr > 1 && !isNaN(mo)) {
      viewYear.value = yr
      viewMonth.value = mo
    } else {
      viewYear.value = today.getFullYear()
      viewMonth.value = today.getMonth()
    }
  } else {
    viewYear.value = today.getFullYear()
    viewMonth.value = today.getMonth()
  }
}

watch(() => props.modelValue, (val) => {
  syncView(val)
}, { immediate: true })

// Display value: convert YYYY-MM-DD -> DD/MM/YYYY
const displayValue = computed(() => {
  const v = props.modelValue
  if (!v || v.startsWith('0001') || v.length < 10) return ''
  const [year, month, day] = v.split('-')
  if (!year || !month || !day) return ''
  const yr = parseInt(year)
  if (isNaN(yr) || yr <= 1) return ''
  return `${day}/${month}/${year}`
})

// Parse typed text DD/MM/YYYY -> YYYY-MM-DD
const onInput = (e) => {
  let raw = e.target.value.replace(/[^0-9]/g, '')
  // Auto-insert slashes
  let formatted = ''
  if (raw.length <= 2) formatted = raw
  else if (raw.length <= 4) formatted = raw.slice(0,2) + '/' + raw.slice(2)
  else formatted = raw.slice(0,2) + '/' + raw.slice(2,4) + '/' + raw.slice(4,8)
  e.target.value = formatted

  // When fully entered (DD/MM/YYYY = 10 chars)
  if (formatted.length === 10) {
    const [dd, mm, yyyy] = formatted.split('/')
    const isoDate = `${yyyy}-${mm}-${dd}`
    const d = new Date(isoDate)
    if (!isNaN(d.getTime())) {
      emit('update:modelValue', isoDate)
      syncView(isoDate)
    }
  } else if (formatted.length === 0) {
    emit('update:modelValue', '')
  }
}

const onFocus = () => {
  syncView(props.modelValue)
  showPicker.value = true
}
const onBlur = () => {
  setTimeout(() => { showPicker.value = false }, 150)
}

// Calendar computation
const monthNames = ['Januari','Februari','Maret','April','Mei','Juni','Juli','Agustus','September','Oktober','November','Desember']
const yearRange = computed(() => {
  const years = []
  for (let y = 1950; y <= today.getFullYear() + 10; y++) years.push(y)
  return years
})

const daysInMonth = computed(() => new Date(viewYear.value, viewMonth.value + 1, 0).getDate())
const firstDayOffset = computed(() => new Date(viewYear.value, viewMonth.value, 1).getDay())

const prevMonth = () => {
  if (viewMonth.value === 0) { viewMonth.value = 11; viewYear.value-- }
  else viewMonth.value--
}
const nextMonth = () => {
  if (viewMonth.value === 11) { viewMonth.value = 0; viewYear.value++ }
  else viewMonth.value++
}
const onViewChange = () => {} // reactive, no extra action needed

const selectDay = (day) => {
  const mm = String(viewMonth.value + 1).padStart(2, '0')
  const dd = String(day).padStart(2, '0')
  const iso = `${viewYear.value}-${mm}-${dd}`
  emit('update:modelValue', iso)
  showPicker.value = false
}

const clearDate = () => {
  emit('update:modelValue', '')
  showPicker.value = false
}

const goToday = () => {
  viewMonth.value = today.getMonth()
  viewYear.value = today.getFullYear()
}

const isSelected = (day) => {
  const v = props.modelValue
  if (!v || v.startsWith('0001')) return false
  const [yr, mo, d] = v.split('-')
  return parseInt(yr) === viewYear.value && parseInt(mo) - 1 === viewMonth.value && parseInt(d) === day
}

const isToday = (day) => {
  return today.getFullYear() === viewYear.value && today.getMonth() === viewMonth.value && today.getDate() === day
}
</script>
