<template>
  <div class="interaction d-flex">
    <div>
    <div class="m-3">
      <div class="d-flex flex-column">
        <v-btn type="button" class="btn btn-secondary ma-1" @click="zooming(1)">+</v-btn>
        <v-btn type="button" class="btn btn-secondary ma-1" @click="zooming(-1)">-</v-btn>
      </div>
    </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "MapInteractionTools",
  props: {
    floors: {
      type: Array
    },
    currentFloor: { 
      default: 1 
    },
    isDataLoading: {
      default: false
    }
  },
  methods: {
    zooming(value) {
      this.$emit("zoomingMap", value)
    },
  }
}
</script>

<style scoped>

</style>