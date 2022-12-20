<template>
  <v-list-item>
    <v-form v-if="companyFiltersPresets == undefined">
      Configuration...
    </v-form>
    <v-form v-else class="mt-2">
      <v-card elevation="0">
        <v-card-text>
      <v-text-field label="Company Name" variant="outlined" v-model="companyName" />
      <v-text-field label="Company SEO" variant="outlined" v-model="companySeo" />
        </v-card-text>
      </v-card>

      <v-card elevation="0">
        <v-card-subtitle>
          Company years
        </v-card-subtitle>
        <v-card-text>
          <v-range-slider
              v-model="yearRange"
              :min="companyFiltersPresets.minDate.substring(0, 4)"
              :max="companyFiltersPresets.maxDate.substring(0, 4)"
              :step="1" hide-details color="blue"
              class="align-center"
              thumb-label="true"
          >
          </v-range-slider>
        </v-card-text>
      </v-card>

      <v-card elevation="0">
        <v-card-subtitle>
          Company staff count
        </v-card-subtitle>
        <v-card-text>
          <v-range-slider
              v-model="staffRange"
              :min="companyFiltersPresets.minStaffSize"
              :max="companyFiltersPresets.maxStaffSize"
              :step="1" hide-details color="blue"
              class="align-center"
              thumb-label="true"
          >
          </v-range-slider>
        </v-card-text>
      </v-card>

      <v-card elevation="0">
        <v-card-text>
          <v-combobox
              v-model="departments"
              :items="companyFiltersPresets.departments.map(el => el.name)"
              label="Destinations"
              multiple variant="outlined"
          ></v-combobox>
        </v-card-text>
      </v-card>


      <v-btn color="blue" class="pl-3 ma-2 mr-3" @click.prevent="submit()">Submit</v-btn>
      <v-btn color="blue" class="ma-2" variant="text" @click.prevent="drop()">Drop</v-btn>
    </v-form>
  </v-list-item>
</template>

<script>
export default {
  name: "CompanyFilters",
  data() {
    return {
      departments: [],
      companyName: "",
      companySeo: "",
      yearRange: [],
      staffRange: [],

    }
  },
  watch: {
    companyFiltersPresets: {
      handler(newValue) {
        this.yearRange = [newValue.minDate.substring(0, 4), newValue.maxDate.substring(0, 4)]
        this.staffRange[0] = newValue.minStaffSize
        this.staffRange[1] = newValue.maxStaffSize
      },
      deep: true
    },
  },
  props: {
    companyFiltersPresets: {
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
        companyName: this.companyName,
        ceo: this.companySeo,
        departments: this.companyFiltersPresets.departments.filter(el => this.departments.includes(el.name)).map(el => el.id),
        minDate: JSON.stringify(this.yearRange[0]) + "-01-01T00:00:00Z",
        maxDate: JSON.stringify(this.yearRange[1]) + "-01-01T00:00:00Z",
        startStaffSize: this.staffRange[0],
        endStaffSize: this.staffRange[1]
      })
    },
    drop: function() {
      this.companyName = ""
      this.companySeo = ""
      this.departments = []
      this.yearRange = [0, 0]
      this.staffRange = [0, 0]
      this.dropHandler()
    }
  }
}
</script>

<style scoped>

</style>