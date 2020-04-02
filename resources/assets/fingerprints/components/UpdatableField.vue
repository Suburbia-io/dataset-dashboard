<script>
    import TextHighlight from 'vue-text-highlight';

    export default {
      components: {
          'text-highlight': TextHighlight,
      },
      props: ['value', 'callback', 'getSuggestions', 'highlights', 'readonly'],
      data: () => ({
        error: false,
        inputModel: '',
        focused: false,
        suggestions: [],
        selectedSuggestionIdx: -1,
        valueChanged: true,
        clickInDropdown: false,
      }),
      mounted() {
        this.inputModel = this.value
        this.valueChanged = false
      },
      methods: {
          isReadOnly() {
            return this.readonly
          },
          blur() {
            this.selectedSuggestionIdx = -1
            this.suggestions = []
            this.inputModel = this.value
            this.error = false
            this.focused = false
          },
          doFocus() {
            if (this.readonly) {
                return;
            }
            this.focused = true
            this.$nextTick(() => {
                this.$refs.input.focus()
                this.$refs.input.select()
            })
          },
          update() {
              if (!this.focused) {
                  return
              }

              let newValue = this.inputModel
              if (this.selectedSuggestionIdx !== -1) {
                  newValue = this.suggestions[this.selectedSuggestionIdx]
              }

              if (newValue === this.value) {
                this.blur()
                this.doFocus()
                return
              }

              if (this.suggestions.length === 0 || (this.selectedSuggestionIdx < 1 && this.suggestions[0] !== newValue)) {
                  this.blur()
                  this.doFocus()
                  return
              }

              this.callback(newValue).then(() => {
                this.value = newValue;
                this.$nextTick(() => {
                  this.blur()
                  this.doFocus()
                });
              }).catch(() => {
                  this.suggestions = []
                  this.focused = true
                  this.error = true
              })
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
          clear() {
              this.callback(null).then(() => {
                  this.blur()
              }).catch(() => {
                  this.focused = true
                  this.error = true
              })
          },
          keyUpPressed() {
            if (this.suggestions.length > 0) {
              this.suggestionMoveDown()
              return;
            }
            this.$emit('up')
          },
          keyDownPressed() {
            if (this.suggestions.length > 0) {
              this.suggestionMoveUp()
              return;
            }
            this.$emit('down')
          },
          keyLeftPressed(e) {
            e.preventDefault()
            this.$emit('left')
          },
          keyRightPressed(e) {
            e.preventDefault()
            this.$emit('right')
          },
        clickSuggestion(s) {
          this.clickInDropdown = true;
          this.$nextTick(() => {
            this.inputModel = s
            this.update()
            this.blur()
          });
        }
      },
      watch: {
        value(newValue) {
            this.valueChanged = true
            this.inputModel = newValue
            this.suggestions = []
        },
        inputModel(newValue, oldValue) {
            if (this.valueChanged) {
              this.suggestions = []
              this.valueChanged = false
              return
            }
            if (!this.focused) {
                this.suggestions = []
                return
            }
            if (newValue.length < 1 || oldValue === '' || newValue === oldValue) {
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
    <div>
        <span v-if="inputModel !== '' && !readonly" class="clear-btn" @click="clear"><i class="fa fa-times"></i></span>
        <div v-show="!focused" @click="doFocus" @focus="doFocus" @keyup.enter=" focus"
              class="clickable"
              tabindex="0">
          <text-highlight :queries="highlights">{{ value || '&nbsp;' }}</text-highlight>
        </div>
        <div class="focused" v-show="focused">
            <input type="text" ref="input" class="text-input" :class="{'error': error}" v-model="inputModel"
                   @focus="doFocus"
                   @blur="blur"
                   @keyup.46="clear"
                   @keydown.down="keyDownPressed"
                   @keydown.up="keyUpPressed"
                   @keydown.left="keyLeftPressed"
                   @keydown.right="keyRightPressed"
                   @keyup.enter="update" @keyup.esc="blur">
            <div class="dropdown" v-if="suggestions.length > 0">
                <ul>
                    <li v-for="(s,i) in suggestions" @mousedown.prevent @click="clickSuggestion(s)" @mouseover="selectedSuggestionIdx = i"
                        :class="{'selected': (i === selectedSuggestionIdx)}">
                        {{ s }}
                    </li>
                </ul>
            </div>
        </div>
    </div>
</template>
