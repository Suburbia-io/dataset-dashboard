<template>
  <table class="table is-striped is-hoverable is-fullwidth is-narrow">
    <thead>
    <tr>
      <th class="action-buttons"></th>
      <th class="google-search"></th>
      <th class="location-string">Location String</th>
      <th class="parsed-country">Parsed Country</th>
      <th class="parsed-postal-code">Parsed Postal Code</th>
      <th class="postal-codes">GoeName Postal Codes</th>
      <th class="geoname">GoeName</th>
      <th class="geoname-id">GeoName ID</th>
      <th class="parent-geoname-id">Parent GeoName ID</th>
      <th class="approved">Approved</th>
    </tr>
    </thead>
    <tbody>
    <locations-table-row
      v-for="l in this.locations"
      :key="l.locationHash"
      :location="l">
    </locations-table-row>
    <tr v-if="this.locations.length === 0">
      <td colspan="10" class="has-text-centered">No results</td>
    </tr>
    </tbody>
    <tfoot></tfoot>
  </table>
</template>

<script>
  import LocationsTableRow from "./LocationsTableRow";
  export default {
    name: "LocationsTable",
    components: {LocationsTableRow},
    beforeCreate() {
      this.locationsJSON = JSON.parse(document.getElementById('locations').getAttribute('data-locations'));
    },
    data() {
      return {
        locations: this.locationsJSON !== null ? this.locationsJSON : [],
      }
    }
  }
</script>
