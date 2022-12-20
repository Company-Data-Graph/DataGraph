<template>
  <div>
    <v-card v-if="isLoading" elevation="0">
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
    <v-card class="ma-2" elevation="0" v-else-if="productData != undefined">
      <v-list>
        <v-list-item v-if="isExist(productData.svg)"
                     :prepend-avatar="productData.svg"
        >
          <v-list-item-title>{{productData.name}} v{{productData.version}}</v-list-item-title>
          <v-list-item-subtitle>by {{productData.company.name}}</v-list-item-subtitle>
        </v-list-item>
        <v-list-item v-else>
          <v-list-item-title>{{productData.name}} v{{productData.version}}</v-list-item-title>
          <v-list-item-subtitle>by {{productData.company.name}}</v-list-item-subtitle>
        </v-list-item>

        <v-list-item v-if="isExist(productData.link)">
          <a class="clearA" :href="productData.link"><v-chip color="green" text-color="white">Верифицированный релиз</v-chip></a>
        </v-list-item>

        <v-list-item v-if="productData.description">
          <span class="font-weight-bold">Release notes:</span> {{productData.description}}
        </v-list-item>

        <v-list-item v-if="productData.year">
          <v-chip-group column>
            <v-chip variant="outlined" v-if="productData.year">Created at {{productData.year}}</v-chip>
          </v-chip-group>
        </v-list-item>

        <v-list-item v-if="productData.departments.length != 0">
          <v-chip-group column>
            <v-chip v-for="department in productData.departments">{{department.name}}</v-chip>
          </v-chip-group>
        </v-list-item>
      </v-list>
    </v-card>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "ProductPage",
  data() {
    return {
      isLoading: false,
      productData: undefined,
    }

  },
  methods: {
    isExist: function(field) {
      return field != "" && field != undefined && field != null
    },
    fetchProduct: function(id) {
      this.isLoading = true
      axios.get("http://141.95.127.215:7328/product?id="+id).then(response => {
        this.productData = response.data
        this.isLoading = false
        this.emitter.emit("select-data", {nodeType: "node", id: this.productData.id})
      })
    }
  },
  beforeRouteUpdate(to, from) {
    this.fetchProduct(to.params.id)
  },
  mounted() {
    this.fetchProduct(this.$route.params.id)
  }
}

</script>

<style scoped>
.clearA {
  text-decoration: none;
}

</style>