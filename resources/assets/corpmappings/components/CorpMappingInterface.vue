<template>
  <div>

    <select v-if="tagTypes" v-model="selectedTagType">
      <option v-for="tagType in tagTypes" :value="tagType.tagTypeID">{{ tagType.tagType }}
      </option>
    </select>

    <select v-if="tagValues" v-model="selectedTagValue" :disabled="tagValues.length === 0">
      <option v-for="tagValue in tagValues" :value="tagValue.tagID">{{ tagValue.tag }}
      </option>
    </select>

    <button @click="insertMapping">Insert Mapping</button>

    <hr>

    <div>

      <nav class="panel is-size-7" style="border: 1px solid #ccc; max-width: 60%; width: 60%;"
           v-for="corpMapping in corpMappings">
        <div class="panel-block">
          <span class="tag is-info">{{ corpMapping.TagType.tagType }}</span>
          <span class="tag is-white">=</span>
          <span class="tag is-success">{{ corpMapping.Tag.tag }}</span>
          <div class="tag-stats is-pulled-right">
            <span>Num Fingerprints: {{ corpMapping.Tag.numFingerprints }}</span>
            <span>Num Line Items: {{ corpMapping.Tag.numLineItems }}</span>
            <span>Included: <template v-if="corpMapping.Tag.isIncluded">Yes</template><template v-else>No</template></span>
          </div>
        </div>

        <table class="table is-narrow is-hoverable is-fullwidth">
          <thead>
          <tr>
            <th>Start Date</th>
            <th>Country</th>
            <th>Corporation</th>
            <th>&nbsp;</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="rule in corpMapping.rules">
            <td>{{ rule.fromDate }}</td>
            <td v-if="rule.country !== ''">{{ rule.country }}</td>
            <td v-if="rule.country === ''">*</td>
            <td>{{ rule.Corporation.name }} ({{ rule.Corporation.exchange }}:{{ rule.Corporation.code }})</td>
            <td>
              <button class="button is-danger is-small" @click="removeCorpMappingRule(rule.corpMappingRuleID)">
                <i class="fas fa-trash"></i>
              </button>
            </td>
          </tr>
          </tbody>
          <tfoot>
          <corpmapping-rule-form :corps="corps" :mapping-id="corpMapping.corpMappingID"
                                 :callback="insertMappingRule"></corpmapping-rule-form>
          </tfoot>
        </table>

      </nav>

    </div>
  </div>
</template>

<script>
  import api from '../js/api'
  import CorpMappingRuleForm from "./CorpMappingRuleForm";

  export default {
    components: {
      'corpmapping-rule-form': CorpMappingRuleForm,
    },
    data: () => ({
      tagTypes: [],
      datasetID: null,
      corpTypeID: null,
      tagValues: [],
      selectedTagType: '',
      selectedTagValue: '',
      corpMappings: [],
      corps: [],
    }),
    mounted() {
      this.corpTypeID = window.corpmappings.corpTypeID;
      this.datasetID = window.corpmappings.datasetID;
      this.tagTypes = window.corpmappings.tagTypes;
      this.corpMappings = window.corpmappings.corpmappings;
      this.corps = window.corpmappings.corps;
    },
    watch: {
      selectedTagType: function () {
        this.getTagValues()
        this.selectedTagValue = ''
      }
    },
    methods: {
      getTagValues() {
        if (!this.selectedTagType) {
          return
        }
        api.listTagsForTagType({DatasetID: this.datasetID, TagTypeID: this.selectedTagType}).then((res) => {
          this.tagValues = res
        })
      },
      insertMapping() {
        api.insertCorpMapping({
          corpTypeID: this.corpTypeID,
          tagTypeID: this.selectedTagType,
          tagID: this.selectedTagValue
        }).then((res) => {
          location.reload()
        }).catch((err) => {
          alert(err)
        })
      },
      insertMappingRule(mappingID, startDate, country, corpID) {
        api.insertCorpMappingRule({
          corpMappingID: mappingID,
          fromDate: startDate,
          country: country,
          corpID: corpID
        }).then((res) => {
          location.reload()
        }).catch((err) => {
          alert(err)
        })
      },
      removeCorpMappingRule(ruleID) {
        if (confirm('Are you sure you want to remove this rule?')) {
          api.deleteCorpMappingRule({corpMappingRuleID: ruleID})
            .then((res) => {
              location.reload()
            })
        }
      }
    }
  }
</script>

<style scoped>

</style>
