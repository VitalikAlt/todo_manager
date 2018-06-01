import Vue from "vue";
import ElementUI from "element-ui";
import "element-ui/lib/theme-chalk/index.css";
import App from "./App.vue";
import VueNativeSock from "vue-native-websocket";

Vue.use(ElementUI);
Vue.use(VueNativeSock, `ws://${window.location.hostname}:8081`, {
  format: "json"
});

new Vue({
  el: "#app",
  render: h => h(App)
});
