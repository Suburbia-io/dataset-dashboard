<template>
  <div>
    <div class="field">
      <input class="is-checkradio is-info" id="mode-checkbox" type="checkbox" v-model="isNested">
      <label for="mode-checkbox">Nested display</label>
    </div>
    <div class="table-container" v-if="isFlat">
      <table class="table is-striped is-hoverable is-fullwidth is-narrow">
        <thead>
        <tr>
          <th>Idx</th>
          <th class="tag-name" @click="orderColumn('tag')" style="white-space: nowrap">
            Tag
            <span class="icon is-small">
            <i class="fas" :class="{'fa-caret-up': orderAsc, 'fa-caret-down': !orderAsc}"
               v-show="orderBy==='tag'"></i>
          </span>
          </th>
          <th class="description">Description</th>
          <th class="internal-notes">Internal Notes</th>
          <th class="grade" @click="orderColumn('grade')" style="white-space: nowrap">
            Grd
            <span class="icon is-small">
            <i class="fas" :class="{'fa-caret-up': orderAsc, 'fa-caret-down': !orderAsc}"
               v-show="orderBy==='grade'"></i>
          </span>
          </th>
          <th class="include">Inc</th>
          <th class="num-fingerprints" @click="orderColumn('numFingerprints')">
            Num. <span style="white-space: nowrap;">Fingerprints
            <span class="icon is-small">
              <i class="fas" :class="{'fa-caret-up': orderAsc, 'fa-caret-down': !orderAsc}"
                 v-show="orderBy==='numFingerprints'"></i>
            </span>
          </span>
          </th>
          <th class="num-line-items" @click="orderColumn('numLineItems')">
            Num. Line <span style="white-space: nowrap;">Items
            <span class="icon is-small">
              <i class="fas" :class="{'fa-caret-up': orderAsc, 'fa-caret-down': !orderAsc}"
                 v-show="orderBy==='numLineItems'"></i>
            </span>
          </span>
          </th>
        </tr>
        </thead>
        <tbody>
        <tags-table-row
          v-for="(tag, index) in this.sortedTags"
          :key="tag.tagID"
          :tag="tag"
          :index="index + 1">
        </tags-table-row>
        <tr v-if="this.tags.length === 0">
          <td colspan="7" class="has-text-centered">No results</td>
        </tr>
        </tbody>
        <tfoot>
        </tfoot>
      </table>
    </div>

    <template v-if="isNested">
      <div class="columns is-mobile nested-header">
        <div class="column is-3">Tag</div>
        <div class="column is-3">Description</div>
        <div class="column is-2">Internal Notes</div>
        <div class="column is-0-point-5">Grd</div>
        <div class="column is-0-point-5">Inc</div>
        <div class="column is-1-point-5">Num. Fingerprints</div>
        <div class="column is-1-point-5">Num. Line Items</div>
      </div>

      <tags-nested-table-row
        v-for="node in this.tagNodeTree.children"
        :key="node.tagID"
        :node="node"
        :datasetID="datasetID"
        :tagTypeID="tagTypeID">
      </tags-nested-table-row>

      <div class="columns is-mobile nested-row"
           style="border-bottom: 0;"
           v-if="this.tags.length === 0">
        <div class="column has-text-centered">No results</div>
      </div>
    </template>
  </div>
</template>

<script>
  import TagsTableRow from "./TagsTableRow";
  import TagsNestedTableRow from "./TagsNestedTableRow";
  import TagNode from "../modules/tagnode";

  const DEFAULT_ORDER_BY = 'numLineItems';
  const DEFAULT_ORDER_ASC = false;
  const NESTED_MODE_ORDER_BY = 'tag';
  const NESTED_MODE_ORDER_ASC = true;

  export default {
    name: "TagsTable",
    components: {TagsNestedTableRow, TagsTableRow},
    beforeCreate() {
      this.tagsJSON = JSON.parse(document.getElementById('tags').getAttribute('data-tags'));
      this.datasetID = document.getElementById('dataset-id').getAttribute('data-dataset-id');
      this.tagTypeID = document.getElementById('tag-type-id').getAttribute('data-tag-type-id');
    },
    data() {
      return {
        isNested: false,
        orderBy: DEFAULT_ORDER_BY,
        orderAsc: DEFAULT_ORDER_ASC,
        tags: this.tagsJSON !== null ? this.tagsJSON : [],
        localStorageNestedKey: `${this.datasetID}.${this.tagTypeID}.nested`,
      }
    },
    computed: {
      sortedTags() {
        if (this.tags.length === 0) {
          return this.tags;
        }

        if (typeof (this.tags[0][this.orderBy]) === 'string') {
          return this.tags.sort((a, b) => {
            if (a[this.orderBy] < b[this.orderBy]) {
              return this.orderAsc ? -1 : 1;
            }
            if (a[this.orderBy] > b[this.orderBy]) {
              return this.orderAsc ? 1 : -1;
            }
            return 0;
          });
        } else {
          return this.tags.sort((a, b) => {
            if (this.orderAsc) {
              return a[this.orderBy] - b[this.orderBy];
            } else {
              return b[this.orderBy] - a[this.orderBy];
            }
          });
        }
      },
      tagNodeTree: function () {
        let tree = new TagNode([''], null, null);
        for (const tag of this.sortedTags) {
          tree.insertTagNode(tag);
        }
        return tree;
      },
      isFlat() {
        return !this.isNested;
      },
    },
    created() {
      const localStorageNestedVal = localStorage.getItem(this.localStorageNestedKey);
      if (localStorageNestedVal) {
        this.isNested = localStorageNestedVal === 'true';
      }
    },
    watch: {
      isNested: function (val) {
        this.orderBy = val ? NESTED_MODE_ORDER_BY : DEFAULT_ORDER_BY;
        this.orderAsc = val ? NESTED_MODE_ORDER_ASC : DEFAULT_ORDER_ASC;
        localStorage.setItem(this.localStorageNestedKey, JSON.stringify(val));
      },
    },
    methods: {
      orderColumn(name) {
        this.orderBy = name;
        this.orderAsc = !this.orderAsc;
      }
    },
  }
</script>
