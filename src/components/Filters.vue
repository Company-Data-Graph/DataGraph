<template>
  <div>
    <v-card elevation="0">
      <v-card-text>
      <v-list>
        <v-list-item>
          <v-list-item-title>FILTERS</v-list-item-title>
        </v-list-item>
        <v-list-item>
          <v-tabs color="blue" v-model="tab">
            <v-tab value="company">Company</v-tab>
            <v-tab value="product">Products and releases</v-tab>
          </v-tabs>

        </v-list-item>
        <v-list-item>
          <v-card v-if="isLoadingFilterPresets" elevation="0">
            <v-list>
              <v-list-item>
                <v-progress-circular
                    :size="24"
                    color="blue"
                    indeterminate
                    class="mr-2"
                ></v-progress-circular>
                <span>Loading...</span>
              </v-list-item>
            </v-list>
          </v-card>

          <v-window v-model="tab" v-else-if="filterPresets != undefined">
            <v-window-item value="company">
              <company-filters
                               :submit-handler="invokeCompanyFilters"
                               :drop-handler="dropCompanyFilters"
                               :company-filters-presets="filterPresets.companyFilters"/>
            </v-window-item>
            <v-window-item value="product">
              <product-filters
                               :drop-handler="dropProductFilters"
                               :submit-handler="invokeProductFilters"
                               :product-filters-presets="filterPresets.productFilters" />
            </v-window-item>
          </v-window>



        </v-list-item>
      </v-list>
      </v-card-text>
    </v-card>
  </div>
</template>

<script>
import axios from "axios";
import CompanyFilters from "@/components/CompanyFilters.vue";
import ProductFilters from "@/components/ProductFilters.vue";

export default {
  name: "Filters",
  components: {ProductFilters, CompanyFilters},
  data() {
    return {
      filterPresets: undefined,
      isLoadingFilterPresets: false,
      dialog: null,
      tab: "company",
    }
  },
  methods: {
    getFilterPresets: function () {
      this.isLoadingFilterPresets = true
      axios.get("http://141.95.127.215:7328/filterPresets").then(response => {
        console.log(response.data)
        this.filterPresets = response.data
        this.isLoadingFilterPresets = false

      })
    },
    invokeCompanyFilters: function (data) {
      axios.post("http://141.95.127.215:7328/filterCompany", data).then(response => {
        console.log(response.data)
        this.dialog = false
        this.emitter.emit('filter-graph-company', response.data)
      })
    },
    invokeProductFilters: function (data) {
      axios.post("http://141.95.127.215:7328/filterProduct", data).then(response => {
        console.log(response.data)
        this.dialog = false
        this.emitter.emit('filter-graph-product', response.data)
      })
    },
    dropCompanyFilters: function() {
      this.dialog = false
      this.emitter.emit('clear-filter-graph-company')
    },

    dropProductFilters: function() {
      this.dialog = false
      this.emitter.emit('clear-filter-graph-product')
    }
  },
  mounted() {
    this.getFilterPresets()
  }
}
</script>

<style scoped>
</style>