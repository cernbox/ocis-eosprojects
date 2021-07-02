import 'regenerator-runtime/runtime'
import App from './components/App.vue'
import store from './store'

const appInfo = {
  name: 'EosProjects',
  id: 'eosprojects',
  icon: 'info',
  isFileEditor: false,
  extensions: []
}

const routes = [
  {
    name: 'projects',
    path: '/',
    components: {
      app: App
    }
  }
]

const navItems = [
  {
    name: 'Projects',
    iconMaterial: appInfo.icon,
    route: {
      name: 'projects',
      path: `/${appInfo.id}/`
    }
  }
]

export default {
  appInfo,
  store,
  routes,
  navItems
}
