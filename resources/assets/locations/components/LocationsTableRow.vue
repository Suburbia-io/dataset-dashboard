<template>
  <tr>
    <td class="action-buttons has-text-centered">

      <div class="field is-grouped is-pulled-left" v-show="areModerateButtonsVisible">
        <p class="control">
          <button class="button is-small is-outlined" @click="showGeonameIdInput">Edit ID</button>
        </p>
        <p class="control">
          <button class="button is-small is-danger" @click="locationApprove(false)">Reject</button>
        </p>
        <p class="control" v-show="geonameID!==''">
          <button class="button is-small is-success" @click="locationApprove(true)">Approve</button>
        </p>
      </div>

      <div class="field is-pulled-left" v-show="!areModerateButtonsVisible && !isGeonameIDInputVisible">
        <p class="control">
          <button class="button is-small is-warning" @click="locationApprove(null)">Reset</button>
        </p>
      </div>
    </td>
    <td class="google-search">
      <a :href="googleSearchUrl" rel="noreferrer nofollow noopener" target="_blank">
        <span class="icon has-text-info">
          <i class="fab fa-google"></i>
        </span>
      </a>
    </td>
    <td class="location-string">{{locationString}}</td>
    <td class="parsed-country">
      <template v-if="parsedCountryCode!==''">{{parsedCountryCode}}</template>
      <template v-else>-</template>
    </td>
    <td class="parsed-postal-code">
      <template v-if="parsedPostalCode!==''">{{parsedPostalCode}}</template>
      <template v-else>-</template>
    </td>
    <td class="postal-codes">
      <div class="postal-code" v-for="pc in geonamesPostalCodes.postalCodes"
                               v-if="geonamesPostalCodes.postalCodes">
        {{pc.placeName}},
        <template v-if="'adminName4' in pc">{{pc.adminName4}}, {{pc.adminName3}},</template>
        <template v-else-if="'adminName3' in pc">{{pc.adminName3}}, {{pc.adminName2}},</template>
        <template v-else-if="'adminName2' in pc">{{pc.adminName2}}, {{pc.adminName1}},</template>
        <template v-else-if="'adminName1' in pc">{{pc.adminName1}},</template>
        {{pc.countryCode}}
      </div>
      <template v-if="!geonamesPostalCodes.postalCodes || geonamesPostalCodes.postalCodes.length === 0">-</template>
    </td>
    <td class="geoname" v-show="!isGeonameIDInputVisible">
      <template v-if="name!==''">{{geoNameDisplay}}</template>
      <template v-else>-</template>
    </td>
    <td class="geoname-id" v-show="!isGeonameIDInputVisible">
      <a v-if="geonameID!==null" :href="geoNameUrl" rel="noreferrer nofollow noopener" target="_blank">{{geonameID}}</a>
      <template v-else>-</template>
    </td>
    <td class="parent-geoname-id" v-show="!isGeonameIDInputVisible">
      <template v-if="this.parentGeoName!==null">{{parentGeoName.geonameId}}</template>
      <template v-else>-</template>
    </td>
    <td class="approved has-text-centered" v-show="!isGeonameIDInputVisible">
      <template v-if="approved===null">-</template>
      <template v-else-if="approved">Yes</template>
      <template v-else>No</template>
    </td>

    <td colspan="7" v-show="isGeonameIDInputVisible">
      <div class="field is-grouped">
        <p class="control">
          <input class="input" :class="{'is-danger': geonameIDError.length > 0}" type="text" pattern="[0-9]+"
                 placeholder="Geoname ID" v-model="geonameIDInputValue">
        </p>
        <p class="control">
          <button class="button is-small is-info" @click="setGeonameID">Save</button>
        </p>
        <p class="control">
          <button class="button is-small is-outlined" @click="hideGeonameIdInput">Cancel</button>
        </p>
      </div>
      <p class="help is-danger">{{geonameIDError}}</p>
    </td>
  </tr>
</template>

<script>
  import api from "../../locations/js/api";

  const geonameIDRegex = /^[0-9]+$/;

  export default {
    name: "LocationsTableRow",

    props: {
      location: Object,
    },

    data() {
      return {
        isGeonameIDInputVisible: false,
        geonameIDError: '',
        geonameIDInputValue: '',

        // location data:
        datasetID: '',
        locationHash: '',
        locationString: '',
        parsedCountryCode: '',
        parsedPostalCode: '',
        geonamesPostalCodes: null,
        geonamesHierarchy: null,
        name: '',
        population: '',
        countryCode: '',
        parentName: '',
        parentPopulation: '',
        geonameID: null,
        approved: null,
      }
    },

    created() {
      this.updateData(this.location);
    },

    computed: {
      googleSearchUrl() {
        let search = this.locationString.replace(/,/g, ' ');
        return `https://www.google.com/search?q=${search}`
      },

      areModerateButtonsVisible() {
        return !this.isGeonameIDInputVisible && this.approved === null;
      },

      geoName() {
        return this.getGeoNameFromHierarchy(1);
      },

      geoNameUrl() {
        return `https://www.geonames.org/${this.geonameID}/`;
      },

      geoNameDisplay() {
        if (parseInt(this.geoName.geonameId, 10) === parseInt(this.geoName.countryId, 10)) {
          return `${this.name} (${this.countryCode}), ${this.parentName}`
        } else {
          return `${this.name}, ${this.parentName}, ${this.grandParentGeoName.name}, ${this.countryCode}`;
        }
      },

      parentGeoName() {
        return this.getGeoNameFromHierarchy(2);
      },

      grandParentGeoName() {
        return this.getGeoNameFromHierarchy(3);
      },
    },

    methods: {
      updateData(location) {
        this.datasetID = location.datasetID;
        this.locationHash = location.locationHash;
        this.locationString = location.locationString;
        this.parsedCountryCode = location.parsedCountryCode;
        this.parsedPostalCode = location.parsedPostalCode;
        this.geonamesPostalCodes = location.geonamesPostalCodes;
        this.geonamesHierarchy = location.geonamesHierarchy;
        this.name = location.name;
        this.population = location.population;
        this.countryCode = location.countryCode;
        this.parentName = location.parentName;
        this.parentPopulation = location.parentPopulation;
        this.geonameID = location.geonameID;
        this.approved = location.approved;
      },

      showGeonameIdInput() {
        if (this.geonameID != null) {
          this.geonameIDInputValue = this.geonameID;
        }
        this.isGeonameIDInputVisible = true;
      },

      hideGeonameIdInput() {
        if (this.geonameID == null) {
          this.geonameIDInputValue = '';
        }
        this.geonameIDError = '';
        this.isGeonameIDInputVisible = false;
      },

      locationApprove(approved) {
        api.locationApprove({
          datasetID: this.datasetID,
          locationHash: this.locationHash,
          approved: approved,
        }).then((location) => {
          this.approved = location.approved;
          if (location.approved === false) {
            this.geonameIDInputValue = '';
            this.setGeonameID();
          }
        }).catch((response) => {
          alert(response);
        });
      },

      getGeoNameFromHierarchy(level) {
        if (!this.geonamesHierarchy.hasOwnProperty('geonames')) {
          return null;
        }
        return this.geonamesHierarchy.geonames[this.geonamesHierarchy.geonames.length - level]
      },

      setGeonameID() {
        this.geonameIDError = '';
        let geonameID = null;
        if (typeof (this.geonameIDInputValue) === 'string') {
          this.geonameIDInputValue = this.geonameIDInputValue.replace(/\s/g, '');
          if (this.geonameIDInputValue !== '' && geonameIDRegex.exec(this.geonameIDInputValue) === null) {
            this.geonameIDError = 'Geoname ID must be a number';
            return
          }
          geonameID = this.geonameIDInputValue !== '' ? parseInt(this.geonameIDInputValue, 10) : null;
        } else {
          geonameID = this.geonameIDInputValue;
        }
        api.locationSetGeonameID({
          datasetID: this.datasetID,
          locationHash: this.locationHash,
          geonameID: geonameID,
        }).then((location) => {
          this.updateData(location);
          this.isGeonameIDInputVisible = false;
        }).catch((response) => {
          this.geonameIDError = response;
        });
      }
    }
  }
</script>
