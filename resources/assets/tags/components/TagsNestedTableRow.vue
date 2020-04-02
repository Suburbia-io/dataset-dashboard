<template>
  <div>
    <div class="columns is-mobile nested-row" :class="{'synthetic-row': isSynthetic}">
      <div class="column is-3" :style="indent" @click="toggleChildren"
           :class="{'has-children': hasChildren}">
        <span class="icon" v-if="hasChildren">
          <i class="far" :class="{'fa-folder': !showChildren, 'fa-folder-open': showChildren}"></i>
        </span>
        {{node.name}}
        <a :href="tagEditUrl"
           class="button is-small is-pulled-right"
           @click.stop
           v-if="!isSynthetic">
          Edit
        </a>
      </div>
      <div class="column is-3 description">{{tag.description}}</div>
      <div class="column is-2 internal-notes">{{tag.internalNotes}}</div>
      <div class="column is-0-point-5 grade ">{{tag.grade}}</div>
      <div class="column is-0-point-5 include">
        <template v-if="tag.isIncluded===true">Yes</template>
        <template v-else-if="tag.isIncluded===false">No</template>
        <template v-else>-</template>
      </div>
      <div class="column is-1-point-5 num-fingerprints">
        <template v-if="!isSynthetic">
          {{totalNumFingerprints}}
          <template v-if="hasChildren">({{numFingerprints}})</template>
        </template>
        <template v-else>-</template>
      </div>
      <div class="column is-1-point-5 num-line-items">
        <template v-if="!isSynthetic">
          {{totalNumLineItems}}
          <template v-if="hasChildren">({{numLineItems}})</template>
        </template>
        <template v-else>-</template>
      </div>
    </div>
    <tags-nested-table-row
      v-if="showChildren"
      v-for="node in node.children"
      :key="node.tagID"
      :node="node"
      :datasetID="datasetID"
      :tagTypeID="tagTypeID">
    </tags-nested-table-row>
  </div>
</template>

<script>
  const nf = new Intl.NumberFormat();

  export default {
    name: "TagsNestedTableRow",
    props: ['node', 'datasetID', 'tagTypeID'],
    data() {
      return {
        showChildren: false,
      }
    },
    created() {
      if (this.isSynthetic) {
        this.showChildren = true;
      }
    },
    computed: {
      tagEditUrl() {
        return `/admin/datasets/${this.datasetID}/tag-types/${this.tagTypeID}/tags/${this.node.tagID}/`;
      },
      indent() {
        let padding = (this.node.depth - 1) * 2.5;
        return `padding-left: calc(${padding}rem + 0.5rem);`;
      },
      hasChildren() {
        return this.node.children.length > 0;
      },
      isSynthetic() {
        return this.node.tag === null;
      },
      tag() {
        return this.node.tag === null ? this.node.path.join('/') : this.node.tag;
      },
      description() {
        return this.node.tag === null ? '-' : this.node.tag.description;
      },
      internalNotes() {
        return this.node.tag === null ? '-' : this.node.tag.internalNotes;
      },
      grade() {
        return this.node.tag === null ? '-' : this.node.tag.grade;
      },
      isIncluded() {
        return this.node.tag === null ? '-' : this.node.tag.isIncluded;
      },
      numFingerprints() {
        return this.node.tag === null ? '-' : nf.format(this.node.tag.numFingerprints);
      },
      numLineItems() {
        return this.node.tag === null ? '-' : nf.format(this.node.tag.numLineItems);
      },
      totalNumFingerprints() {
        return nf.format(this.node.totalNumFingerprints);
      },
      totalNumLineItems() {
        return nf.format(this.node.totalNumLineItems);
      },
    },
    methods: {
      toggleChildren() {
        this.showChildren = !this.showChildren;
      },
    }
  }
</script>
