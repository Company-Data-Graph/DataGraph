<template>
<v-card elevation="0">
    <v-card elevation="0" class="d-flex pa-2 justify-end">
    <v-btn variant="plain" @click.prevent="setPanelState(false)"><v-icon>mdi-close</v-icon></v-btn>
    </v-card>
  <router-view></router-view>
</v-card>
</template>

<script>
export default {
  data() {
    return {
      currentParams: '',
      isFullScreen: true,
    }
  },
  props: {
    setPanelState: {
      required: true,
      type: Function
    }
  },
  methods: {
    changePanelSize: function(newValue) {
      this.isFullScreen = newValue
      this.emitter.emit('change-panel-size', this.isFullScreen)
    }
  },
  created() {
    this.emitter.on('change-panel-size', (evt) => {
      this.isFullScreen = evt
    })
  },
  mounted() {
    this.changePanelSize(false)
  },
  name: "Panel",
}
</script>

<style scoped>
</style>