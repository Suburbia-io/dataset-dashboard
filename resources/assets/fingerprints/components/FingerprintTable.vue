<script>
  import api from '../js/api'
  import TextHighlight from 'vue-text-highlight';
  import UpdatableField from "./UpdatableField"
  import BulkUpsertModal from "./BulkUpsertModal"
  import FilterField from "./FilterField"
  import confetti from 'canvas-confetti'

  export default {
    components: {
      'text-highlight': TextHighlight,
      'filterfield': FilterField,
      'updatefield': UpdatableField,
      'bulkupsert': BulkUpsertModal,
    },
    data: () => ({
      datasetID: null,
      tagAppID: '00000000-0000-0000-0000-100000000000',

      // modals
      showUpdateModalForTagType: null,

      // remote data
      fingerprints: [],
      tagTypes: [],
      tagApps: [],

      // filter inputs
      fingerprintIncludes: [],
      fingerprintExcludes: [],
      rawTextIncludes: [],
      rawTextExcludes: [],
      tagIncludes: {},
      tagExcludes: {},
      consTagIncludes: {},
      consTagExcludes: {},

      // order settings
      orderAsc: false,
      orderCol: 'count',

      more: false,
      impact: 0,
      addPercentages: true,
      toggleSettingsDropdown: false,

      bulkExcludes: [],

      // ui settings
      stickyHeader: false,
      countThreshold: 100,
      loading: 0,
      columns: {},
      rawTextWidth: 330,
    }),
    computed: {
      readOnly: function () {
        return this.tagAppID !== '00000000-0000-0000-0000-100000000000'
      }
    },
    created() {
      this.datasetID = window.datasetID
      this.getApps()
      this.getList()
      window.onscroll = this.handleScroll

      const localStorageColumnsKey = this.datasetID + '.columns'
      const localStorageColumnsVal = localStorage.getItem(localStorageColumnsKey)
      if (localStorageColumnsVal) {
        this.columns = JSON.parse(localStorageColumnsVal);
      }

      const localStorageWidthKey = this.datasetID + '.rawTextWidth'
      const localStorageWidthVal = localStorage.getItem(localStorageWidthKey)
      if (localStorageWidthVal) {
        this.rawTextWidth = localStorageWidthVal;
      }

      const localStorageCountThresholdKey = this.datasetID + '.countTreshold'
      const localStorageCountThresholdVal = localStorage.getItem(localStorageCountThresholdKey)
      if (localStorageCountThresholdVal) {
        this.countThreshold = parseInt(localStorageCountThresholdVal);
      }
    },
    watch: {
      columns: {
        handler(val) {
          const localStorageKey = this.datasetID + '.columns'
          localStorage.setItem(localStorageKey, JSON.stringify(val));
        },
        deep: true
      },
      rawTextWidth(val) {
        const localStorageKey = this.datasetID + '.rawTextWidth'
        localStorage.setItem(localStorageKey, val);
      },
      countThreshold(val) {
        const localStorageKey = this.datasetID + '.countTreshold'
        localStorage.setItem(localStorageKey, val);
        this.countThreshold = parseInt(val)
      }
    },
    methods: {
      toggleAddPercentages(e) {
        const evt = window.event ? event : e
        if (evt.keyCode === 13 && evt.shiftKey) {
          this.addPercentages = !this.addPercentages
        }
      },
      handleScroll() {
        const scrollPos = window.pageYOffset || document.documentElement.scrollTop
        this.stickyHeader = scrollPos >= 149;
      },
      getTextsToHighlight(field) {
        if (!field) {
          return []
        }
        const l = []
        field.forEach(function (part, index, arr) {
          if (part.length > 0 && part[0] === '%') {
            part = part.slice(1)
          }
          if (part.length > 0 && part[part.length - 1] === '%') {
            part = part.slice(0, part.length - 1)
          }
          l.push(part)
        });
        return l
      },
      celebrateGoodTimes() {
        const duration = 800;
        const end = Date.now() + duration;
        const interval = setInterval(function () {
          if (Date.now() > end) {
            return clearInterval(interval);
          }

          confetti({
            startVelocity: 30,
            spread: 100,
            ticks: 60,
            shapes: ['square'],
            origin: {
              x: Math.random(),
              // since they fall down, start a bit higher than random
              y: Math.random() - 0.2
            }
          });
        }, 80);
      },
      getList() {
        this.bulkExcludes = []
        this.loading++
        api.fingerprintList({
          datasetID: this.datasetID,
          tagAppID: this.tagAppID,
          limit: 35,
          offset: 0,
          countThreshold: this.countThreshold,
          orderBy: this.orderCol,
          orderAsc: this.orderAsc,
          fingerprintIncludes: this.fingerprintIncludes,
          fingerprintExcludes: this.fingerprintExcludes,
          consTagIncludes: this.consTagIncludes,
          consTagExcludes: this.consTagExcludes,
          rawTextIncludes: this.rawTextIncludes,
          rawTextExcludes: this.rawTextExcludes,
          tagIncludes: this.tagIncludes,
          tagExcludes: this.tagExcludes,
        }).then((resp) => {
          this.tagTypes = resp.tagTypes
          this.fingerprints = resp.rows
          this.more = resp.more
          this.impact = resp.totalCount
          this.loading--
        }).catch((err) => {
          this.loading--
          alert("Something went wrong: " + err)
        })
      },
      getNext() {
        this.loading++
        api.fingerprintList({
          datasetID: this.datasetID,
          tagAppID: this.tagAppID,
          limit: 35,
          countThreshold: this.countThreshold,
          offset: this.fingerprints.length,
          orderBy: this.orderCol,
          orderAsc: this.orderAsc,
          fingerprintIncludes: this.fingerprintIncludes,
          fingerprintExcludes: this.fingerprintExcludes,
          consTagIncludes: this.consTagIncludes,
          consTagExcludes: this.consTagExcludes,
          rawTextIncludes: this.rawTextIncludes,
          rawTextExcludes: this.rawTextExcludes,
          tagIncludes: this.tagIncludes,
          tagExcludes: this.tagExcludes,
        }).then((resp) => {
          if (resp.rows) {
            this.fingerprints.push(...resp.rows)
          }
          this.loading--
        }).catch((err) => {
          this.loading--
          alert("Something went wrong: " + err)
        })
      },
      getTagTypes() {
        this.loading++
        api.datasetTagTypeList({
          datasetID: this.datasetID,
        }).then((tagTypes) => {
          this.tagTypes = tagTypes
          this.loading--
        })
      },
      getApps() {
        this.loading++
        api.fingerprintsListTagApps()
          .then((apps) => {
            this.tagApps = apps
            this.loading--
          })
      },
      findTagValueForFingerprint(tagTypeID, fingerprint) {
        if (null === fingerprint.fingerprintTags) {
          return ''
        }
        for (let i = 0; i < fingerprint.fingerprintTags.length; i++) {
          if (fingerprint.fingerprintTags[i].tagTypeID === tagTypeID) {
            return fingerprint.fingerprintTags[i].tag
          }
        }
        return ''
      },
      updateSingleTag(fingerprint, tagTypeID, tagValue) {
        this.loading++
        return api.fingerprintUpsertTags({
          datasetID: this.datasetID,
          tagTypeID: tagTypeID,
          tagValue: tagValue,
          fingerprints: [fingerprint]
        }).then((data) => {
          let count = 0
          this.fingerprints.forEach((fp, idx) => {
            if (fp.fingerprint === fingerprint) {
              this.fingerprints[idx] = data[0]
              this.$set(this.fingerprints, idx, data[0])
              if (tagValue !== null) {
                count = fp.count
              }
            }
          })
          this.loading--

          if (count > 100000) {
            this.celebrateGoodTimes()
          }
        })
      },
      bulkUpdateTag(tagTypeID, tagValue, force) {
        if (!tagValue) {
          return
        }
        this.loading++
        let fingerprints = []

        let tagIdx = null

        for (let i = 0; i < this.tagTypes.length; i++) {
          const tt = this.tagTypes[i]
          if (tt.tagTypeID === tagTypeID) {
            tagIdx = i
            break
          }
        }

        let count = 0

        this.fingerprints.forEach((fp, fpIdx) => {
          let tagAlreadySet = true
          if (!fp.tagIDs[tagIdx]) {
            tagAlreadySet = false
          }

          const excluded = (this.bulkExcludes.indexOf(fpIdx) !== -1)

          if (excluded) {
            return
          }

          if (!tagAlreadySet || force) {
            count += fp.count
            fingerprints.push(fp.fingerprint)
          }
        })
        api.fingerprintUpsertTags({
          datasetID: this.datasetID,
          tagTypeID: tagTypeID,
          tagValue: tagValue,
          fingerprints: fingerprints,
        }).then((resp) => {
          this.fingerprints.forEach((fp, idx) => {
            resp.forEach((rfp) => {
              if (fp.fingerprint === rfp.fingerprint) {
                console.log(fp.fingerprint, rfp.fingerprint)
                this.fingerprints[idx] = rfp
                this.$set(this.fingerprints, idx, rfp)
              }
            })
          })
          this.loading--

          if (count > 100000) {
            this.celebrateGoodTimes()
          }
        })
      },
      getTagSuggestions(tagTypeID, searchValue) {
        return api.fingerprintTagSuggestions({
          datasetID: this.datasetID,
          tagTypeID: tagTypeID,
          search: searchValue
        })
      },
      keyMove(direction, fingerprintIdx, columnIdx) {
        if (direction === 'up') {
          fingerprintIdx--
        } else if (direction === 'down') {
          fingerprintIdx++
        } else if (direction === 'left') {
          columnIdx--
        } else if (direction === 'right') {
          columnIdx++
        }
        const ref = this.$refs[fingerprintIdx + ':' + columnIdx]
        if (ref && ref[0]) {
          if (ref[0].isReadOnly()) {
            this.keyMove(direction, fingerprintIdx, columnIdx)
            return;
          }
          ref[0].doFocus();
        }
      },
      bulkExcludeCheck(checked, fpIdx) {
        if (checked) {
          this.bulkExcludes.push(fpIdx)
          return
        }
        const index = this.bulkExcludes.indexOf(fpIdx);
        if (index > -1) {
          this.bulkExcludes.splice(index, 1);
        }
      }
    }
  }
</script>

<template>
  <div class="Fingerprinting">
    <div class="loading-overlay" v-show="loading>0"></div>
    <div class="Toolbar">
      <div class="select is-small">
        <select v-model="orderCol" @change="getList()">
          <option value="fingerprint">Fingerprints</option>
          <option value="count">Count</option>
          <template v-for="tagType in tagTypes">
            <option :value="tagType.tagType + '_cons_confidence'">{{tagType.tagType }} cons confidence</option>
          </template>
        </select>
      </div>
      <div class="select is-small">
        <select v-model="orderAsc" @change="getList">
          <option :value="true">Asc</option>
          <option :value="false">Desc</option>
        </select>
      </div>
      <button class="button is-small" @click="getList"
              content="Refresh content while keeping filters & sort options intact"
              v-tippy="{ theme : 'dark', interactive: false, duration:0}"><i class="fas fa-sync"></i></button>
      <div class="select is-small">
        <select v-model="addPercentages">
          <option :value="true">Add wildcards to filters</option>
          <option :value="false">Don't add wildcards to filters</option>
        </select>
      </div>
      <div class="dropdown" :class="{'is-active': toggleSettingsDropdown}">
        <div class="dropdown-trigger is-small">
          <button class="button is-small" aria-haspopup="true" aria-controls="dropdown-menu"
                  @click="toggleSettingsDropdown = !toggleSettingsDropdown">
            <span><i class="fas fa-hammer"></i> Settings</span>
            <span class="icon is-small">
                <i class="fas fa-angle-down" aria-hidden="true"></i>
              </span>
          </button>
        </div>
        <div class="dropdown-menu" id="dropdown-menu" role="menu">
          <div class="dropdown-content" style="width: 250px;">
            <div class="dropdown-item">
              <label class="checkbox">
                <input type="checkbox" @change="$set(columns, 'fingerprint', !columns['fingerprint'])"
                       :checked="columns['.fingerprint']">
                Show fingerprint
              </label>
            </div>
            <div class="dropdown-item">
              <label class="checkbox">
                <input type="checkbox" @change="$set(columns, 'annotations', !columns['annotations'])"
                       :checked="columns['.annotations']">
                Show annotations
              </label>
            </div>
            <hr style="margin:0">
            <div v-for="(tagType, idx) in tagTypes">
              <div class="dropdown-item">
                <label class="checkbox">
                  <input type="checkbox" @change="$set(columns, idx + '.input', !columns[idx + '.input'])"
                         :checked="columns[idx + '.input']">
                  Show {{ tagType.tagType }} input
                </label>
              </div>
              <div class="dropdown-item">
                <label class="checkbox">
                  <input type="checkbox" @change="$set(columns, idx + '.updated', !columns[idx + '.updated'])"
                         :checked="columns[idx + '.updated']">
                  Show {{ tagType.tagType }} updated
                </label>
              </div>
              <div class="dropdown-item">
                <label class="checkbox">
                  <input type="checkbox" @change="$set(columns, idx + '.consensus', !columns[idx + '.consensus'])"
                         :checked="columns[idx + '.consensus']">
                  Show {{ tagType.tagType }} consensus
                </label>
              </div>
              <div class="dropdown-item">
                <label class="checkbox">
                  <input type="checkbox" @change="$set(columns, idx + '.confidence', !columns[idx + '.confidence'])"
                         :checked="columns[idx + '.confidence']">
                  Show {{ tagType.tagType }} confidence
                </label>
              </div>
              <hr style="margin:0">
            </div>
          </div>
        </div>
        <div class="field" style="margin:0 6px">
          <div class="control has-icons-left" content="Width of the Raw Text field"
               v-tippy="{ theme : 'dark', interactive: false, duration:0}">
            <input class="input is-small" type="number" step="10" style="width: 85px;" v-model="rawTextWidth" min="320">
            <span class="icon is-small is-left">
                <i class="fas fa-text-width"></i>
              </span>
          </div>
        </div>
        <div class="field" style="margin:0 6px">
          <div class="control has-icons-left" content="Count threshold"
               v-tippy="{ theme : 'dark', interactive: false, duration:0}">
            <input class="input is-small" type="number" step="1" style="width: 90px;" v-model="countThreshold" min="0">
            <span class="icon is-small is-left">
                <i class="fas fa-filter"></i>
              </span>
          </div>
        </div>
        <div class="select is-small">
          <select v-model="tagAppID" @change="getList">
            <option :value="app.tagAppID" v-for="app in tagApps">{{ app.name }}</option>
          </select>
        </div>
      </div>
    </div>

    <bulkupsert v-if="showUpdateModalForTagType !== null"
                :close="() => {this.showUpdateModalForTagType = null}"
                :callback="(newValue, force) => {return bulkUpdateTag(showUpdateModalForTagType, newValue, force)}"
                :ignorecount="bulkExcludes.length"
                :get-suggestions="(typingValue) => {return getTagSuggestions(showUpdateModalForTagType, typingValue)}"></bulkupsert>

    <table class="Spreadsheet" :class="{'has-sticky-header': stickyHeader}">
      <thead :class="{'sticky-header': stickyHeader}">
      <tr>
        <th class="medium has-maxwidth" v-if="columns['fingerprint']">Fingerprint</th>
        <th class="small" align="right" :content="'Total impact of filtered results: ' + impact"
            v-tippy="{ theme : 'dark', interactive: true, duration:0}">Count <i class="fas fa-info-circle"></i></th>
        <th class="large has-maxwidth"
            :style="'max-width: ' + rawTextWidth + 'px; width: ' + rawTextWidth + 'px; min-width: ' + rawTextWidth + 'px;' ">
          Raw Text
        </th>
        <th class="large has-maxwidth" v-if="columns['annotations']">Annotations</th>
        <th class="micro" content="Exclude fingerprint when bulk tagging"
            v-tippy="{ theme : 'dark', interactive: false, duration:0}"><i class="fas fa-ban"></i></th>
        <template v-for="(tagType,idx) in tagTypes">
          <th @click="showUpdateModalForTagType = tagType.tagTypeID"
              class="clickable" v-if="columns[idx+'.input']">
            {{ tagType.tagType }}
            <div class="tag is-link is-normal"
                 style="font-size: 10px; height: auto; line-height: 16px; padding-left: 5px; padding-right: 5px;">Tag
            </div>
          </th>

          <th class="medium" v-if="columns[idx+'.updated']">
            {{ tagType.tagType }}
            <div class="tag"
                 style="font-size: 10px; height: auto; line-height: 16px; padding-left: 5px; padding-right: 5px;">
              Updated
            </div>
          </th>

          <th class="medium" v-if="columns[idx+'.consensus']">
            {{ tagType.tagType }}
            <div class="tag"
                 style="font-size: 10px; height: auto; line-height: 16px; padding-left: 5px; padding-right: 5px;">
              Consensus
            </div>
          </th>

          <th class="medium" v-if="columns[idx+'.confidence']">
            {{ tagType.tagType }}
            <div class="tag"
                 style="font-size: 10px; height: auto; line-height: 16px; padding-left: 5px; padding-right: 5px;">
              Confid.
            </div>
          </th>

        </template>

        <th class="fill"></th>
      </tr>

      <tr class="filters-row">
        <td v-if="columns['fingerprint']">
          <!-- Fingerprint ix -->
          <filterfield
            :callback="(includes, excludes) => {this.fingerprintIncludes = includes; this.fingerprintExcludes = excludes; getList()}"
            :percentages="addPercentages"></filterfield>
        </td>
        <td></td>
        <td :colspan="(columns['annotations']) ? 2 : 1">
          <!-- RawText ix -->
          <filterfield
            :callback="(includes, excludes) => {this.rawTextIncludes = includes; this.rawTextExcludes = excludes; getList()}"
            :percentages="addPercentages"></filterfield>
        </td>

        <th class="micro"></th>

        <template v-for="(tagType,idx) in tagTypes">
          <td class="medium" v-if="columns[idx+'.input']">
            <!-- Tag Includes -->
            <filterfield
              :callback="(includes, excludes) => {tagIncludes[tagType.tagType] = includes; tagExcludes[tagType.tagType] = excludes; getList()}"
              :percentages="addPercentages"></filterfield>
          </td>

          <td class="medium" v-if="columns[idx+'.updated']">

          </td>

          <td class="medium" v-if="columns[idx+'.consensus']">
            <!-- Consensus Includes -->
            <filterfield
              :callback="(includes, excludes) => {consTagIncludes[tagType.tagType] = includes; consTagExcludes[tagType.tagType] = excludes; getList()}"
              :percentages="addPercentages"></filterfield>
          </td>

          <td class="medium" v-if="columns[idx+'.confidence']">

          </td>

        </template>

        <td class="fill"></td>
      </tr>

      </thead>
      <tbody>
      <tr v-for="(fp,fpIdx) in fingerprints" :key="fp.fingerprintID"
          :class="{excluded: (bulkExcludes.indexOf(fpIdx) !== -1)}">
        <td class="medium has-maxwidth clickable" :class="{excluded: (bulkExcludes.indexOf(fpIdx) !== -1)}"
            v-if="columns['fingerprint']" :content="fp.fingerprint"
            v-tippy="{ theme : 'dark', trigger : 'click', interactive: true, duration:0}">
          <text-highlight :queries="getTextsToHighlight(fingerprintIncludes)">{{ fp.fingerprint }}</text-highlight>
        </td>
        <td class="small" style="text-align: right;" :class="{excluded: (bulkExcludes.indexOf(fpIdx) !== -1)}"><span>{{ fp.count }}</span>
        </td>
        <td class="large has-maxwidth clickable" :class="{excluded: (bulkExcludes.indexOf(fpIdx) !== -1)}"
            :style="'max-width: ' + rawTextWidth + 'px; width: ' + rawTextWidth + 'px; min-width: ' + rawTextWidth + 'px;'"
            :content="fp.rawText" v-tippy="{ theme : 'dark', trigger : 'click', interactive: true, duration:0}">
          <text-highlight :queries="getTextsToHighlight(rawTextIncludes)">{{ fp.rawText }}</text-highlight>
        </td>
        <td class="large has-maxwidth clickable" :class="{excluded: (bulkExcludes.indexOf(fpIdx) !== -1)}"
            v-if="columns['annotations']" :content="fp.annotations"
            v-tippy="{ theme : 'dark', trigger : 'click', interactive: true, duration:0}">
          <text-highlight :queries="getTextsToHighlight(rawTextIncludes)">{{ fp.annotations }}</text-highlight>
        </td>

        <td class="micro">
          <input type="checkbox" v-model="bulkExcludes" :value="fpIdx" style="margin:4px auto; display: block;" number>
        </td>

        <template v-for="(tagType,idx) in tagTypes">

          <td class="medium" v-if="columns[idx+'.input']">
            <!--<td v-for="tagType in tagTypes" :key="tagType.tagTypeID" class="medium">  -->
            <updatefield
              :ref="fpIdx + ':' + idx"
              @up="keyMove('up', fpIdx, idx)"
              @down="keyMove('down', fpIdx, idx)"
              @left="keyMove('left', fpIdx, idx)"
              @right="keyMove('right', fpIdx, idx)"
              :highlights="getTextsToHighlight(tagIncludes[tagType.tagType])"
              :value="fp.tags[idx]"
              :callback="(newValue) => {return updateSingleTag(fp.fingerprint, tagType.tagTypeID, newValue)}"
              :get-suggestions="(typingValue) => {return getTagSuggestions(tagType.tagTypeID, typingValue)}"
              :readonly="readOnly || (bulkExcludes.indexOf(fpIdx) !== -1)">
            </updatefield>
          </td>

          <td class="medium" v-if="columns[idx+'.updated']">
            <span>{{ fp.updatedAtsFormatted[idx] }}</span>
          </td>

          <td class="medium" v-if="columns[idx+'.consensus']">
            <span>{{ fp.consTags[idx] }}</span>
          </td>

          <td class="medium" align="right" v-if="columns[idx+'.confidence']">
            <span v-if="fp.consConfidences[idx]">{{ fp.consConfidences[idx].toFixed(2) }}</span>
          </td>


        </template>
        <td class="fill"></td>
      </tr>
      </tbody>
    </table>

    <div class="Toolbar">
      <div class="button is-small" @click="getNext">More</div>
    </div>

  </div>
</template>
