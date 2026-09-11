<template>
    <el-card>
        <template #header>
            操作
        </template>
        <p><el-input placeholder="管理员用户名"></el-input></p>
        <p><el-input placeholder="管理员UUHash"></el-input></p>
        <p>权限筛选：</p>
        <el-checkbox label="是否可以添加管理员"></el-checkbox>
        <el-checkbox label="是否可以删除管理员"></el-checkbox>
        <el-checkbox label="是否可以修改管理员"></el-checkbox>
        <el-checkbox label="是否可以获取管理员列表"></el-checkbox>
        <el-checkbox label="是否可以操作用户"></el-checkbox>
        <el-checkbox label="是否可以操作角色"></el-checkbox>
        <template #footer>
            <div style="text-align: right">
                <el-button type="primary" @click="handleSearch">查询</el-button>
                <el-button type="success" @click="handleAdd">添加管理员</el-button>
            </div>
        </template>
    </el-card>
    <br>
    <el-card>
        <template #header>
            管理员列表
        </template>
        <div style="height: 600px">
            <el-auto-resizer>
                <template #default="{ height, width }">
                    <el-table-v2 :columns="columns" :data="adminList" :width="width" :height="height" fixed />
                </template>
            </el-auto-resizer>
        </div>
    </el-card>
</template>
<script setup>
import { post } from '.././../../lib/request.js'
import rangeAdmin from '.././../../lib/rangeAdmin.js'
import { h, ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElCheckbox, ElButton } from 'element-plus'

const router = useRouter()
const adminInfo = ref({})
const adminList = ref([])

const canGetAdmin = computed(() => {
    if (adminInfo.value.is_root) return true
    return Boolean(adminInfo.value.permission?.can_get_admin)
})

// 权限列统一渲染：命中为“是”，否则“否”
const permCell = (key) => ({ rowData }) => h('span', rowData.Permission?.[key] ? '是' : '否')

// 勾选状态：以行 ID 记录
const selectedIds = ref([])
const isSelected = (id) => selectedIds.value.includes(id)
const allSelected = computed(
    () => adminList.value.length > 0 && selectedIds.value.length === adminList.value.length
)

const toggleRow = (id) => {
    selectedIds.value = isSelected(id)
        ? selectedIds.value.filter((item) => item !== id)
        : [...selectedIds.value, id]
}

const toggleAll = () => {
    selectedIds.value = allSelected.value ? [] : adminList.value.map((item) => item.ID)
}

// 行操作：待接入 /admin/delete 与 /admin/edit
const handleDelete = (row) => {
}

const handleEdit = (row) => {
}

const columns = [
    {
        key: 'selection',
        width: 50,
        headerRenderer: () => h(ElCheckbox, {
            modelValue: allSelected.value,
            'onUpdate:modelValue': toggleAll,
        }),
        cellRenderer: ({ rowData }) => h(ElCheckbox, {
            modelValue: isSelected(rowData.ID),
            'onUpdate:modelValue': () => toggleRow(rowData.ID),
        }),
    },
    { key: 'ID', dataKey: 'ID', title: 'ID', width: 70 },
    { key: 'Username', dataKey: 'Username', title: '用户名', width: 140 },
    { key: 'UUHash', dataKey: 'UUHash', title: 'UUHash', width: 550 },
    { key: 'can_add_admin', title: '添加管理员', width: 110, cellRenderer: permCell('can_add_admin') },
    { key: 'can_delete_admin', title: '删除管理员', width: 110, cellRenderer: permCell('can_delete_admin') },
    { key: 'can_edit_admin', title: '修改管理员', width: 110, cellRenderer: permCell('can_edit_admin') },
    { key: 'can_get_admin', title: '查看管理员', width: 110, cellRenderer: permCell('can_get_admin') },
    { key: 'can_operate_user', title: '操作用户', width: 100, cellRenderer: permCell('can_operate_user') },
    { key: 'can_operate_character', title: '操作角色', width: 100, cellRenderer: permCell('can_operate_character') },
    {
        key: 'actions',
        title: '操作',
        width: 200,
        cellRenderer: ({ rowData }) => h('div', { style: 'display: flex; gap: 8px;' }, [
            h(ElButton, { type: 'danger', size: 'small', onClick: () => handleDelete(rowData) }, () => '删除用户'),
            h(ElButton, { type: 'primary', size: 'small', onClick: () => handleEdit(rowData) }, () => '编辑用户'),
        ]),
    },
]

const range = async () => {
    try {
        const list = await rangeAdmin({}, 1, 100)
        adminList.value = list.sort((a, b) => a.ID - b.ID)
    } catch {
        // 401 已由 request 封装自动跳转登录页
    }
}

onMounted(async () => {
    try {
        const res = await post('/admin/info')
        adminInfo.value = res.data
    } catch {
        // 401 已由 request 封装自动跳转登录页，网络错误时不再继续鉴权
        return
    }
    if (!canGetAdmin.value) {
        router.push('/')
        return
    }
    range()
})
</script>
