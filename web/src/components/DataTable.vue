<template>
  <div class="table-responsive-wrapper w-full overflow-x-auto rounded-xl border border-slate-800 bg-slate-900/40 backdrop-blur-md">
    <DataTable
      :data="data"
      :columns="columns"
      :options="tableOptions"
      class="display w-full border-collapse"
      ref="tableRef"
    >
      <slot></slot>
    </DataTable>
  </div>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import DataTable from 'datatables.net-vue3'
import DataTablesLib from 'datatables.net-dt'

DataTable.use(DataTablesLib)

const props = defineProps({
  columns: {
    type: Array,
    required: true
  },
  ajaxUrl: {
    type: String,
    required: false
  },
  data: {
    type: Array,
    required: false,
    default: () => []
  },
  serverSide: {
    type: Boolean,
    default: true
  },
  extraParams: {
    type: Object,
    default: () => ({})
  }
})

const tableRef = ref(null)

const tableOptions = computed(() => {
  const options = {
    responsive: true,
    autoWidth: false,
    paging: true,
    pageLength: 10,
    lengthMenu: [5, 10, 25, 50],
    language: {
      search: 'Cari:',
      lengthMenu: 'Tampilkan _MENU_ data',
      info: 'Menampilkan _START_ s/d _END_ dari _TOTAL_ data',
      infoEmpty: 'Menampilkan 0 data',
      infoFiltered: '(disaring dari _MAX_ total data)',
      emptyTable: 'Tidak ada data yang tersedia di tabel',
      zeroRecords: 'Tidak ditemukan data yang sesuai',
      paginate: {
        first: 'Pertama',
        previous: 'Sebelumnya',
        next: 'Selanjutnya',
        last: 'Terakhir'
      }
    }
  }

  if (props.serverSide && props.ajaxUrl) {
    options.serverSide = true
    options.processing = true
    options.ajax = (dataTablesParams, callback) => {
      // Build query string matching Go DataTables parameters
      const params = new URLSearchParams()
      params.append('draw', dataTablesParams.draw)
      params.append('start', dataTablesParams.start)
      params.append('length', dataTablesParams.length)
      params.append('search[value]', dataTablesParams.search?.value || '')

      // Ordering
      if (dataTablesParams.order && dataTablesParams.order.length > 0) {
        params.append('order[0][column]', dataTablesParams.order[0].column)
        params.append('order[0][dir]', dataTablesParams.order[0].dir)
      }

      // Add extra custom parameters
      for (const [key, value] of Object.entries(props.extraParams)) {
        if (value !== undefined && value !== null) {
          params.append(key, value)
        }
      }

      const token = localStorage.getItem('sipecut_token')
      fetch(`${props.ajaxUrl}?${params.toString()}`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      })
        .then(response => response.json())
        .then(res => {
          callback({
            draw: res.draw,
            recordsTotal: res.recordsTotal,
            recordsFiltered: res.recordsFiltered,
            data: res.data || []
          })
        })
        .catch(err => {
          console.error('DataTables server-side request failed:', err)
          callback({
            draw: dataTablesParams.draw,
            recordsTotal: 0,
            recordsFiltered: 0,
            data: []
          })
        })
    }
  }

  return options
})

// Expose table reload method to parent component
const reload = () => {
  if (tableRef.value && tableRef.value.dt) {
    tableRef.value.dt.ajax.reload()
  }
}

defineExpose({
  reload
})
</script>

<style>
/* Scoped override for nested elements */
.dt-processing {
  background: rgba(139, 92, 246, 0.15) !important;
  color: #8b5cf6 !important;
  font-weight: 600 !important;
  border-radius: 0.375rem !important;
}
</style>
