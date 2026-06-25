import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './assets/index.css'

// Import DataTables CSS
import 'datatables.net-dt/css/dataTables.dataTables.css'
// Import SweetAlert2 CSS
import 'sweetalert2/dist/sweetalert2.min.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')
