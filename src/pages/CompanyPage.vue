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


    <v-card class="ma-2" elevation="0" v-else-if="companyData != undefined">
     <v-list>
       <v-list-item v-if="companyData.svg != undefined"
           :prepend-avatar="companyData.svg">
          <v-list-item-title>{{companyData.name}}</v-list-item-title>
          <v-list-item-subtitle>by {{companyData.ceo}}</v-list-item-subtitle>
       </v-list-item>
       <v-list-item v-else>
         <v-list-item-title>{{companyData.name}}</v-list-item-title>
         <v-list-item-subtitle>by {{companyData.ceo}}</v-list-item-subtitle>
       </v-list-item>

       <v-spacer></v-spacer>
       <v-list-item v-if="companyData.description">
         {{companyData.description}}
       </v-list-item>

       <v-list-item v-if="companyData.year || companyData.staffSize || companyData.products.length > 0">
         <v-chip-group column>
           <v-chip variant="outlined" v-if="companyData.year"><v-icon class="mr-2">mdi-calendar-range</v-icon>Created at {{companyData.year}}</v-chip>
           <v-chip variant="outlined" v-if="companyData.staffSize"><v-icon class="mr-2">mdi-account-group</v-icon>Staff size: {{companyData.staffSize}}</v-chip>
           <v-chip variant="outlined" v-if="companyData.products.length > 0"><v-icon class="mr-2">mdi-checkbox-multiple-blank-circle</v-icon>{{companyData.products.length}} products</v-chip>
         </v-chip-group>
       </v-list-item>

       <v-list-item v-if="companyData.departments.length != 0">
         <v-chip-group column>
           <v-chip v-for="department in companyData.departments">{{department}}</v-chip>
         </v-chip-group>
       </v-list-item>

       <v-card class="ma-3">
         <v-card-title>Продукты</v-card-title>
         <v-list v-if="companyData.products.length == 0">
           <v-list-item>
           <v-list-item-title>Продукты не найдены...</v-list-item-title>
           </v-list-item>
         </v-list>
         <v-list v-else>
           <v-list-item v-for="product in companyData.products">
             <v-list-item-title>
               <span v-if="product.isVerified" class="text-green">{{product.name}}</span>
               <span v-else>{{product.name}}</span>
             </v-list-item-title>
             <v-list-item-subtitle>since {{product.year}}</v-list-item-subtitle>
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
  name: "CompanyPage",
  data() {
    return {
      currentId: undefined,
      isLoading: false,
      isLoadingFailure: false,
      companyData: undefined,
    }

  },
  methods: {
    fetchCompany: function(id) {
      this.currentId = id
      this.isLoading = true
      this.isLoadingFailure = false
      axios.get("http://141.95.127.215:7328/company?id="+id).then(response => {
        this.companyData = response.data
        this.isLoading = false
        this.emitter.emit("select-data", {nodeType: "node", id: this.companyData.id})
      }).catch(e => {
        console.log("error")
        this.isLoading = false
        this.isLoadingFailure = true
      })
    }
  },
  beforeRouteUpdate(to, from) {
    this.fetchCompany(to.params.id)
  },
  mounted() {
    this.fetchCompany(this.$route.params.id)
  }
}

</script>

<style scoped>
  .title {
    word-break: break-all;
  }
</style>