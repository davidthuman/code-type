<script setup>

import { ref, computed, onMounted } from 'vue';
import { useMagicKeys, whenever, useEventListener } from '@vueuse/core';

const { current } = useMagicKeys();

class Character {

    /**
     * 
     * @param {string} character 
     */
    constructor(character, index, extra = false) {
        this.characterRaw = character;
        this.index = index;
        this.input = null;
        this.extra = extra;
    }

    finished() {
        return this.input != null;
    }

    setInput(character) {
        this.input = character;
    }

    backspace() {
        this.input = null;
    }

    correct() {
        return !this.extra && this.input == this.characterRaw;
    }
}

class Word {
    /**
     * 
     * @param {string} word 
     */
    constructor(word, index) {
        this.wordRaw = word;
        this.index = index;
        this.characters = word.split('').map((character, index) => new Character(character, index));
        this.currentCharacter = 0;
        this.extraCharacters = [];
        this.passToExtraCharacters = false;
    }

    finished() {
        return this.characters.reduce(
            (accumulator, currentCharacter) => accumulator && currentCharacter.finished(),
            true
        )
    }

    setInput(character) {

        if (!this.passToExtraCharacters && this.currentCharacter < this.characters.length) {
            this.characters[this.currentCharacter].setInput(character);
            if (this.currentCharacter == this.characters.length - 1) {
                this.passToExtraCharacters = true;
                return true;
            } else {
                this.currentCharacter++;
            }
        } else {
            this.extraCharacters.push(new Character(character, this.characters.length + this.extraCharacters.length + 1, true));
        }
    }

    backspace() {

        // If current character is the first character in the word
        if (this.currentCharacter == 0) {
            return true;
        } else if (this.extraCharacters.length == 0) {
            if (this.passToExtraCharacters) {
                this.characters[this.currentCharacter].backspace();
                this.passToExtraCharacters = false;
            } else {
                this.characters[this.currentCharacter - 1].backspace();
                this.currentCharacter--;
            }
            
        } else {
            this.extraCharacters.pop();
        }
    }

    correct() {
        return this.characters.reduce(
            (accumulator, currentCharacter) => accumulator && currentCharacter.correct(),
            true
        )
    }

    allCharacters() {
        return [...this.characters, ...this.extraCharacters]
    }
}

class Line {

    /**
     * 
     * @param {string} line
     */
    constructor(line, index) {
        this.lineRaw = line;
        this.index = index;
        this.words = line.split(' ').flatMap((word, index) => [new Word(word, index*2), new Word(' ', (index*2)+1)]).filter((word) => word.wordRaw != '');
        this.words[this.words.length -1] = new Word(' \n', this.words[this.words.length -1].index)
        this.currentWord = 0;
    }

    finished() {
        return this.words.reduce(
            (accumulator, currentWord) => accumulator && currentWord.finished(),
            true
        )
    }

    setInput(character) {

        // If the current word is a Space, pass the input to the previous word
        if (this.words[this.currentWord].wordRaw == ' ') {
            this.words[this.currentWord - 1].setInput(character);
        } else {
            // Pass the input to the current word
            const goToNextWord = this.words[this.currentWord].setInput(character);
            // If the current word is filled, go to next word
            if (goToNextWord) {
                this.currentWord++;
            }
        }    
    }

    backspace() {

        // Current word is not the first on the line AND is a Space AND the previous word has extra characters
        if (this.currentWord != 0 && this.words[this.currentWord].wordRaw == ' ' && this.words[this.currentWord - 1].extraCharacters.length > 0) {
            const previousWord = this.words[this.currentWord - 1].backspace();
            return false;
        } 
        // Current word is not the first on the line AND is a Space and the previous word has no extra characters
        else if (this.currentWord != 0 && this.words[this.currentWord].wordRaw == ' ' && this.words[this.currentWord - 1].extraCharacters.length == 0) {
            const previousWord = this.words[this.currentWord - 1].backspace();
            this.currentWord--;
            return false;
        } else {
            // Current word is regular
            const previousWord = this.words[this.currentWord].backspace();

            if (this.currentWord == 0 && previousWord) {
                return true;
            } else if (previousWord) {
                this.currentWord--;
                return false;
            }  else {
                return false;
            }
        }
    }

    space() {
        // Space should go to the next non-Space word

        if (this.currentWord == this.words.length - 1) {
            return;
        }

        // Go to next word
        this.currentWord++;
        // If the next word is a Space, go to the next word after that
        if (this.words[this.currentWord].wordRaw == ' ') {
            this.currentWord++;
        }
    }

    tab() {
        // For Tab-Length (4) go to next work if the current word is a Space
        for (let i = 0; i < 4; i++) {
            if (this.words[this.currentWord].wordRaw == ' ') {
                this.currentWord++;
            }
        }
    }

    correct() {
        return this.words.reduce(
            (accumulator, currentWord) => accumulator && currentWord.correct(),
            true
        )
    }

}

class Prompt {

    /**
     * 
     * @param {string} prompt 
     */
    constructor(prompt) {
        this.promptRaw = prompt;
        this.lines = prompt.split('\n').map((line, index) => new Line(line, index));
        this.currentLine = 0;
    }

    finished() {
        return this.lines.reduce(
            (accumulator, currentLine) => accumulator && currentLine.finished(),
            true
        )
    }

    setInput(character) {
        // Pass input to current line
        this.lines[this.currentLine].setInput(character);
    }

    backspace() {
        // Pass backspace to current line
        const goToPreviousLine = this.lines[this.currentLine].backspace();

        // Go to the previous line if the current line is not the first AND current line has no inputs
        if (this.currentLine != 0 && goToPreviousLine) {
            this.currentLine--;
        }
    }

    enter() {
        // Go to next line if it is not the last line
        if (this.currentLine < this.lines.length - 1) {
            this.currentLine++;
        }
    }

    space() {
        // Pass space to current line
        this.lines[this.currentLine].space();
    }

    tab() {
        // Pass tab to current line
        this.lines[this.currentLine].tab();
    }

    correct() {
        return this.lines.reduce(
            (accumulator, currentLine) => accumulator && currentLine.correct(),
            true
        )
    }

    currentCharacter(lineIndex, wordIndex, characterIndex) {
        const line = this.lines[this.currentLine];
        const word = line.words[line.currentWord];
        var character;
        if (word.extraCharacters.length > 0) {
            character = word.extraCharacters[word.extraCharacters.length - 1];
        } else {
            character = word.characters[word.currentCharacter];
        }
        return lineIndex == line.index && wordIndex == word.index && characterIndex == character.index;
    }

    passCharacter(lineIndex, wordIndex, characterIndex) {
        const line = this.lines[this.currentLine];
        const word = line.words[line.currentWord];
        var character;
        if (word.extraCharacters.length > 0) {
            character = word.extraCharacters[word.extraCharacters.length - 1];
        } else {
            character = word.characters[word.currentCharacter];
        }

        if (lineIndex < line.index) {
            return true;
        } else if (lineIndex > line.index) {
            return false;
        } else {

            if (wordIndex < word.index) {
                return true;
            } else if (wordIndex > word.index) {
                return false;
            } else {

                if (characterIndex < character.index) {
                    return true;
                } else if (characterIndex > character.index) {
                    return false;
                }

            }

        }
    }

}

const response = await fetch(
  '/test'
)

const test = await response.text();

const prompt = ref(new Prompt(test));

console.log(prompt.value);

useEventListener(window, 'keydown', (e) => {

    //var [lineIndex, wordIndex, characterIndex] = prompt.value.currentCharacter();
    //var characterEl = document.getElementById(`line-${lineIndex}-word-${wordIndex}-character-${characterIndex}`);
    //characterEl.classList.remove("current-character");

    switch (e.key) {
        case "Tab":
            e.preventDefault();
            prompt.value.tab();
            break;
        case "Backspace":
            prompt.value.backspace();
            break;
        case "Enter":
            prompt.value.enter();
            break;
        case " ":
            prompt.value.space();
            break;
        case "Shift":
            break;
        case "CapsLock":
            break;
        case "'":
            e.preventDefault()
            prompt.value.setInput("'");
            break;
        default:
            prompt.value.setInput(e.key);
            break;
    }

    const lastLine = prompt.value.lines.length - 1;
    const lastWord = prompt.value.lines[lastLine].words.length - 1;
    const lastCharacter = prompt.value.lines[lastLine].words[lastWord].characters.length - 1;
    if (prompt.value.currentCharacter(lastLine, lastWord, lastCharacter)) {
        console.log("Finished!");
    }

    //var [lineIndex, wordIndex, characterIndex] = prompt.value.currentCharacter();
    //var characterEl = document.getElementById(`line-${lineIndex}-word-${wordIndex}-character-${characterIndex}`);
    //characterEl.classList.add("current-character");

    //console.log(characterEl);
});

import hljs from 'highlight.js/lib/common';

/*
const code = computed(() => {
  const typed =  keys.value.join('')
  const all = 
    hljs.highlight(typed, { language: 'python'}).value +
    `<span style="animation: blink 0.60s infinite alternate;" class="cursor">${test[typed.length]}</span>` +
    `<span class="opacity-50">${hljs.highlight(test.slice(typed.length+1), { language: 'plaintext'}).value}</span>`;
  return all;
})
*/


onMounted(() => {
    //hljs.highlightAll();
    //const characterEl = document.getElementById(`line-0-word-0-character-0`);
    //characterEl.classList.add("current-character");
})


</script>

<template>
  <div class="p-5">
    <div class="font-bold">Code Type</div>
    <br />
    <div>Time: </div>
    <br />
    <div class="font-mono" style="white-space: pre-wrap;">
        <div v-for="line in prompt.lines" :id="`line-${line.index}`" :key="line.index">
            <span v-for="word in line.words" :id="`line-${line.index}-word-${word.index}`" :key="word.index">
                <span 
                    v-for="character in word.allCharacters()" 
                    :id="`line-${line.index}-word-${word.index}-character-${character.index}`" 
                    :key="character.index" 
                    v-bind:class="{
                        'finished': character.finished() || prompt.passCharacter(line.index, word.index, character.index), 
                        'not-finished': (!character.finished() && !character.extra) || !prompt.passCharacter(line.index, word.index, character.index), 
                        'correct': character.correct() && (character.finished() || prompt.passCharacter(line.index, word.index, character.index)), 
                        'not-correct': !character.correct() && (character.finished() || prompt.passCharacter(line.index, word.index, character.index)), 
                        'extra': character.extra, 
                        'current-character': prompt.currentCharacter(line.index, word.index, character.index)
                    }">
                    {{ character.finished() ? character.input : character.characterRaw }}
                </span>
            </span>
        </div>
    </div>
  </div>
</template>

<style>

.current-character {
    animation: blink 0.60s infinite alternate;
}

@keyframes blink {
  from { background-color: rgba(255, 255, 255, 0); }
  to { background-color: rgb(183, 215, 255); }
}

.not-finished {
    color: rgba(0, 0, 0, 0.5);
}

.not-correct {
    color: rgb(255, 94, 94);
}

.extra {
    color: rgb(255, 94, 94);
}

</style>
