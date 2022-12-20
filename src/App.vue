<template>
  <div class="app">
    <v-app>
      <v-layout style="z-index: 0">
        <v-app-bar color="blue">
          <v-app-bar-title>Data Graph</v-app-bar-title>
          <v-spacer />
          <div>
            <v-btn @click.prevent="filtersShow = !filtersShow"><v-icon>mdi-filter</v-icon></v-btn>
          </div>
        </v-app-bar>
        <v-navigation-drawer v-model="filtersShow" :width="panelSize" permanent location="right">
          <filters  />
        </v-navigation-drawer>
        <v-navigation-drawer v-model="panelShow" :width="panelSize" permanent location="left">
          <panel :set-panel-state="setPanelState" />
        </v-navigation-drawer>

        <v-main>
          <v-container fluid>
            <interactive-map
              :graph-interactions="graphInteractions"
              :set-panel-state="setPanelState"
              :selected-node="selectedNode"
              :current-filters-company="currentFiltersCompany"
              :current-filters-product="currentFiltersProduct"
              @setPanelState="setPanelState"
            />
            <map-interaction-tools
                class="interaction-tools"
                @zoomingMap="zoomingMap"
            />
          </v-container>
        </v-main>
      </v-layout>
    </v-app>

  </div>
</template>

<script>
import InteractiveMap from "@/components/Map";
import Panel from "@/components/Panel";
import MapInteractionTools from "@/components/MapInteractionTools";
import Filters from "@/components/Filters.vue";


export default {
  components: {Filters, InteractiveMap, Panel, MapInteractionTools },
  data() {
    return {
      graphInteractions: {
        'newScaleFactor': 0,
        'oldScaleFactor': 0,
      },
      currentFiltersCompany: undefined,
      currentFiltersProduct: undefined,
      panelShow: null,
      filtersShow: null,
      minPanelSize: 400,
      panelSize: 200,
      selectedNode: undefined
    }
  },
  methods: {
    setPanelState: function(value) {
      this.panelShow = value
    },
    zoomingMap: function(value) {
      this.graphInteractions.oldScaleFactor = this.graphInteractions.newScaleFactor
      this.graphInteractions.newScaleFactor = this.graphInteractions.oldScaleFactor + value
    },
    getPanelSize: function() {
      let width = window.innerWidth * 0.4
      width = (width < this.minPanelSize) ? this.minPanelSize : width
      return width
    },

  },

  mounted() {
    this.panelSize = this.getPanelSize()
    window.addEventListener('resize', () => {
      this.panelSize = this.getPanelSize()
      console.log(this.panelSize)
    })
  },

  created (){
    this.emitter.on('filter-graph-company', (evt) => {
      this.currentFiltersCompany = evt;
    })

    this.emitter.on('clear-filter-graph-company', () => {
      this.currentFiltersCompany = undefined;
    })
    this.emitter.on('filter-graph-product', (evt) => {
      this.currentFiltersProduct = evt;
    })

    this.emitter.on('clear-filter-graph-product', () => {
      this.currentFiltersProduct = undefined;
    })



    this.emitter.on('select-data', (evt) => {
      console.log("new select")
      this.selectedNode = evt;
    })
    },
}
</script>

<style scoped>
.interaction-tools {
  position: absolute;
  z-index: 50;
  right: 0em;
  bottom: 0em;
}
</style>