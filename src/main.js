import { createApp } from 'vue'
import App from './App.vue'
import router from "@/router/panel-router";
import component from "@/components/UIElements"
import mitt from 'mitt';

import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'


const vuetify = createVuetify({
    components,
    directives,
})


const emitter = mitt()
const app = createApp(App)

component.forEach(component => {
    app.component(component.name, component)
})

app.config.globalProperties.emitter = emitter

app
    .use(router)
    .use(vuetify)
    .mount('#app')
