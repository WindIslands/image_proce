import {createApp} from 'vue'
import naive from 'naive-ui'
import {router} from './routers'

import App from './App.vue'
import './style.css';

// 定义特性标志
window.__VUE_PROD_DEVTOOLS__ = false;
window.__VUE_PROD_HYDRATION_MISMATCH_DETAILS__ = false;


const app = createApp(App)
app.use(naive)
app.use(router)
app.mount('#app')
