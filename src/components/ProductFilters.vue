<template>
  <v-list-item>
    <v-form v-if="productFiltersPresets == undefined">
      Configuration...
    </v-form>
    <v-form v-else class="mt-2">
      <v-card elevation="0">
        <v-card-text>
      <v-text-field label="Product Name" variant="outlined" v-model="productName" />
        </v-card-text>
      </v-card>

      <v-card elevation="0">
        <v-card-subtitle>
          Product years
        </v-card-subtitle>
        <v-card-text>
          <v-range-slider
              v-model="yearRange"
              :min="productFiltersPresets.minDate.substring(0, 4)"
              :max="productFiltersPresets.maxDate.substring(0, 4)"
              :step="1" hide-details color="blue"
              class="align-center"
              thumb-label="true"
          >
          </v-range-slider>
        </v-card-text>
      </v-card>

      <v-card elevation="0">
        <v-card-text>
          <v-checkbox
              v-model="isVerified"
              label="Is Verified release"
              color="blue"
              hide-details
          ></v-checkbox>
        </v-card-text>
      </v-card>

      <v-btn color="blue" class="pl-3 ma-2 mr-3" @click.prevent="submit()">Submit</v-btn>
      <v-btn color="blue" class="ma-2" variant="text" @click.prevent="drop()">Drop</v-btn>
    </v-form>
  </v-list-item>
</template>

<script>
export default {
  name: "ProductFilters",
  data() {
    return {
      productName: "",
      yearRange: [],
      isVerified: true
    }
  },
  watch: {
    productFiltersPresets: {
      handler(newValue) {
        this.yearRange = [newValue.minDate.substring(0, 4), newValue.maxDate.substring(0, 4)]
      },
      deep: true
    },
  },
  props: {
    productFiltersPresets: {
      required: true,
      default: undefined,
    },
    submitHandler: {
      type: Function,
      required: true
    },
    dropHandler: {
      type: Function,
      required: true
    }
  },
  methods: {
    submit: function () {
      this.submitHandler({
        productName: this.productName,
        minDate: JSON.stringify(this.yearRange[0]) + "-01-01T00:00:00Z",
        maxDate: JSON.stringify(this.yearRange[1]) + "-01-01T00:00:00Z",
        isVerified: this.isVerified
      })
    },
    drop: function () {
      this.productName = ""
      this.yearRange = [0, 0]
      this.isVerified = false
      this.dropHandler()
    }
  }
}
</script>

<style scoped>

</style>