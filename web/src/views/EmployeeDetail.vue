<template>
  <div class="space-y-6">
    <!-- Header profile banner -->
    <div class="glass-panel p-6 rounded-2xl flex flex-col md:flex-row items-center gap-6 relative overflow-hidden">
      <!-- Glow effect -->
      <div class="absolute w-64 h-64 rounded-full bg-primary-600/10 blur-[60px] -left-20 -top-20"></div>

      <!-- Foto Profil -->
      <div class="w-24 h-24 rounded-full bg-slate-800 border-2 border-primary-500/50 flex items-center justify-center font-bold text-3xl text-primary-400 overflow-hidden relative z-10 shrink-0 shadow-lg">
        <img v-if="employee.foto_profil" :src="employee.foto_profil" class="w-full h-full object-cover" />
        <span v-else>{{ employee.nama?.charAt(0) }}</span>
      </div>

      <div class="flex-1 text-center md:text-left relative z-10 space-y-1">
        <h1 class="text-xl font-bold text-white">{{ employee.nama }}</h1>
        <p class="text-xs text-slate-400 font-medium">{{ employee.jabatan }} &bull; {{ employee.tempat_tugas }}</p>
        <div class="flex flex-wrap justify-center md:justify-start gap-2 mt-2">
          <span class="px-2 py-0.5 rounded bg-slate-800 border border-slate-700 text-[10px] text-slate-350">NIP: {{ employee.nip }}</span>
          <span class="px-2 py-0.5 rounded text-[10px]"
                :class="employee.status_kepegawaian === 'Aktif' ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20' : 'bg-amber-500/10 text-amber-400 border border-amber-500/20'">
            {{ employee.status_kepegawaian }}
          </span>
          <span class="px-2 py-0.5 rounded bg-primary-500/10 text-primary-400 border border-primary-500/20 text-[10px]">{{ employee.jenis_jabatan }}</span>
        </div>
      </div>

      <!-- Admin Actions -->
      <div v-if="authStore.role === 'admin'" class="flex flex-col sm:flex-row gap-2 relative z-10 shrink-0 w-full md:w-auto">
        <button @click="showEditModal = true" class="px-3 py-2 bg-slate-800 hover:bg-slate-750 text-white rounded-xl text-xs font-bold transition-all border border-slate-700">
          Ubah Profil
        </button>
        <button @click="openKgbModal" class="px-3 py-2 bg-amber-600 hover:bg-amber-700 text-white rounded-xl text-xs font-bold transition-all shadow-md shadow-amber-500/10">
          Update KGB
        </button>
        <button @click="openPangkatModal" class="px-3 py-2 bg-violet-600 hover:bg-violet-700 text-white rounded-xl text-xs font-bold transition-all shadow-md shadow-violet-500/10">
          Update Pangkat
        </button>
      </div>
    </div>

    <!-- Tabs Navigation -->
    <div class="flex border-b border-slate-800 gap-4">
      <button v-for="tab in tabs" :key="tab.id" @click="currentTab = tab.id"
              class="pb-3 text-xs font-semibold uppercase tracking-wider relative transition-all"
              :class="currentTab === tab.id ? 'text-primary-400' : 'text-slate-500 hover:text-slate-300'">
        {{ tab.name }}
        <div v-if="currentTab === tab.id" class="absolute bottom-0 left-0 right-0 h-0.5 bg-primary-500"></div>
      </button>
    </div>

    <!-- TAB: Detail Informasi -->
    <div v-if="currentTab === 'info'" class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div class="glass-panel p-6 rounded-2xl space-y-4">
        <h3 class="text-sm font-bold text-white border-b border-slate-800 pb-2">Biodata Diri</h3>
        <div class="grid grid-cols-3 gap-2 text-xs">
          <span class="text-slate-500 font-semibold">Tempat, Tgl Lahir</span>
          <span class="col-span-2 text-slate-200">{{ employee.tempat_lahir }}, {{ formatDate(employee.tanggal_lahir) }}</span>

          <span class="text-slate-500 font-semibold">Umur Pegawai</span>
          <span class="col-span-2 text-slate-200">{{ ageStr }}</span>

          <span class="text-slate-500 font-semibold">Unit Kerja</span>
          <span class="col-span-2 text-slate-200">{{ employee.tempat_tugas }} ({{ employee.jenis_tempat }})</span>

          <span class="text-slate-500 font-semibold">Pengangkatan</span>
          <span class="col-span-2 text-slate-200">{{ employee.pengangkatan || '-' }}</span>

          <span class="text-slate-500 font-semibold">Jenis Kepegawaian</span>
          <span class="col-span-2 text-slate-200">{{ employee.jenis_pengangkatan }}</span>

          <span class="text-slate-500 font-semibold">No. WhatsApp</span>
          <span class="col-span-2 text-slate-200">{{ employeeUserPhone || '-' }}</span>
        </div>
      </div>

      <div class="glass-panel p-6 rounded-2xl space-y-4">
        <h3 class="text-sm font-bold text-white border-b border-slate-800 pb-2">Jadwal Target & Siklus KGB / Pangkat</h3>
        <div class="grid grid-cols-3 gap-2 text-xs">
          <!-- KGB rows: hanya tampil jika ada tanggal terakhir -->
          <template v-if="hasDate(employee.tanggal_kgb_terakhir)">
            <span class="text-slate-500 font-semibold">KGB Terakhir</span>
            <span class="col-span-2 text-slate-250">{{ formatDate(employee.tanggal_kgb_terakhir) }}</span>

            <span class="text-slate-500 font-semibold">KGB Berikutnya</span>
            <span class="col-span-2 text-amber-400 font-bold">
              {{ hasDate(employee.tanggal_kgb_berikutnya) ? formatDate(employee.tanggal_kgb_berikutnya) : '-' }}
            </span>
          </template>
          <template v-else>
            <span class="text-slate-500 font-semibold">KGB Terakhir</span>
            <span class="col-span-2 text-slate-600 italic">Belum diisi</span>
          </template>

          <!-- Pangkat rows: hanya tampil jika ada tanggal terakhir -->
          <template v-if="hasDate(employee.tanggal_kenaikan_pangkat_terakhir)">
            <span class="text-slate-500 font-semibold">Pangkat Terakhir</span>
            <span class="col-span-2 text-slate-250">{{ formatDate(employee.tanggal_kenaikan_pangkat_terakhir) }}</span>

            <span class="text-slate-500 font-semibold">Pangkat Berikutnya</span>
            <span class="col-span-2 text-violet-400 font-bold">
              {{ hasDate(employee.tanggal_kenaikan_pangkat_berikutnya) ? formatDate(employee.tanggal_kenaikan_pangkat_berikutnya) : '-' }}
            </span>
          </template>
          <template v-else>
            <span class="text-slate-500 font-semibold">Pangkat Terakhir</span>
            <span class="col-span-2 text-slate-600 italic">Belum diisi</span>
          </template>

          <span class="text-slate-500 font-semibold">Tanggal Pensiun</span>
          <span class="col-span-2 text-red-400 font-bold">{{ formatDate(employee.tanggal_pensiun) }}</span>
        </div>
      </div>
    </div>

    <!-- TAB: Berkas Dokumen SK -->
    <div v-if="currentTab === 'docs'" class="glass-panel p-6 rounded-2xl space-y-4">
      <h3 class="text-sm font-bold text-white border-b border-slate-800 pb-2">Berkas Pendukung Terunggah</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <!-- Helper component pattern for each SK item -->
        <div v-for="doc in skDocuments" :key="doc.key"
             class="p-4 bg-slate-900/40 border border-slate-800 rounded-xl flex items-center justify-between gap-3">
          <div class="flex-1 min-w-0">
            <p class="text-xs font-bold text-white">{{ doc.label }}</p>
            <p class="text-[10px] text-slate-500">{{ doc.desc }}</p>
          </div>
          <div class="flex items-center gap-1.5 shrink-0">
            <template v-if="employee[doc.key]">
              <!-- Preview Button -->
              <button @click="openPreview(employee[doc.key], doc.label)"
                      class="px-2.5 py-1 bg-primary-600/20 text-primary-400 hover:bg-primary-600/30 rounded border border-primary-500/30 text-[10px] font-semibold flex items-center gap-1 transition-all">
                <Eye class="w-3 h-3" />
                Preview
              </button>
              <!-- Download Button -->
              <a :href="employee[doc.key]" target="_blank"
                 class="px-2.5 py-1 bg-slate-800 text-slate-300 hover:text-white rounded border border-slate-700 text-[10px] font-semibold">
                Unduh
              </a>
            </template>
            <span v-else class="text-[10px] text-slate-600">Belum diunggah</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Preview Modal (Lightbox) -->
    <div v-if="previewUrl" @click.self="closePreview"
         class="fixed inset-0 bg-black/90 backdrop-blur-sm flex items-center justify-center p-4 z-[100]">
      <div class="relative w-full max-w-4xl max-h-[90vh] flex flex-col">
        <!-- Header -->
        <div class="flex items-center justify-between bg-slate-900 rounded-t-xl px-4 py-2 shrink-0">
          <span class="text-sm font-bold text-white truncate">{{ previewLabel }}</span>
          <div class="flex items-center gap-2">
            <a :href="previewUrl" target="_blank"
               class="px-3 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded text-xs font-semibold flex items-center gap-1">
              <Download class="w-3.5 h-3.5" />Unduh
            </a>
            <button @click="closePreview" class="text-slate-400 hover:text-white p-1 rounded">
              <X class="w-5 h-5" />
            </button>
          </div>
        </div>
        <!-- Preview Content -->
        <div class="flex-1 bg-white rounded-b-xl overflow-hidden min-h-0">
          <!-- PDF -->
          <iframe v-if="previewIsPdf" :src="previewUrl" class="w-full h-full min-h-[75vh]" />
          <!-- Image -->
          <div v-else class="flex items-center justify-center h-full min-h-[75vh] bg-slate-950">
            <img :src="previewUrl" class="max-w-full max-h-[75vh] object-contain" :alt="previewLabel" />
          </div>
        </div>
      </div>
    </div>

    <!-- TAB: Histori Karir KGB & Pangkat -->
    <div v-if="currentTab === 'history'" class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- KGB History -->
      <div class="glass-panel p-6 rounded-2xl space-y-4">
        <h3 class="text-sm font-bold text-amber-400 border-b border-slate-800 pb-2">Riwayat Kenaikan Gaji Berkala (KGB)</h3>
        <div v-if="kgbHistory.length === 0" class="text-xs text-slate-500 py-4 text-center">Belum ada riwayat KGB tercatat.</div>
        <div v-else class="relative border-l border-slate-800 pl-4 space-y-4 py-2">
          <div v-for="item in kgbHistory" :key="item.id" class="relative">
            <!-- marker dot -->
            <div class="absolute -left-[21px] top-1.5 w-3.5 h-3.5 rounded-full bg-amber-500/20 border-2 border-amber-500"></div>
            <p class="text-xs font-bold text-white">{{ formatDate(item.tanggal_kgb) }}</p>
            <div class="flex items-center gap-2 mt-1.5">
              <span class="text-[10px] text-slate-500">Tercatat: {{ formatDateTime(item.created_at) }}</span>
              <a v-if="item.file_sk_kgb" :href="item.file_sk_kgb" target="_blank" class="text-[10px] text-primary-400 hover:underline">Download SK</a>
            </div>
          </div>
        </div>
      </div>

      <!-- Pangkat History -->
      <div class="glass-panel p-6 rounded-2xl space-y-4">
        <h3 class="text-sm font-bold text-violet-400 border-b border-slate-800 pb-2">Riwayat Kenaikan Pangkat</h3>
        <div v-if="pangkatHistory.length === 0" class="text-xs text-slate-500 py-4 text-center">Belum ada riwayat pangkat tercatat.</div>
        <div v-else class="relative border-l border-slate-800 pl-4 space-y-4 py-2">
          <div v-for="item in pangkatHistory" :key="item.id" class="relative">
            <!-- marker dot -->
            <div class="absolute -left-[21px] top-1.5 w-3.5 h-3.5 rounded-full bg-violet-500/20 border-2 border-violet-500"></div>
            <p class="text-xs font-bold text-white">{{ formatDate(item.tanggal_kenaikan_pangkat) }}</p>
            <div class="flex items-center gap-2 mt-1.5">
              <span class="text-[10px] text-slate-500">Tercatat: {{ formatDateTime(item.created_at) }}</span>
              <a v-if="item.file_sk_pangkat" :href="item.file_sk_pangkat" target="_blank" class="text-[10px] text-primary-400 hover:underline">Download SK</a>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- UPDATE KGB MODAL (Admin) -->
    <div v-if="showKgbModal" class="fixed inset-0 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4 z-50">
      <div class="glass-panel w-full max-w-md rounded-2xl p-6 space-y-4">
        <div class="flex items-center justify-between border-b border-slate-800 pb-2">
          <h2 class="text-sm font-bold text-white">Input Kenaikan Gaji Berkala Baru</h2>
          <button @click="showKgbModal = false" class="text-slate-400 hover:text-white"><X class="w-5 h-5" /></button>
        </div>
        <form @submit.prevent="submitKgb" class="space-y-4">
          <div>
            <label class="block text-xs font-semibold text-slate-350 mb-1">Tanggal KGB Baru*</label>
            <DatePicker
              v-model="kgbForm.tanggal_kgb"
              inputClass="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 rounded-lg text-xs text-slate-100 outline-none"
            />
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-350 mb-1">Upload File SK KGB (PDF/Gambar)*</label>
            <input type="file" @change="onHistoryFileChange($event, 'file_sk_kgb')" accept=".pdf,image/*" required class="w-full text-xs text-slate-400" />
          </div>
          <div class="flex justify-end gap-2 pt-2">
            <button type="button" @click="showKgbModal = false" class="px-3 py-1.5 bg-slate-800 hover:bg-slate-750 text-slate-300 rounded-lg text-xs font-semibold">Batal</button>
            <button type="submit" :disabled="modalSubmitting" class="px-3 py-1.5 bg-amber-600 hover:bg-amber-700 text-white rounded-lg text-xs font-semibold flex items-center gap-1">
              <Loader2 v-if="modalSubmitting" class="w-3 animate-spin" />
              <span>Simpan KGB</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- UPDATE PANGKAT MODAL (Admin) -->
    <div v-if="showPangkatModal" class="fixed inset-0 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4 z-50">
      <div class="glass-panel w-full max-w-md rounded-2xl p-6 space-y-4">
        <div class="flex items-center justify-between border-b border-slate-800 pb-2">
          <h2 class="text-sm font-bold text-white">Input Kenaikan Pangkat Baru</h2>
          <button @click="showPangkatModal = false" class="text-slate-400 hover:text-white"><X class="w-5 h-5" /></button>
        </div>
        <form @submit.prevent="submitPangkat" class="space-y-4">
          <div>
            <label class="block text-xs font-semibold text-slate-350 mb-1">Tanggal Kenaikan Pangkat Baru*</label>
            <DatePicker
              v-model="pangkatForm.tanggal_kenaikan_pangkat"
              inputClass="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 rounded-lg text-xs text-slate-100 outline-none"
            />
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-350 mb-1">Upload File SK Kenaikan Pangkat*</label>
            <input type="file" @change="onHistoryFileChange($event, 'file_sk_pangkat')" accept=".pdf,image/*" required class="w-full text-xs text-slate-400" />
          </div>
          <div class="flex justify-end gap-2 pt-2">
            <button type="button" @click="showPangkatModal = false" class="px-3 py-1.5 bg-slate-800 hover:bg-slate-750 text-slate-300 rounded-lg text-xs font-semibold">Batal</button>
            <button type="submit" :disabled="modalSubmitting" class="px-3 py-1.5 bg-violet-600 hover:bg-violet-700 text-white rounded-lg text-xs font-semibold flex items-center gap-1">
              <Loader2 v-if="modalSubmitting" class="w-3 animate-spin" />
              <span>Simpan Pangkat</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- EDIT PROFILE MODAL (Admin) -->
    <div v-if="showEditModal" class="fixed inset-0 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4 z-50 overflow-y-auto">
      <div class="glass-panel w-full max-w-2xl rounded-2xl p-6 space-y-4 max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between border-b border-slate-800 pb-2">
          <h2 class="text-sm font-bold text-white">Ubah Data Pegawai</h2>
          <button @click="showEditModal = false" class="text-slate-400 hover:text-white"><X class="w-5 h-5" /></button>
        </div>
        <form @submit.prevent="submitEditProfile" class="space-y-4">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-semibold text-slate-350 mb-1">Nama Lengkap*</label>
              <input v-model="editForm.nama" type="text" required class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 rounded-lg text-xs text-slate-100 outline-none" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-350 mb-1">No. WhatsApp Pegawai</label>
              <input v-model="editForm.no_hp" type="text" placeholder="Contoh: 08123456789" class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 rounded-lg text-xs text-slate-100 outline-none" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-350 mb-1">Jenis Jabatan*</label>
              <select v-model="editForm.jenis_jabatan" required class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 rounded-lg text-xs text-slate-300 outline-none">
                <option value="Fungsional">Fungsional</option>
                <option value="Pelaksana">Pelaksana</option>
                <option value="Struktural">Struktural</option>
              </select>
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-350 mb-1">Jabatan*</label>
              <input v-model="editForm.jabatan" type="text" placeholder="Contoh: Guru Madya, Staf IT" required class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 rounded-lg text-xs text-slate-100 outline-none" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-350 mb-1">Tempat Tugas*</label>
              <input v-model="editForm.tempat_tugas" type="text" placeholder="Unit Kerja / Sekolah" required class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 rounded-lg text-xs text-slate-100 outline-none" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-350 mb-1">Jenis Unit Kerja*</label>
              <select v-model="editForm.jenis_tempat" required class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 rounded-lg text-xs text-slate-300 outline-none">
                <option value="Dinas">Dinas</option>
                <option value="Sekolah">Sekolah</option>
              </select>
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-350 mb-1">Tempat Lahir</label>
              <input v-model="editForm.tempat_lahir" type="text" placeholder="Kota Lahir" class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 rounded-lg text-xs text-slate-100 outline-none" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-350 mb-1">Tanggal Lahir</label>
              <DatePicker
                v-model="editForm.tanggal_lahir"
                inputClass="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 rounded-lg text-xs text-slate-100 outline-none"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-350 mb-1">Pengangkatan</label>
              <input v-model="editForm.pengangkatan" type="text" placeholder="Contoh: CPNS 2019" class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 rounded-lg text-xs text-slate-100 outline-none" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-350 mb-1">Jenis Kepegawaian</label>
              <select v-model="editForm.jenis_pengangkatan" class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 rounded-lg text-xs text-slate-300 outline-none">
                <option value="PNS">PNS</option>
                <option value="PPPK">PPPK</option>
              </select>
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-350 mb-1">Tanggal KGB Terakhir</label>
              <DatePicker
                v-model="editForm.tanggal_kgb_terakhir"
                inputClass="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 rounded-lg text-xs text-slate-100 outline-none"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-350 mb-1">Tanggal Pangkat Terakhir</label>
              <DatePicker
                v-model="editForm.tanggal_kenaikan_pangkat_terakhir"
                inputClass="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 rounded-lg text-xs text-slate-100 outline-none"
              />
            </div>
          </div>

          <!-- Document Overwrite Uploads -->
          <div class="border-t border-slate-800 pt-4 space-y-3">
            <h3 class="text-xs font-bold text-white uppercase tracking-wider">Unggah File Lampiran SK (PDF/JPG/PNG Max 5MB)</h3>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-[10px] text-slate-400 mb-1">Foto Profil</label>
                <input type="file" @change="onEditFileChange($event, 'foto_profil')" accept="image/*" class="text-xs text-slate-400" />
              </div>
              <div>
                <label class="block text-[10px] text-slate-400 mb-1">SK CPNS / PPPK</label>
                <input type="file" @change="onEditFileChange($event, 'sk_cpns_pppk_file')" accept=".pdf,image/*" class="text-xs text-slate-400" />
              </div>
              <div>
                <label class="block text-[10px] text-slate-400 mb-1">SK PNS (Jika Ada)</label>
                <input type="file" @change="onEditFileChange($event, 'sk_pns_file')" accept=".pdf,image/*" class="text-xs text-slate-400" />
              </div>
              <div>
                <label class="block text-[10px] text-slate-400 mb-1">SK KGB Terbaru</label>
                <input type="file" @change="onEditFileChange($event, 'sk_kgb_file')" accept=".pdf,image/*" class="text-xs text-slate-400" />
              </div>
              <div>
                <label class="block text-[10px] text-slate-400 mb-1">SK Pangkat Terbaru</label>
                <input type="file" @change="onEditFileChange($event, 'sk_pangkat_file')" accept=".pdf,image/*" class="text-xs text-slate-400" />
              </div>
              <div>
                <label class="block text-[10px] text-slate-400 mb-1">SK Pensiun</label>
                <input type="file" @change="onEditFileChange($event, 'sk_pensiun_file')" accept=".pdf,image/*" class="text-xs text-slate-400" />
              </div>
              <div class="sm:col-span-2">
                <label class="block text-[10px] text-slate-400 mb-1">Dokumen Pemberhentian Pembayaran (Opsi B)</label>
                <input type="file" @change="onEditFileChange($event, 'dokumen_pemberhentian_pembayaran')" accept=".pdf,image/*" class="w-full text-xs text-slate-400" />
              </div>
            </div>
          </div>

          <div class="flex justify-end gap-2 pt-2 border-t border-slate-800">
            <button type="button" @click="showEditModal = false" class="px-3 py-1.5 bg-slate-800 hover:bg-slate-750 text-slate-300 rounded-lg text-xs font-semibold">Batal</button>
            <button type="submit" :disabled="modalSubmitting" class="px-3 py-1.5 bg-primary-600 hover:bg-primary-700 text-white rounded-lg text-xs font-semibold flex items-center gap-1">
              <Loader2 v-if="modalSubmitting" class="w-3 animate-spin" />
              <span>Simpan Perubahan</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import api from '@/services/api'
import { X, Loader2, Eye, Download } from 'lucide-vue-next'
import DatePicker from '@/components/DatePicker.vue'
import Swal from 'sweetalert2'

const route = useRoute()
const authStore = useAuthStore()

const nip = route.params.nip
const employee = ref({})
const ageStr = ref('')
const employeeUserPhone = ref('')
const kgbHistory = ref([])
const pangkatHistory = ref([])

const currentTab = ref('info')
const tabs = [
  { id: 'info', name: 'Detail Profil' },
  { id: 'docs', name: 'Berkas SK' },
  { id: 'history', name: 'Riwayat KGB & Pangkat' }
]

// Modal triggers
const showKgbModal = ref(false)
const showPangkatModal = ref(false)
const showEditModal = ref(false)
const modalSubmitting = ref(false)

const kgbForm = ref({ tanggal_kgb: '' })
const pangkatForm = ref({ tanggal_kenaikan_pangkat: '' })

// Preview state
const previewUrl = ref(null)
const previewLabel = ref('')
const previewIsPdf = computed(() => previewUrl.value?.toLowerCase().includes('.pdf'))

const openPreview = (url, label) => {
  previewUrl.value = url
  previewLabel.value = label
}
const closePreview = () => { previewUrl.value = null }

// SK Documents config for rendering tab berkas
const skDocuments = computed(() => [
  { key: 'sk_cpns_pppk_file', label: 'SK CPNS / PPPK Pertama', desc: 'Berkas awal pengangkatan pegawai' },
  { key: 'sk_pns_file', label: 'SK PNS (Pegawai Negeri Sipil)', desc: 'Dokumen sumpah PNS 100%' },
  { key: 'sk_kgb_file', label: 'SK KGB Terbaru', desc: 'Dokumen Kenaikan Gaji Berkala aktif' },
  { key: 'sk_pangkat_file', label: 'SK Kenaikan Pangkat Terbaru', desc: 'SK golongan pangkat/ruang terakhir' },
  { key: 'sk_pensiun_file', label: 'SK Pensiun', desc: 'SK pemberhentian hormat pensiun' },
  { key: 'dokumen_pemberhentian_pembayaran', label: 'Dokumen Pemberhentian Pembayaran', desc: 'Berkas syarat final status pensiun (Opsi B)' },
])

// Helper: check if date is valid (not zero/null)
// Go's zero time serializes as "0001-01-01T00:00:00Z" or "0001-01-01"
const hasDate = (dateStr) => {
  if (!dateStr) return false
  if (dateStr === '' || dateStr === null || dateStr === undefined) return false
  if (dateStr.startsWith('0001-01-01') || dateStr.startsWith('0000')) return false
  // Ensure year is reasonable
  const year = parseInt(dateStr.substring(0, 4))
  if (isNaN(year) || year <= 1) return false
  return true
}
const editForm = ref({
  nama: '',
  no_hp: '',
  jenis_jabatan: 'Fungsional',
  jabatan: '',
  tempat_tugas: '',
  jenis_tempat: 'Dinas',
  tempat_lahir: '',
  tanggal_lahir: '',
  pengangkatan: '',
  jenis_pengangkatan: 'PNS',
  tanggal_kgb_terakhir: '',
  tanggal_kenaikan_pangkat_terakhir: ''
})

const kgbFile = ref(null)
const pangkatFile = ref(null)
const editFiles = ref({
  foto_profil: null,
  sk_cpns_pppk_file: null,
  sk_pns_file: null,
  sk_kgb_file: null,
  sk_pangkat_file: null,
  sk_pensiun_file: null,
  dokumen_pemberhentian_pembayaran: null
})

const fetchData = async () => {
  try {
    const response = await api.get(`/employees/${nip}`)
    employee.value = response.data.employee
    ageStr.value = response.data.umur
    
    // Populate edit form values
    const safeDate = (d) => {
      if (!d || d.startsWith('0001-01-01') || d.startsWith('0000')) return ''
      const year = parseInt(d.substring(0, 4))
      if (isNaN(year) || year <= 1) return ''
      return d.substring(0, 10)
    }

    editForm.value = {
      nama: employee.value.nama,
      jenis_jabatan: employee.value.jenis_jabatan,
      jabatan: employee.value.jabatan,
      tempat_tugas: employee.value.tempat_tugas,
      jenis_tempat: employee.value.jenis_tempat,
      tempat_lahir: employee.value.tempat_lahir,
      tanggal_lahir: safeDate(employee.value.tanggal_lahir),
      pengangkatan: employee.value.pengangkatan || '',
      jenis_pengangkatan: employee.value.jenis_pengangkatan || 'PNS',
      tanggal_kgb_terakhir: safeDate(employee.value.tanggal_kgb_terakhir),
      tanggal_kenaikan_pangkat_terakhir: safeDate(employee.value.tanggal_kenaikan_pangkat_terakhir)
    }

    // Load Histories
    const kgbRes = await api.get(`/employees/${nip}/kgb-history`)
    kgbHistory.value = kgbRes.data.data || []

    const pangkatRes = await api.get(`/employees/${nip}/pangkat-history`)
    pangkatHistory.value = pangkatRes.data.data || []

    // Fetch user phone from credential user table
    const allEmpResponse = await api.get('/employees')
    const match = allEmpResponse.data.data?.find(u => u.nip === nip)
    if (match) {
      employeeUserPhone.value = match.no_hp
      editForm.value.no_hp = match.no_hp
    } else if (authStore.nip === nip) {
      employeeUserPhone.value = authStore.employee?.no_hp || authStore.no_hp
      editForm.value.no_hp = employeeUserPhone.value
    }
  } catch (error) {
    console.error('Gagal mengambil detail pegawai:', error)
  }
}

const openKgbModal = () => {
  kgbForm.value.tanggal_kgb = ''
  kgbFile.value = null
  showKgbModal.value = true
}

const openPangkatModal = () => {
  pangkatForm.value.tanggal_kenaikan_pangkat = ''
  pangkatFile.value = null
  showPangkatModal.value = true
}

const onHistoryFileChange = (e, field) => {
  const file = e.target.files[0]
  if (file) {
    if (field === 'file_sk_kgb') kgbFile.value = file
    if (field === 'file_sk_pangkat') pangkatFile.value = file
  }
}

const onEditFileChange = (e, field) => {
  const file = e.target.files[0]
  if (file) {
    editFiles.value[field] = file
  }
}

const submitKgb = async () => {
  modalSubmitting.value = true
  const formData = new FormData()
  formData.append('tanggal_kgb', kgbForm.value.tanggal_kgb)
  formData.append('file_sk_kgb', kgbFile.value)

  try {
    await api.post(`/employees/${nip}/kgb-history`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    showKgbModal.value = false
    await fetchData()
    Swal.fire({
      icon: 'success',
      title: 'Berhasil',
      text: 'Data KGB berhasil disimpan.',
      timer: 1500,
      showConfirmButton: false
    })
  } catch (error) {
    Swal.fire({
      icon: 'error',
      title: 'Gagal',
      text: error.response?.data?.error || 'Gagal menyimpan KGB baru',
      confirmButtonColor: '#1A365D'
    })
  } finally {
    modalSubmitting.value = false
  }
}

const submitPangkat = async () => {
  modalSubmitting.value = true
  const formData = new FormData()
  formData.append('tanggal_kenaikan_pangkat', pangkatForm.value.tanggal_kenaikan_pangkat)
  formData.append('file_sk_pangkat', pangkatFile.value)

  try {
    await api.post(`/employees/${nip}/pangkat-history`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    showPangkatModal.value = false
    await fetchData()
    Swal.fire({
      icon: 'success',
      title: 'Berhasil',
      text: 'Data Pangkat berhasil disimpan.',
      timer: 1500,
      showConfirmButton: false
    })
  } catch (error) {
    Swal.fire({
      icon: 'error',
      title: 'Gagal',
      text: error.response?.data?.error || 'Gagal menyimpan Pangkat baru',
      confirmButtonColor: '#1A365D'
    })
  } finally {
    modalSubmitting.value = false
  }
}

const submitEditProfile = async () => {
  modalSubmitting.value = true
  const formData = new FormData()
  for (const [key, value] of Object.entries(editForm.value)) {
    // Jangan kirim undefined/null sebagai string — kirim string kosong
    formData.append(key, value ?? '')
  }
  for (const [key, file] of Object.entries(editFiles.value)) {
    if (file) {
      formData.append(key, file)
    }
  }

  try {
    await api.put(`/employees/${nip}`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    showEditModal.value = false
    await fetchData()
    Swal.fire({
      icon: 'success',
      title: 'Berhasil',
      text: 'Perubahan profil berhasil disimpan.',
      timer: 1500,
      showConfirmButton: false
    })
  } catch (error) {
    Swal.fire({
      icon: 'error',
      title: 'Gagal',
      text: error.response?.data?.error || 'Gagal menyimpan perubahan profil',
      confirmButtonColor: '#1A365D'
    })
  } finally {
    modalSubmitting.value = false
  }
}

const formatDate = (dateStr) => {
  if (!dateStr || dateStr.startsWith('0001-01-01')) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' })
}

const formatDateTime = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return `${d.toLocaleDateString('id-ID')} ${d.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })}`
}

onMounted(() => {
  fetchData()
})
</script>
