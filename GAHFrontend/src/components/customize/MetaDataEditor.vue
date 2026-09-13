<template>
    <div class="meta-data-editor">
        <el-tabs v-model="currentLanguage" class="meta-data-editor__tabs">
            <el-tab-pane v-for="item in LANGUAGE_OPTIONS" :key="item.value" :label="item.label" :name="item.value" />
        </el-tabs>
        <div ref="containerRef" class="meta-data-editor__body" :style="{ height }"></div>
    </div>
</template>

<script setup>
import * as monaco from 'monaco-editor/esm/vs/editor/editor.api'
import 'monaco-editor/esm/vs/editor/editor.all'
// Basic languages (xml, yaml, ini, etc.) only need syntax highlighting, no worker required
import 'monaco-editor/esm/vs/basic-languages/monaco.contribution'
// json highlighting and validation are provided by the language service, which needs the json worker
import 'monaco-editor/esm/vs/language/json/monaco.contribution'
import jsonWorker from 'monaco-editor/esm/vs/language/json/json.worker?worker'
import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker'
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'

// Only the default worker and the json worker are used
self.MonacoEnvironment = {
    getWorker(_workerId, label) {
        return label === 'json' ? new jsonWorker() : new editorWorker()
    },
}

// monaco has no built-in toml, so add syntax highlighting here
const tomlLanguage = {
    defaultToken: '',
    tokenPostfix: '.toml',
    brackets: [
        { open: '[', close: ']', token: 'delimiter.square' },
        { open: '{', close: '}', token: 'delimiter.curly' },
    ],
    tokenizer: {
        root: [
            [/[ \t\r\n]+/, 'white'],
            [/#.*$/, 'comment'],
            // Table headers are only allowed at the start of a line: [table] and [[array of tables]]
            [/^\s*\[\[[^\]]*\]\]/, 'type'],
            [/^\s*\[[^\]]*\]/, 'type'],
            // Key, i.e. the key in key = value
            [/[A-Za-z0-9_-]+(?=\s*=)/, 'key'],
            // Date-time (RFC 3339), must come before numbers
            [/\d{4}-\d{2}-\d{2}(?:[Tt ]\d{2}:\d{2}(?::\d{2}(?:\.\d+)?)?(?:[Zz]|[+-]\d{2}:\d{2})?)?/, 'number'],
            [/\d{2}:\d{2}:\d{2}(?:\.\d+)?/, 'number'],
            [/(?:true|false)(?![A-Za-z0-9_])/, 'keyword'],
            // Numbers: decimal, hexadecimal, octal, binary
            [/[+-]?(?:0[xX][0-9A-Fa-f_]+|0[oO][0-7_]+|0[bB][01_]+|(?:0|[1-9][0-9_]*)(?:\.[0-9_]+)?(?:[eE][+-]?[0-9_]+)?)/, 'number'],
            // Multi-line strings must come before single-line strings
            [/"""/, { token: 'string.quote', next: '@multiLineStringDouble' }],
            [/'''/, { token: 'string.quote', next: '@multiLineStringSingle' }],
            [/"/, { token: 'string.quote', next: '@stringDouble' }],
            [/'/, { token: 'string.quote', next: '@stringSingle' }],
            [/[[\]{}]/, '@brackets'],
            [/[=,.]/, 'delimiter'],
        ],
        stringDouble: [
            [/[^\\"\n]+/, 'string'],
            [/\\(?:[btnfr"\\]|u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8})/, 'string.escape'],
            [/"/, { token: 'string.quote', next: '@pop' }],
            // Unterminated strings return to the normal state after a line break, to avoid polluting subsequent content
            [/\r?\n/, { token: 'string', next: '@pop' }],
        ],
        stringSingle: [
            [/[^'\n]+/, 'string'],
            [/'/, { token: 'string.quote', next: '@pop' }],
            [/\r?\n/, { token: 'string', next: '@pop' }],
        ],
        multiLineStringDouble: [
            [/[^"\\]+/, 'string'],
            [/\\(?:[btnfr"\\]|u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8})/, 'string.escape'],
            [/"""/, { token: 'string.quote', next: '@pop' }],
            [/"/, 'string'],
        ],
        multiLineStringSingle: [
            [/[^']+/, 'string'],
            [/'''/, { token: 'string.quote', next: '@pop' }],
            [/'/, 'string'],
        ],
    },
}

monaco.languages.register({ id: 'toml' })
monaco.languages.setMonarchTokensProvider('toml', tomlLanguage)

const LANGUAGE_OPTIONS = [
    { label: 'JSON', value: 'json' },
    { label: 'XML', value: 'xml' },
    { label: 'YAML', value: 'yaml' },
    { label: 'INI', value: 'ini' },
    { label: 'TOML', value: 'toml' },
    { label: 'TXT', value: 'txt' },
]

// txt uses monaco's built-in plain text language
const toMonacoLanguage = (value) => (value === 'txt' ? 'plaintext' : value)

const props = defineProps({
    modelValue: {
        type: String,
        default: '',
    },
    // Syntax type: json / xml / yaml / ini / toml / txt
    language: {
        type: String,
        default: 'json',
    },
    height: {
        type: String,
        default: '300px',
    },
    readonly: {
        type: Boolean,
        default: false,
    },
})

const emit = defineEmits(['update:modelValue', 'update:language'])

const containerRef = ref(null)
const currentLanguage = ref(props.language)

let editor = null
let model = null

onMounted(() => {
    model = monaco.editor.createModel(props.modelValue ?? '', toMonacoLanguage(currentLanguage.value))
    editor = monaco.editor.create(containerRef.value, {
        model,
        theme: 'vs',
        automaticLayout: true,
        minimap: { enabled: false },
        fontSize: 13,
        tabSize: 2,
        scrollBeyondLastLine: false,
        wordWrap: 'on',
        readOnly: props.readonly,
    })
    editor.onDidChangeModelContent(() => {
        emit('update:modelValue', editor.getValue())
    })
})

// New value passed from the parent (sync only when different from the editor content, to avoid cursor jumps)
watch(
    () => props.modelValue,
    (val) => {
        if (!editor || val === editor.getValue()) return
        editor.setValue(val ?? '')
    }
)

watch(currentLanguage, (val) => {
    if (!model) return
    monaco.editor.setModelLanguage(model, toMonacoLanguage(val))
    emit('update:language', val)
})

watch(
    () => props.language,
    (val) => {
        if (val !== currentLanguage.value) currentLanguage.value = val
    }
)

watch(
    () => props.readonly,
    (val) => editor?.updateOptions({ readOnly: val })
)

onBeforeUnmount(() => {
    editor?.dispose()
    model?.dispose()
    editor = null
    model = null
})
</script>

<style scoped>
/* Inside a flex container (such as el-form-item__content) a definite width is required:
   otherwise the width is content-driven, and monaco writes the pixel width it measured back into the DOM,
   creating positive feedback that widens it frame by frame */
.meta-data-editor {
    width: 100%;
    min-width: 0;
}

/* The tab bar is only used to switch highlighting; the editor is a separate area, so tighten the default spacing */
.meta-data-editor__tabs :deep(.el-tabs__header) {
    margin-bottom: 8px;
}

.meta-data-editor__body {
    width: 100%;
    /* The project has no global box-sizing reset; use border-box so the 1px border does not add extra width */
    box-sizing: border-box;
    overflow: hidden;
    border: 1px solid var(--el-border-color);
    border-radius: var(--el-border-radius-base);
}
</style>
