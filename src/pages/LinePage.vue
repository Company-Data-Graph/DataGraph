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
    <v-card class="ma-2" elevation="0" v-else-if="linkData != undefined">
      <v-list>
        <v-list-item>
          <v-list-item-title>Time section</v-list-item-title>
          <v-list-item-subtitle>from {{linkData.source.year}} to {{linkData.target.year}}</v-list-item-subtitle>
        </v-list-item>

        <v-card class="ma-3">
          <v-card-title>Linked elements</v-card-title>
          <v-list>
            <v-list-item>
              <v-list-item-title>{{this.linkData.source.name}}</v-list-item-title>
              <v-list-item-subtitle>since {{this.linkData.source.year}}</v-list-item-subtitle>
            </v-list-item>
            <v-list-item>
              <v-list-item-title>{{this.linkData.target.name}}</v-list-item-title>
              <v-list-item-subtitle>since {{this.linkData.source.year}}</v-list-item-subtitle>
            </v-list-item>
          </v-list>
        </v-card>
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
      linkData: undefined,
    }

  },
  methods: {
    isExist: function(field) {
      return field != "" && field != undefined && field != null
    },
    fetchLine: function(source, target, mode) {
      console.log(mode)
      this.isLoading = true
      axios.get("http://141.95.127.215:7328/link/" + (mode=="company" ? "company" : "products") + "?source=" + source + "&target=" + target).then(response => {
        this.linkData = response.data
        this.isLoading = false
        this.emitter.emit("select-data", {nodeType: "line", source: this.linkData.source.id, target: this.linkData.target.id})
      })
    }
  },
  beforeRouteUpdate(to, from) {
    console.log(to)
    this.fetchLine(to.query.source, to.query.target, to.query.mode)
  },
  mounted() {
    this.fetchLine(this.$route.query.source, this.$route.query.target, this.$route.query.mode)
  }
}

</script>

<style scoped>
.clearA {
  text-decoration: none;
}

</style>