<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-xl font-bold text-white">Arsip Pegawai Pensiun</h1>
      <p class="text-xs text-slate-400 mt-1">Daftar seluruh pegawai berstatus Pensiun. Akun login mereka dinonaktifkan, namun histori arsip data tetap disimpan.</p>
    </div>

    <!-- DataTable container -->
    <div class="glass-panel p-4 rounded-2xl" @click="handleTableClick">
      <DataTable
        ref="tableRef"
        :columns="tableColumns"
        ajaxUrl="/api/employees/retired"
      />
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import DataTable from '@/components/DataTable.vue'

const router = useRouter()
const tableRef = ref(null)

const tableColumns = [
  { data: 'nip', title: 'NIP' },
  { data: 'nama', title: 'Nama Pegawai' },
  { data: 'jenis_jabatan', title: 'Jenis Jabatan' },
  { data: 'jabatan', title: 'Jabatan' },
  { data: 'tempat_tugas', title: 'Tempat Tugas terakhir' },
  {
    data: 'umur',
    title: 'Umur saat ini',
    orderable: false
  },
  {
    data: null,
    title: 'Aksi',
    orderable: false,
    render: (data, type, row) => {
      return `
        <button data-nip="${row.nip}" class="btn-detail px-3 py-1.5 bg-primary-600 hover:bg-primary-700 text-white font-bold rounded-lg text-xs transition-colors">
          Detail Arsip
        </button>
      `
    }
  }
]

// Delegate row action click to Vue Router push (avoiding full page reload)
const handleTableClick = (e) => {
  const target = e.target
  if (target.classList.contains('btn-detail')) {
    const nip = target.getAttribute('data-nip')
    router.push(`/pegawai/${nip}`)
  }
}
</script>
