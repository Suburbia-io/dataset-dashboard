<script>
    import VueTagsInput from '@johmun/vue-tags-input';

    export default {
        components: {
            VueTagsInput,
        },
        props: ['callback', 'percentages'],
        data: () => ({
            tag: '',
            tags: [],
            validation: [{
                classes: 'exclude',
                rule: tag => tag.text.charAt(0) === '!',
            },
            {
                classes: 'include',
                rule: tag => tag.text.charAt(0) !== '!',
            }],
            inverseBehaviour: false ,
        }),
        mounted() {
            window.addEventListener("keyup", this.keyUp)
            window.addEventListener("keydown", this.keyDown)
        },
        destroyed() {
            window.removeEventListener("keyup", this.keyUp)
            window.removeEventListener("keydown", this.keyDown)
        },
        methods: {
            keyDown(e) {
                const evt = window.event ? event : e
                if (evt.keyCode === 16) {
                    this.inverseBehaviour = true;
                }
            },
            keyUp(e) {
                const evt = window.event ? event : e
                if (evt.keyCode === 16) {
                    this.inverseBehaviour = false;
                }
            },
            addPercentages(tagObj) {
                if (this.percentages && !this.inverseBehaviour || !this.percentages && this.inverseBehaviour) {
                      if (!tagObj.tag.text.includes('NULL')) {
                          if (tagObj.tag.text.charAt(0) === '!') {
                              tagObj.tag.text = '!%' + tagObj.tag.text.substr(1).toLowerCase() + '%'
                          } else {
                              tagObj.tag.text = '%' + tagObj.tag.text.toLowerCase() + '%'
                          }
                      }
                }
                tagObj.addTag();
            }
        },
        watch: {
            tags() {
              let includeValues = []
              let excludeValues = []
              this.tags.forEach((tag, idx) => {
                let tagText = tag.text;
                if (tag.text.charAt(0) === '!') {
                  tagText = tag.text.substr(1)
                  excludeValues.push(tagText)
                } else {
                  includeValues.push(tagText)
                }
              })
              this.callback(includeValues,excludeValues)
            }
        }
    }
</script>

<template>
  <div>
    <vue-tags-input
      v-model="tag"
      :validation="validation"
      :tags="tags"
      :separators="['||']"
      @before-adding-tag="addPercentages"
      @tags-changed="newTags => tags = newTags"
      placeholder="Add filter..."
    />
  </div>
</template>

<style lang="scss">
  .vue-tags-input {
    max-width:100% !important;
    width: 100%;
    .ti-input {
      border: 0 !important;
    }

    .ti-tags .ti-tag {
      position: relative;
      margin-right: 2px;
      color: #fff;

      &.exclude {
        background: hsl(348, 100%, 61%)
      }

      &.include {
        background: hsl(141, 53%, 53%);
      }
    }
  }
</style>
