<script>
    export default {
      props: ['getSuggestions', 'callback', 'ignorecount', 'close'],
      data: () => ({
        inputModel: '',
        forceModel: false,
        suggestions: [],
        selectedSuggestionIdx: -1
      }),
      mounted() {
          this.$nextTick(() => {
              this.$refs.input.focus()
          })
      },
      methods: {
          update() {
              let newValue = this.inputModel
              if (this.selectedSuggestionIdx !== -1) {
                  newValue = this.suggestions[this.selectedSuggestionIdx]
              }

              if (this.suggestions.length === 0 || (this.selectedSuggestionIdx < 1 && this.suggestions[0] !== newValue)) {
                  return
              }

              this.callback(newValue, this.forceModel)
              this.close()
          },
          suggestionMoveUp() {
              if (this.selectedSuggestionIdx === (this.suggestions.length - 1)) {
                  return
              }
              this.selectedSuggestionIdx++
          },
          suggestionMoveDown() {
              if (this.selectedSuggestionIdx === -1) {
                  return
              }
              this.selectedSuggestionIdx--
          },
      },
        watch: {
          inputModel(newValue, oldValue) {
              if (newValue.length < 1 || oldValue === '') {
                  this.suggestions = []
                  return
              }
              this.getSuggestions(newValue).then((suggestions) => {
                  this.selectedSuggestionIdx = -1
                  if (null === suggestions) {
                      this.suggestions = []
                      return
                  }
                  this.suggestions = suggestions
                  if (this.suggestions.length === 1) {
                      this.selectedSuggestionIdx = 0
                  }
              }).catch(() => {
                  this.suggestions = []
              })
          }
      }
    }
</script>

<template>
  <div class="modal is-active" @keyup.esc="close">
    <div class="modal-background" @click="close"></div>
    <div class="modal-content">
      <div class="box">
        <label class="label">Set all tag values at once</label>

        <label class="checkbox has-text-danger has-text-weight-bold">
          <input type="checkbox" v-model="forceModel">
          Also overwrite fields that already have a value.
        </label>

        <hr>

        <template v-if="ignorecount > 0">
          <article class="message is-warning">
            <div class="message-body">
              You have <strong>excluded {{ ignorecount }} fingerprint rows</strong>. These will not be updated by this bulk action.
            </div>
          </article>
          <hr>
        </template>

        <input class="input" type="text" ref="input" v-model="inputModel" placeholder="Start typing to search for values..."
               @keyup.down="suggestionMoveUp"
               @keyup.up="suggestionMoveDown"
               @keyup.esc="close"
               @keyup.enter="update">
        <ul class="menu-list">
          <li v-for="(s,i) in suggestions" @click="update" @mouseover="selectedSuggestionIdx = i">
            <a :class="{'is-active': (i === selectedSuggestionIdx)}">{{ s }}</a>
          </li>
        </ul>
      </div>
    </div>
    <button class="modal-close is-large" aria-label="close" @click="close"></button>
  </div>
</template>
