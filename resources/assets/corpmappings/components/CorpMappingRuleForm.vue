<template>
  <tr>
    <th>
      <input class="input  is-small" type="date" min="1990-01-01" v-model="startDate">
    </th>
    <th>
      <v-select :options="countries" v-model="country" placeholder="Empty means *">
      </v-select>
    </th>
    <th>
      <v-select :options="corps" v-model="corporation" label="name">
      </v-select>
    </th>
    <th>
      <button class="button is-success is-small" @click="addCorpMappingRule"><i class="fas fa-plus"></i></button>
    </th>
  </tr>
</template>

<script>
  import isoCountries from "../../shared/js/isoCountries";

  const defaultDate = '1990-01-01';

  export default {
    props: ['corps', 'mappingId', 'callback'],
    data: () => ({
      startDate: defaultDate,
      country: '',
      corporation: null,
      countries: []
    }),
    mounted() {
      this.countries = isoCountries;
    },
    methods: {
      addCorpMappingRule() {
        this.callback(this.mappingId, this.startDate, this.country, this.corporation.corporationID);
        this.startDate = defaultDate;
        this.country = '';
        this.corporationID = '';
      }
    }
  }
</script>

<style scoped>

</style>
