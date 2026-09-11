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
// 基础语言（xml、yaml、ini 等）仅做词法高亮，不需要 worker
import 'monaco-editor/esm/vs/basic-languages/monaco.contribution'
// json 高亮与校验由 language service 提供，需要 json worker
import 'monaco-editor/esm/vs/language/json/monaco.contribution'
import jsonWorker from 'monaco-editor/esm/vs/language/json/json.worker?worker'
import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker'
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'

// 只用到默认 worker 与 json worker
self.MonacoEnvironment = {
    getWorker(_workerId, label) {
        return label === 'json' ? new jsonWorker() : new editorWorker()
    },
}

// monaco 未内置 toml，这里补一个词法高亮
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
            // 表头只能在行首：[table] 与 [[array of tables]]
            [/^\s*\[\[[^\]]*\]\]/, 'type'],
            [/^\s*\[[^\]]*\]/, 'type'],
            // 键，即 key = value 中的 key
            [/[A-Za-z0-9_-]+(?=\s*=)/, 'key'],
            // 日期时间（RFC 3339），必须在数字之前
            [/\d{4}-\d{2}-\d{2}(?:[Tt ]\d{2}:\d{2}(?::\d{2}(?:\.\d+)?)?(?:[Zz]|[+-]\d{2}:\d{2})?)?/, 'number'],
            [/\d{2}:\d{2}:\d{2}(?:\.\d+)?/, 'number'],
            [/(?:true|false)(?![A-Za-z0-9_])/, 'keyword'],
            // 数字：十进制、十六进制、八进制、二进制
            [/[+-]?(?:0[xX][0-9A-Fa-f_]+|0[oO][0-7_]+|0[bB][01_]+|(?:0|[1-9][0-9_]*)(?:\.[0-9_]+)?(?:[eE][+-]?[0-9_]+)?)/, 'number'],
            // 多行字符串要放在单行字符串之前
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
            // 未闭合的字符串换行后回到普通状态，避免污染后续内容
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

// txt 用 monaco 内置的纯文本语言
const toMonacoLanguage = (value) => (value === 'txt' ? 'plaintext' : value)

const props = defineProps({
    modelValue: {
        type: String,
        default: '',
    },
    // 语法类型：json / xml / yaml / ini / toml / txt
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

// 父组件传入的新值（与编辑器当前内容不同才同步，避免光标跳动）
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
/* 标签栏只用来切换高亮，编辑器是独立区域，收紧默认间距 */
.meta-data-editor__tabs :deep(.el-tabs__header) {
    margin-bottom: 8px;
}

.meta-data-editor__body {
    width: 100%;
    overflow: hidden;
    border: 1px solid var(--el-border-color);
    border-radius: var(--el-border-radius-base);
}
</style>
