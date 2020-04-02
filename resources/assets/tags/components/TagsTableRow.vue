<template>
  <tr>
    <td>
      {{index}}
    </td>
    <td>
      <a :href="tagEditUrl">{{tag.tag}}</a>
    </td>
    <td>{{tag.description}}</td>
    <td>{{tag.internalNotes}}</td>
    <td>{{tag.grade}}</td>
    <td>
      <template v-if="tag.isIncluded">Yes</template>
      <template v-else>No</template>
    </td>
    <td class="num-fingerprints">{{numFingerprints}}</td>
    <td class="num-line-items">{{numLineItems}}</td>
  </tr>
</template>

<script>
  const nf = new Intl.NumberFormat();

  export default {
    name: "TagsTableRow",
    props: ['tag', 'index'],
    computed: {
      tagEditUrl() {
        return `/admin/datasets/${this.tag.datasetID}/tag-types/${this.tag.tagTypeID}/tags/${this.tag.tagID}/`;
      },
      numFingerprints() {
        return nf.format(this.tag.numFingerprints);
      },
      numLineItems() {
        return nf.format(this.tag.numLineItems);
      },
    },
  }
</script>
